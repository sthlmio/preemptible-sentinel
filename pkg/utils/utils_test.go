package utils

import (
	types "k8s.io/api/core/v1"
	"testing"
)

func TestIsNodeReady(t *testing.T) {
	conditions := []types.NodeCondition{
		{Type: types.NodeReady, Status: types.ConditionTrue},
		{Type: types.NodeDiskPressure, Status: types.ConditionTrue},
		{Type: types.NodeMemoryPressure, Status: types.ConditionUnknown},
	}

	if !IsNodeReady(types.NodeStatus{Conditions: conditions}) {
		t.Error("Node is not ready")
	}
}
