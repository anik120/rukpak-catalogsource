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

	opm "github.com/operator-framework/operator-registry/alpha/action"
	bundleProperty "github.com/operator-framework/operator-registry/alpha/property"
	"github.com/operator-framework/operator-registry/pkg/image/containerdregistry"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"rukpak-catalogsource/api/v1alpha1"
	corev1alpha1 "rukpak-catalogsource/api/v1alpha1"
)

// CatalogSourceReconciler reconciles a CatalogSource object
type CatalogSourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.rukpak.io,resources=catalogsources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.rukpak.io,resources=catalogsources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.rukpak.io,resources=catalogsources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CatalogSource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *CatalogSourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	catalogSource := v1alpha1.CatalogSource{}
	if err := r.Client.Get(ctx, req.NamespacedName, &catalogSource); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// IMPORTANT TODO: This implementation of containerdregistry requires privileged perm to create a CacheDir
	// Figure out a way to not use the CacheDir so that container can be run in non-privileged mode.
	reg, err := containerdregistry.NewRegistry()
	// defer reg.Destroy()
	if err != nil {
		return ctrl.Result{}, err
	}
	imageRenderer := opm.Render{Refs: []string{catalogSource.Spec.Image}, Registry: reg}
	declCfg, err := imageRenderer.Run(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}
	cache := v1alpha1.CatalogCache{ObjectMeta: metav1.ObjectMeta{
		Name:      req.Name,
		Namespace: req.Namespace,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: catalogSource.APIVersion,
				Kind:       catalogSource.Kind,
				Name:       req.Name,
				UID:        catalogSource.UID,
			},
		}}}
	for _, bundle := range declCfg.Bundles {
		operator := v1alpha1.Operator{
			Name:       bundle.Name,
			Package:    bundle.Package,
			BundlePath: bundle.Image,
		}
		props, _ := bundleProperty.Parse(bundle.Properties)
		providedGVKs := []corev1alpha1.APIKey{}
		for _, gvk := range props.GVKs {
			providedGVKs = append(providedGVKs, corev1alpha1.APIKey{Group: gvk.Group, Kind: gvk.Kind, Version: gvk.Version})
		}
		requiredGVKs := []corev1alpha1.APIKey{}
		for _, gvk := range props.GVKsRequired {
			requiredGVKs = append(requiredGVKs, corev1alpha1.APIKey{Group: gvk.Group, Kind: gvk.Kind, Version: gvk.Version})
		}
		operator.ProvidedAPIs = providedGVKs
		operator.RequiredAPIs = requiredGVKs
		cache.Spec.Operators = append(cache.Spec.Operators, operator)
	}

	if err := r.Client.Create(ctx, &cache); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CatalogSourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.CatalogSource{}).
		Complete(r)
}
