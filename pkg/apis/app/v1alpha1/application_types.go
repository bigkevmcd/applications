package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplicationSpec defines the desired state of Application
// +k8s:openapi-gen=true
type ApplicationSpec struct {
	Environment map[string]string `json:"environment,omitempty"`

	// +kubebuilder:validation:MinItems=1
	Processes []ProcessSpec `json:"processes,omitempty"`

	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// ProcessSpec defines the state for a process.
// +k8s:openapi-gen=true
type ProcessSpec struct {
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Pattern=.+:.+
	Image string `json:"image,omitempty"`
	Port  int32  `json:"port"`
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas"`
}

// ApplicationStatus defines the observed state of Application
// +k8s:openapi-gen=true
type ApplicationStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Application is the Schema for the applications API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
