// SPDX-License-Identifier: BSD-3-Clause

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/netbirdio/kubernetes-operator/internal/util"
)

// NBGroupSpec defines the desired state of NBGroup.
type NBGroupSpec struct {
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// NBGroupStatus defines the observed state of NBGroup.
type NBGroupStatus struct {
	// +optional
	GroupID *string `json:"groupID"`
	// +optional
	Conditions []NBCondition `json:"conditions,omitempty"`
}

// Equal returns if NBGroupStatus is equal to this one. Condition timestamps are
// ignored, since they are restamped on every reconcile and writing them alone would
// trigger the next reconcile.
func (a NBGroupStatus) Equal(b NBGroupStatus) bool {
	return ptr.Equal(a.GroupID, b.GroupID) &&
		util.Equivalent(withoutConditionTimes(a.Conditions), withoutConditionTimes(b.Conditions))
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NBGroup is the Schema for the nbgroups API.
type NBGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NBGroupSpec   `json:"spec,omitempty"`
	Status NBGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NBGroupList contains a list of NBGroup.
type NBGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NBGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NBGroup{}, &NBGroupList{})
}
