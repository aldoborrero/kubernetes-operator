// SPDX-License-Identifier: BSD-3-Clause

package v1

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNBResourceStatusEqualComparesPointersByValue(t *testing.T) {
	t.Parallel()

	resource := "resource"
	policy := "policy"
	status := NBResourceStatus{
		NetworkResourceID: &resource,
		PolicyName:        &policy,
		TCPPorts:          []int32{80, 443},
		Groups:            []string{"a"},
		Conditions:        NBConditionTrue(),
	}

	// The reconciler compares a deep copy with the status it rebuilt, so equal values
	// behind different pointers must match.
	require.True(t, status.Equal(*status.DeepCopy()), "status with equal values should be equal")
	require.True(t, NBResourceStatus{}.Equal(NBResourceStatus{}), "empty statuses should be equal")

	changed := *status.DeepCopy()
	other := "other-resource"
	changed.NetworkResourceID = &other
	require.False(t, status.Equal(changed), "different network resource ID should not be equal")

	changed = *status.DeepCopy()
	changed.PolicyName = nil
	require.False(t, status.Equal(changed), "missing policy name should not be equal")
}

func TestNBResourceStatusEqualIgnoresConditionTimes(t *testing.T) {
	t.Parallel()

	status := NBResourceStatus{
		Conditions: NBConditionTrue(),
	}

	// The reconciler stamps fresh condition times on every pass.
	restamped := *status.DeepCopy()
	restamped.Conditions[0].LastProbeTime = metav1.NewTime(time.Now().Add(time.Hour))
	restamped.Conditions[0].LastTransitionTime = metav1.NewTime(time.Now().Add(time.Hour))
	require.True(t, status.Equal(restamped), "status differing only in condition times should be equal")

	changed := *status.DeepCopy()
	changed.Conditions = NBConditionFalse("internalError", "boom")
	require.False(t, status.Equal(changed), "different condition status should not be equal")
}
