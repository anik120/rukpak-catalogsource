/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"rukpak-catalogsource/api/v1alpha1"
	corev1alpha1 "rukpak-catalogsource/api/v1alpha1"
)

const (
	conditionLastUpdateTime = "CacheUpdatedAt"
)

// CatalogCacheReconciler reconciles a CatalogCache object
type CatalogCacheReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.rukpak.io,resources=catalogcaches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.rukpak.io,resources=catalogcaches/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.rukpak.io,resources=catalogcaches/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CatalogCache object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *CatalogCacheReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	cache := v1alpha1.CatalogCache{}
	if err := r.Client.Get(ctx, req.NamespacedName, &cache); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(cache.Status.Conditions) == 0 {
		cache.Status.Conditions = append(cache.Status.Conditions, v1.Condition{Type: conditionLastUpdateTime, LastTransitionTime: v1.Now()})
	} else {
		for _, cond := range cache.Status.Conditions {
			if cond.Type == conditionLastUpdateTime {
				cond.LastTransitionTime = v1.Now()
			}
		}
	}

	r.Client.Status().Update(ctx, &cache)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CatalogCacheReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.CatalogCache{}).
		Complete(r)
}
