// SPDX-License-Identifier: BSD-3-Clause

package v1

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNBGroupStatusEqualComparesPointersByValue(t *testing.T) {
	t.Parallel()

	group := "group"
	status := NBGroupStatus{
		GroupID:    &group,
		Conditions: NBConditionTrue(),
	}

	require.True(t, status.Equal(*status.DeepCopy()), "status with equal values should be equal")
	require.True(t, NBGroupStatus{}.Equal(NBGroupStatus{}), "status without pointers set should be equal")

	changed := *status.DeepCopy()
	other := "other-group"
	changed.GroupID = &other
	require.False(t, status.Equal(changed), "different group ID should not be equal")

	changed = *status.DeepCopy()
	changed.GroupID = nil
	require.False(t, status.Equal(changed), "missing group ID should not be equal")
}

func TestNBGroupStatusEqualIgnoresConditionTimes(t *testing.T) {
	t.Parallel()

	status := NBGroupStatus{
		Conditions: NBConditionTrue(),
	}

	// The reconciler stamps fresh condition times on every pass.
	restamped := *status.DeepCopy()
	restamped.Conditions[0].LastProbeTime = metav1.NewTime(time.Now().Add(time.Hour))
	restamped.Conditions[0].LastTransitionTime = metav1.NewTime(time.Now().Add(time.Hour))
	require.True(t, status.Equal(restamped), "status differing only in condition times should be equal")

	changed := *status.DeepCopy()
	changed.Conditions = NBConditionFalse("APIError", "boom")
	require.False(t, status.Equal(changed), "different condition status should not be equal")
}
