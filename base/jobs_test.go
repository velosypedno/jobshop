package base_test

import (
	"testing"
	"time"

	"github.com/velosypedno/resource-allocation/base"
)

func TestCreateJob_Flattening(t *testing.T) {
	template := base.JobTemplate{
		Name: "Test Robot",
		Operations: []base.OperationTemplate{
			{
				Name:           "Final Assembly (Root)",
				MachineType:    1,
				ProcessingTime: 10 * time.Minute,
				Children: []base.OperationTemplate{
					{
						Name:           "Legs (Child 1)",
						MachineType:    2,
						ProcessingTime: 20 * time.Minute,
						Children: []base.OperationTemplate{
							{
								Name:           "Feet (GrandChild 1)",
								MachineType:    3,
								ProcessingTime: 5 * time.Minute,
							},
						},
					},
					{
						Name:           "Head (Child 2)",
						MachineType:    2,
						ProcessingTime: 15 * time.Minute,
					},
				},
			},
		},
	}

	job := base.CreateJob(1, template)

	if job.Name != "Test Robot" {
		t.Errorf("Expected name 'Test Robot', got '%s'", job.Name)
	}

	expectedCount := 4
	if count := job.OperationsCount(); count != expectedCount {
		t.Fatalf("Expected %d operations in flattened list, got %d", expectedCount, count)
	}

	expectedOrder := []string{
		"Feet (GrandChild 1)",
		"Legs (Child 1)",
		"Head (Child 2)",
		"Final Assembly (Root)",
	}

	for i, name := range expectedOrder {
		op := job.GetOperation(i)
		if op == nil {
			t.Fatalf("Operation at index %d is nil", i)
		}
		if op.Name != name {
			t.Errorf("At index %d: expected operation '%s', got '%s'", i, name, op.Name)
		}
	}

	for i := 0; i < job.OperationsCount(); i++ {
		op := job.GetOperation(i)
		for _, child := range op.ChildOperations {
			childIdx := -1
			for j := 0; j < job.OperationsCount(); j++ {
				if job.GetOperation(j).ID == child.ID {
					childIdx = j
					break
				}
			}

			if childIdx > i {
				t.Errorf("Constraint violation: Child '%s' (idx %d) is after parent '%s' (idx %d)",
					child.Name, childIdx, op.Name, i)
			}
		}
	}
}

func TestCreateJob_ComplexFlattening(t *testing.T) {
	template := base.JobTemplate{
		Name: "Complex Product",
		Operations: []base.OperationTemplate{
			{
				Name:           "Main System (Root A)",
				MachineType:    1,
				ProcessingTime: 60 * time.Minute,
				Children: []base.OperationTemplate{
					{
						Name:           "Sub-Assembly 1",
						MachineType:    2,
						ProcessingTime: 30 * time.Minute,
						Children: []base.OperationTemplate{
							{
								Name:           "Part 1.1",
								MachineType:    3,
								ProcessingTime: 10 * time.Minute,
							},
							{
								Name:           "Part 1.2",
								MachineType:    3,
								ProcessingTime: 12 * time.Minute,
								Children: []base.OperationTemplate{
									{
										Name:           "Raw Material 1.2.1",
										MachineType:    4,
										ProcessingTime: 5 * time.Minute,
									},
								},
							},
						},
					},
					{
						Name:           "Sub-Assembly 2",
						MachineType:    2,
						ProcessingTime: 20 * time.Minute,
						Children: []base.OperationTemplate{
							{
								Name:           "Part 2.1",
								MachineType:    3,
								ProcessingTime: 15 * time.Minute,
							},
						},
					},
				},
			},
			{
				Name:           "Accessory (Root B)",
				MachineType:    5,
				ProcessingTime: 10 * time.Minute,
			},
		},
	}

	job := base.CreateJob(7, template)

	expectedCount := 8
	if count := job.OperationsCount(); count != expectedCount {
		t.Fatalf("Expected %d operations, got %d", expectedCount, count)
	}

	allOpsMap := make(map[base.OperationID]bool)
	for i := 0; i < job.OperationsCount(); i++ {
		op := job.GetOperation(i)
		if allOpsMap[op.ID] {
			t.Errorf("Duplicate operation ID %d found at index %d", op.ID, i)
		}
		allOpsMap[op.ID] = true
	}

	treeIDs := make(map[base.OperationID]bool)
	var collectIDs func(*base.Operation)
	collectIDs = func(o *base.Operation) {
		treeIDs[o.ID] = true
		for _, child := range o.ChildOperations {
			collectIDs(child)
		}
	}
	for _, root := range job.Operations {
		collectIDs(root)
	}

	if len(allOpsMap) != len(treeIDs) {
		t.Errorf("Mismatch: Tree has %d operations, Flattened list has %d", len(treeIDs), len(allOpsMap))
	}

	for id := range treeIDs {
		if !allOpsMap[id] {
			t.Errorf("Operation ID %d from tree is missing in flattened list", id)
		}
	}
}
