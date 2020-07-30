package ratelimiter

import (
	corev1 "k8s.io/api/core/v1"
	v1 "ratelimit-operator/pkg/apis/operators/v1"
	"ratelimit-operator/pkg/constants"
	"ratelimit-operator/pkg/utils"
)

func (r *ReconcileRateLimiter) BuildRedisContainer() corev1.Container {
	return corev1.Container{
		Name:    "redis",
		Image:   constants.REDIS_IMAGE,
		Command: []string{"redis-server"},
		Args:    []string{"--save", "\"\"", "--appendonly", "no"},
	}
}

func (r *ReconcileRateLimiter) BuildServiceContainer(instance *v1.RateLimiter) corev1.Container {
	defaultRedisUrl := r.buildRedisUrl()
	return corev1.Container{
		Name:  instance.Name,
		Image: constants.RATE_LIMITER_SERVICE_IMAGE,
		Ports: []corev1.ContainerPort{{
			ContainerPort: r.buildRateLimiterServicePort(instance),
			Protocol:      corev1.ProtocolTCP,
		}},
		Env: []corev1.EnvVar{
			{
				Name:  "LOG_LEVEL",
				Value: utils.DefaultIfEmpty(string(*instance.Spec.LogLevel), string(v1.INFO)),
			},
			{
				Name:  "REDIS_SOCKET_TYPE",
				Value: "TCP",
			},
			{
				Name:  "REDIS_URL",
				Value: defaultRedisUrl,
			},
			{
				Name:  "RUNTIME_IGNOREDOTFILES",
				Value: "true",
			},
			{
				Name:  "RUNTIME_ROOT",
				Value: "/home/user/src/runtime/data",
			},
			{
				Name:  "RUNTIME_SUBDIRECTORY",
				Value: "ratelimit",
			},
			{
				Name:  "RUNTIME_WATCH_ROOT",
				Value: "false",
			},
			{
				Name:  "USE_STATSD",
				Value: "false",
			},
		},
		VolumeMounts: []corev1.VolumeMount{{
			Name:      "config",
			MountPath: "/home/user/src/runtime/data/ratelimit/config",
		}},
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
	}
}
