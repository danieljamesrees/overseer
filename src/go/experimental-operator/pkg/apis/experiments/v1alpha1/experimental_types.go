package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ExperimentalSpec defines the desired state of Experimental
type ExperimentalSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	IsPublic         bool  `json:"isPublic"`
	NumberOfReplicas int32 `json:"numberOfReplicas"`
}

// ExperimentalStatus defines the observed state of Experimental
type ExperimentalStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	NumberOfReplicas      int32    `json:"numberOfReplicas"`
	PersistentVolumesUsed []string `json:"persistentVolumesUsed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Experimental is the Schema for the experimentals API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=experimentals,scope=Namespaced
type Experimental struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExperimentalSpec   `json:"spec,omitempty"`
	Status ExperimentalStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExperimentalList contains a list of Experimental
type ExperimentalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Experimental `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Experimental{}, &ExperimentalList{})
}
