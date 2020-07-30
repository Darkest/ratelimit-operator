package controller

import (
	"ratelimit-operator/pkg/controller/ratelimiterconfig"
	"ratelimit-operator/pkg/controller/ratelimiter"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager) error {
	AddToManagerFuncs = append(AddToManagerFuncs, ratelimiter.Add)
	AddToManagerFuncs = append(AddToManagerFuncs, ratelimiterconfig.Add)

	for _, f := range AddToManagerFuncs {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
