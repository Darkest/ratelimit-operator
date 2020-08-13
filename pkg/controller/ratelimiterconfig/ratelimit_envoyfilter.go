package ratelimiterconfig

import (
	"context"
	"fmt"
	"github.com/champly/lib4go/encoding"
	"github.com/ghodss/yaml"
	proto_types "github.com/gogo/protobuf/types"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/controller/ratelimiterconfig/envoyfilter_types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileRateLimiterConfig) reconcileEnvoyFilter(ctx context.Context, instance *v1.RateLimiterConfig) (reconcile.Result, error) {
	envoyFilterFromInstance := r.buildEnvoyFilter(instance)
	reqLogger := log.WithValues("Instance.Name", envoyFilterFromInstance.Name)

	foundEnvoyFilter := &v1alpha3.EnvoyFilter{}

	err := r.client.Get(ctx, types.NamespacedName{Name: envoyFilterFromInstance.Name, Namespace: envoyFilterFromInstance.Namespace}, foundEnvoyFilter)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new EnvoyFilter")
			err = r.client.Create(ctx, envoyFilterFromInstance)
			if err != nil {
				reqLogger.Error(err, "Failed to create new EnvoyFilter")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		} else {
			reqLogger.Error(err, "Failed to get EnvoyFilter")
			return reconcile.Result{}, err
		}
	}

	if !equality.Semantic.DeepEqual(foundEnvoyFilter.Spec, envoyFilterFromInstance.Spec) {
		foundEnvoyFilter.Spec = envoyFilterFromInstance.Spec
		r.client.Update(ctx, foundEnvoyFilter)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileRateLimiterConfig) buildEnvoyFilter(instance *v1.RateLimiterConfig) *v1alpha3.EnvoyFilter {
	envoyFilter := &v1alpha3.EnvoyFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: networking.EnvoyFilter{
			WorkloadSelector: &networking.WorkloadSelector{
				Labels: r.buildWorkloadSelectorLabels(instance),
			},
			ConfigPatches: []*networking.EnvoyFilter_EnvoyConfigObjectPatch{
				{
					ApplyTo: networking.EnvoyFilter_HTTP_FILTER,
					Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: r.buildContext(instance),
						ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
							Listener: &networking.EnvoyFilter_ListenerMatch{
								FilterChain: &networking.EnvoyFilter_ListenerMatch_FilterChainMatch{
									Filter: &networking.EnvoyFilter_ListenerMatch_FilterMatch{
										Name: "envoy.http_connection_manager",
										SubFilter: &networking.EnvoyFilter_ListenerMatch_SubFilterMatch{
											Name: "envoy.router",
										},
									},
								},
							},
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_INSERT_BEFORE,
						Value:     convertYaml2Struct(r.buildHttpFilterPatch(instance)),
					},
				},
				{
					ApplyTo: networking.EnvoyFilter_CLUSTER,
					Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
						ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{
							Cluster: &networking.EnvoyFilter_ClusterMatch{
								Service: r.buildRateLimiterServiceName(),
							},
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_MERGE,
						Value:     convertYaml2Struct(r.buildClusterPatch()),
					},
				},
				{
					ApplyTo: networking.EnvoyFilter_VIRTUAL_HOST,
					Match: &networking.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: r.buildContext(instance),
						ObjectTypes: &networking.EnvoyFilter_EnvoyConfigObjectMatch_RouteConfiguration{
							RouteConfiguration: &networking.EnvoyFilter_RouteConfigurationMatch{
								Vhost: &networking.EnvoyFilter_RouteConfigurationMatch_VirtualHostMatch{
									Name: r.buildVirtualHostName(instance),
									Route: &networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch{
										Action: networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch_ANY,
									},
								},
							},
						},
					},
					Patch: &networking.EnvoyFilter_Patch{
						Operation: networking.EnvoyFilter_Patch_MERGE,
						Value:     convertYaml2Struct(r.buildVirtualHostPatch(instance)),
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(instance, envoyFilter, r.scheme)
	return envoyFilter
}

func convertYaml2Struct(str string) *proto_types.Struct {
	res, _ := encoding.YAML2Struct(str)
	return res
}

func (r *ReconcileRateLimiterConfig) buildHttpFilterPatch(instance *v1.RateLimiterConfig) string {
	values := envoyfilter_types.HttpFilterPatchValues{
		Name: "envoy.rate_limit",
		Config: envoyfilter_types.Config{
			Domain:          instance.Spec.RateLimitProperty.Domain,
			FailureModeDeny: instance.Spec.FailureModeDeny,
			RateLimitService: envoyfilter_types.RateLimitService{
				GrpcService: envoyfilter_types.GrpcService{
					Timeout: "0.25s",
					EnvoyGrpc: envoyfilter_types.EnvoyGrpc{
						ClusterName: r.buildWorkAroundServiceName(),
					},
				},
			},
		},
	}

	res, err := yaml.Marshal(&values)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml for http filter patch")
	}
	return string(res)
}

func (r *ReconcileRateLimiterConfig) buildClusterPatch() string {
	values := envoyfilter_types.ClusterPatchValues{
		Name: r.buildWorkAroundServiceName(),
	}

	res, err := yaml.Marshal(&values)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml for cluster patch")
	}
	return string(res)
}

func (r *ReconcileRateLimiterConfig) buildVirtualHostPatch(instance *v1.RateLimiterConfig) string {
	var rateLimits []envoyfilter_types.RateLimit

	for _, d := range instance.Spec.RateLimitProperty.Descriptors {
		var action envoyfilter_types.Action

		if d.Key == "header_match" {
			action = envoyfilter_types.Action{
				HeaderValueMatch: &envoyfilter_types.HeaderValueMatch{
					DescriptorValue: d.Value,
					ExpectMatch:     true,
					Headers: []envoyfilter_types.Header{{
						ExactMatch: d.Value,
						Name:       ":path",
					}},
				},
				RequestHeaders: nil,
			}
		} else {
			action = envoyfilter_types.Action{
				HeaderValueMatch: nil,
				RequestHeaders: &envoyfilter_types.RequestHeader{
					DescriptorKey: d.Key,
					HeaderName:    d.Key,
				},
			}
		}

		rateLimits = append(rateLimits,
			envoyfilter_types.RateLimit{
				Actions: []envoyfilter_types.Action{action},
			})
	}

	values := envoyfilter_types.VirtualHostPatchValues{RateLimits: rateLimits}

	res, err := yaml.Marshal(&values)
	if err != nil {
		log.Error(err, "Failed to convert object to yaml for virtual host patch")
	}
	return string(res)
}

func (r *ReconcileRateLimiterConfig) buildWorkloadSelectorLabels(instance *v1.RateLimiterConfig) map[string]string {
	if instance.Spec.ApplyTo == v1.GATEWAY {
		return map[string]string{
			"istio": "ingressgateway",
		}
	} else {
		return *instance.Spec.WorkloadSelectorLabels
	}
}

func (r *ReconcileRateLimiterConfig) buildContext(instance *v1.RateLimiterConfig) networking.EnvoyFilter_PatchContext {
	if instance.Spec.ApplyTo == v1.SIDECAR {
		return networking.EnvoyFilter_SIDECAR_OUTBOUND
		//} else if instance.Spec.ApplyTo == v1.SIDECAR_INBOUND {
		//	return networking.EnvoyFilter_SIDECAR_INBOUND
	} else {
		return networking.EnvoyFilter_GATEWAY
	}
}

func (r *ReconcileRateLimiterConfig) buildVirtualHostName(instance *v1.RateLimiterConfig) string {
	return fmt.Sprintf("%s:%d", instance.Spec.Host, instance.Spec.Port)
}

func (r *ReconcileRateLimiterConfig) buildRateLimiterServiceName() string {
	return fmt.Sprintf("%s.%s.%s", r.rateLimiter.Name, r.rateLimiter.Namespace, "svc.cluster.local")
}

func (r *ReconcileRateLimiterConfig) buildWorkAroundServiceName() string {
	return fmt.Sprintf("%s.%s.%s.%s", "patched", r.rateLimiter.Name, r.rateLimiter.Namespace, "svc.cluster.local")
}