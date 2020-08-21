package ratelimiterconfig

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	networking "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/utils"
	"testing"
)

func Test_BuildEnvoyFilter_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build envoy filter", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.GATEWAY,
				RateLimitProperty: v1.RateLimitProperty{
					Domain: utils.BuildRandomString(3),
				},
				FailureModeDeny: true,
			},
		}
		rateLimiter := buildRateLimiter()

		actualPatch := buildEnvoyFilter(rateLimiterConfig, rateLimiter)

		a.Equal(rateLimiterConfig.Name, actualPatch.ObjectMeta.Name)
		a.Equal(rateLimiterConfig.Namespace, actualPatch.ObjectMeta.Name)
		a.Equal(buildWorkloadSelectorLabels(rateLimiterConfig), actualPatch.Spec.WorkloadSelector.Labels)
		a.Equal(3, len(actualPatch.Spec.ConfigPatches))
		a.Equal(buildHttpFilterPatch(rateLimiterConfig, rateLimiter), actualPatch.Spec.ConfigPatches[0])
		a.Equal(buildClusterPatch(rateLimiter), actualPatch.Spec.ConfigPatches[1])
		a.Equal(buildVirtualHostPatch(rateLimiterConfig), actualPatch.Spec.ConfigPatches[2])
	})
}

func Test_BuildHttpFilterPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for http filter", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.GATEWAY,
				RateLimitProperty: v1.RateLimitProperty{
					Domain: utils.BuildRandomString(3),
				},
				FailureModeDeny: true,
			},
		}
		rateLimiter := buildRateLimiter()

		expectedObjectTypes := &networking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
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
		}

		actualPatch := buildHttpFilterPatch(rateLimiterConfig, rateLimiter)

		a.Equal(networking.EnvoyFilter_HTTP_FILTER, actualPatch.ApplyTo)
		a.IsType(&networking.EnvoyFilter_EnvoyConfigObjectMatch_Listener{}, actualPatch.Match.ObjectTypes)
		a.Equal(expectedObjectTypes, actualPatch.Match.ObjectTypes)
		a.Equal(networking.EnvoyFilter_Patch_INSERT_BEFORE, actualPatch.Patch.Operation)
		a.Equal(convertYaml2Struct(buildHttpFilterPatchValue(rateLimiterConfig, rateLimiter)), actualPatch.Patch.Value)
	})
}

func Test_BuildHttpFilterPatchValue_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch value for http filter", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				RateLimitProperty: v1.RateLimitProperty{
					Domain: utils.BuildRandomString(3),
				},
				FailureModeDeny: true,
			},
		}
		rateLimiter := buildRateLimiter()

		expectedPatchValue := fmt.Sprintf(`
          config:
            domain: %s
            failure_mode_deny: true
            rate_limit_service:
              grpc_service:
                envoy_grpc:
                  cluster_name: %s
                timeout: 0.25s
          name: envoy.rate_limit`,
			rateLimiterConfig.Spec.RateLimitProperty.Domain,
			buildWorkAroundServiceName(rateLimiter))

		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildHttpFilterPatchValue(rateLimiterConfig, rateLimiter)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(expectedPatch, actualPatch)
	})
}

func Test_BuildClusterPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch for cluster", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedObjectTypes := &networking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{
			Cluster: &networking.EnvoyFilter_ClusterMatch{
				Service: buildRateLimiterServiceName(rateLimiter),
			},
		}

		actualPatch := buildClusterPatch(rateLimiter)

		a.Equal(networking.EnvoyFilter_CLUSTER, actualPatch.ApplyTo)
		a.IsType(&networking.EnvoyFilter_EnvoyConfigObjectMatch_Cluster{}, actualPatch.Match.ObjectTypes)
		a.Equal(expectedObjectTypes, actualPatch.Match.ObjectTypes)
		a.Equal(networking.EnvoyFilter_Patch_MERGE, actualPatch.Patch.Operation)
		a.Equal(convertYaml2Struct(buildClusterPatchValue(rateLimiter)), actualPatch.Patch.Value)
	})
}

func Test_BuildClusterPatchValue_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch value for cluster", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedPatchValue := fmt.Sprintf("name: %s", buildWorkAroundServiceName(rateLimiter))
		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildClusterPatchValue(rateLimiter)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(expectedPatch, actualPatch)
	})
}

func Test_BuildVirtualHostPatch_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build virtual host patch", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				RateLimitProperty: v1.RateLimitProperty{
					Descriptors: []v1.Descriptor{{
						Key: utils.BuildRandomString(3),
					}},
				},
			},
		}

		expectedObjectTypes := &networking.EnvoyFilter_EnvoyConfigObjectMatch_RouteConfiguration{
			RouteConfiguration: &networking.EnvoyFilter_RouteConfigurationMatch{
				Vhost: &networking.EnvoyFilter_RouteConfigurationMatch_VirtualHostMatch{
					Name: buildVirtualHostName(rateLimiterConfig),
					Route: &networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch{
						Action: networking.EnvoyFilter_RouteConfigurationMatch_RouteMatch_ANY,
					},
				},
			},
		}

		actualPatch := buildVirtualHostPatch(rateLimiterConfig)

		a.Equal(networking.EnvoyFilter_VIRTUAL_HOST, actualPatch.ApplyTo)
		a.Equal(buildContext(rateLimiterConfig), actualPatch.Match.Context)
		a.IsType(&networking.EnvoyFilter_EnvoyConfigObjectMatch_RouteConfiguration{}, actualPatch.Match.ObjectTypes)
		a.Equal(expectedObjectTypes, actualPatch.Match.ObjectTypes)
		a.Equal(networking.EnvoyFilter_Patch_MERGE, actualPatch.Patch.Operation)
		a.Equal(convertYaml2Struct(buildVirtualHostPatchValue(rateLimiterConfig)), actualPatch.Patch.Value)
	})
}

func Test_BuildVirtualHostPatchValue_HeaderSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch value for virtual host (header)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				RateLimitProperty: v1.RateLimitProperty{
					Descriptors: []v1.Descriptor{{
						Key: utils.BuildRandomString(3),
					}},
				},
			},
		}

		var expectedPatchValue = fmt.Sprintf(`
          rate_limits:
            - actions:
                - request_headers:
                    descriptor_key: %s
                    header_name: %s`,
			rateLimiterConfig.Spec.RateLimitProperty.Descriptors[0].Key,
			rateLimiterConfig.Spec.RateLimitProperty.Descriptors[0].Key)

		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildVirtualHostPatchValue(rateLimiterConfig)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(expectedPatch, actualPatch)
	})
}

func Test_BuildVirtualHostPatchValue_PathSuccess(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build patch value for virtual host (path)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				RateLimitProperty: v1.RateLimitProperty{
					Descriptors: []v1.Descriptor{{
						Key:   "header_match",
						Value: utils.BuildRandomString(3),
					}},
				},
			},
		}

		var expectedPatchValue = fmt.Sprintf(`
          rate_limits:
            - actions:
                - header_value_match: 
                    descriptor_value: %s
                    expect_match: true
                    headers:
                    - exact_match: %s
                      name: ":path"`,
			rateLimiterConfig.Spec.RateLimitProperty.Descriptors[0].Value,
			rateLimiterConfig.Spec.RateLimitProperty.Descriptors[0].Value)

		expectedPatch := convertYaml2Struct(expectedPatchValue)

		actualPatchValue := buildVirtualHostPatchValue(rateLimiterConfig)
		actualPatch := convertYaml2Struct(actualPatchValue)

		a.Equal(expectedPatch, actualPatch)
	})
}

func Test_BuildVirtualHostName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build virtual host name", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				Host: utils.BuildRandomString(3),
				Port: int32(utils.BuildRandomInt(4)),
			},
		}

		expectedResult := fmt.Sprintf("%s:%d", rateLimiterConfig.Spec.Host, rateLimiterConfig.Spec.Port)
		actualResult := buildVirtualHostName(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildContext_Gateway(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build context (gateway)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.GATEWAY,
			},
		}

		expectedResult := networking.EnvoyFilter_GATEWAY
		actualResult := buildContext(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildContext_Sidecar(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build context (sidecar)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo:                v1.SIDECAR,
				WorkloadSelectorLabels: &map[string]string{utils.BuildRandomString(3): utils.BuildRandomString(3)},
			},
		}

		expectedResult := networking.EnvoyFilter_SIDECAR_OUTBOUND
		actualResult := buildContext(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildWorkloadSelectorLabels_IngressGateway(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build workload selector labels (ingressgateway)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo: v1.GATEWAY,
			},
		}

		expectedResult := map[string]string{"istio": "ingressgateway"}
		actualResult := buildWorkloadSelectorLabels(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildWorkloadSelectorLabels_Sidecar(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build workload selector labels (sidecar)", func(t *testing.T) {
		rateLimiterConfig := &v1.RateLimiterConfig{
			Spec: v1.RateLimiterConfigSpec{
				ApplyTo:                v1.SIDECAR,
				WorkloadSelectorLabels: &map[string]string{utils.BuildRandomString(3): utils.BuildRandomString(3)},
			},
		}

		expectedResult := *rateLimiterConfig.Spec.WorkloadSelectorLabels
		actualResult := buildWorkloadSelectorLabels(rateLimiterConfig)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildRateLimiterServiceName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build rate limiter service name", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedResult := fmt.Sprintf("%s.%s.%s", rateLimiter.Name, rateLimiter.Namespace, "svc.cluster.local")
		actualResult := buildRateLimiterServiceName(rateLimiter)

		a.Equal(expectedResult, actualResult)
	})
}

func Test_BuildWorkAroundServiceName_Success(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("success build work around service name", func(t *testing.T) {
		rateLimiter := buildRateLimiter()

		expectedResult := fmt.Sprintf("%s.%s.%s.%s", "patched", rateLimiter.Name, rateLimiter.Namespace, "svc.cluster.local")
		actualResult := buildWorkAroundServiceName(rateLimiter)

		a.Equal(expectedResult, actualResult)
	})
}

func buildRateLimiter() *v1.RateLimiter {
	return &v1.RateLimiter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.BuildRandomString(3),
			Namespace: utils.BuildRandomString(3),
		},
	}
}
