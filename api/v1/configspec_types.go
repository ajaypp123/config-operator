/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConfigSpecSpec defines the desired state of ConfigSpec
type ConfigSpecSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// foo is an example field of ConfigSpec. Edit configspec_types.go to remove/update
	// +optional
	Value string `json:"value"`

	// Time is the timestamp of creation/update in RFC3339 format
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=date-time
	Time string `json:"time"`

	// Status indicates the current status of the configuration
	// +kubebuilder:validation:Enum=Pending;Applied;Error
	// +kubebuilder:default=Pending
	Status string `json:"status"`
}

// ConfigSpecStatus defines the observed state of ConfigSpec.
type ConfigSpecStatus struct {
	// Conditions represent the latest available observations of a ConfigSpec's current state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// LastAppliedTime is the timestamp when the configuration was last successfully applied
	// +optional
	LastAppliedTime *metav1.Time `json:"lastAppliedTime,omitempty"`

	// ConfigMapName is the name of the created ConfigMap
	// +optional
	ConfigMapName string `json:"configMapName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ConfigSpec is the Schema for the configspecs API
type ConfigSpec struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of ConfigSpec
	// +required
	Spec ConfigSpecSpec `json:"spec"`

	// status defines the observed state of ConfigSpec
	// +optional
	Status ConfigSpecStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// ConfigSpecList contains a list of ConfigSpec
type ConfigSpecList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigSpec `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigSpec{}, &ConfigSpecList{})
}
