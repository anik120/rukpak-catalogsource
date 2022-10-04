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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CatalogCacheSpec defines the desired state of CatalogCache
type CatalogCacheSpec struct {
	Operators []Operator `json:"operators,omitempty"`
}

// CatalogCacheStatus defines the observed state of CatalogCache
type CatalogCacheStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CatalogCache is the Schema for the catalogcaches API
type CatalogCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CatalogCacheSpec   `json:"spec,omitempty"`
	Status CatalogCacheStatus `json:"status,omitempty"`
}

// Operator
type Operator struct {
	Name         string   `json:"name"`
	Package      string   `json:"package"`
	Version      string   `json:"version"`
	ProvidedAPIs []APIKey `json:"providedapis"`
	RequiredAPIs []APIKey `json:"requiredapis,omitempty"`
	BundlePath   string   `json:"bundlepath"`
}

// APIKey
type APIKey struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

//+kubebuilder:object:root=true

// CatalogCacheList contains a list of CatalogCache
type CatalogCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CatalogCache `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CatalogCache{}, &CatalogCacheList{})
}
