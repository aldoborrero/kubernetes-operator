// SPDX-License-Identifier: BSD-3-Clause

package v1

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNBRoutingPeerStatusEqualComparesPointersByValue(t *testing.T) {
	t.Parallel()

	network := "network"
	setupKey := "setup-key"
	router := "router"
	status := NBRoutingPeerStatus{
		NetworkID:  &network,
		SetupKeyID: &setupKey,
		RouterID:   &router,
		Conditions: NBConditionTrue(),
	}

	// The reconciler compares a deep copy with the status it rebuilt, so equal values
	// behind different pointers must match.
	require.True(t, status.Equal(*status.DeepCopy()), "status with equal values should be equal")
	require.True(t, NBRoutingPeerStatus{}.Equal(NBRoutingPeerStatus{}), "status without pointers set should be equal")

	changed := *status.DeepCopy()
	other := "other-network"
	changed.NetworkID = &other
	require.False(t, status.Equal(changed), "different network ID should not be equal")

	changed = *status.DeepCopy()
	changed.SetupKeyID = nil
	require.False(t, status.Equal(changed), "missing setup key ID should not be equal")

	changed = *status.DeepCopy()
	other = "other-router"
	changed.RouterID = &other
	require.False(t, status.Equal(changed), "different router ID should not be equal")
}

func TestNBRoutingPeerStatusEqualIgnoresConditionTimes(t *testing.T) {
	t.Parallel()

	// No pointers are set, so only the conditions can tell the statuses apart.
	status := NBRoutingPeerStatus{
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
