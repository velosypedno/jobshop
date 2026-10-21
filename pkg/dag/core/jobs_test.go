package core

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// buildLinearTemplate creates a chain: 1 → 2 → … → n
func buildLinearTemplate(n int) JobGraphTemplate {
	t := InitJobtemplate("linear")
	for i := 1; i <= n; i++ {
		t.AddOpTemplate(OpTemplate{
			ID:             OpTemplateID(i),
			Name:           "op",
			MachineType:    1,
			ProcessingTime: time.Second,
		})
	}
	for i := 1; i < n; i++ {
		t.AddEdge(OpTemplateID(i), OpTemplateID(i+1))
	}
	return t
}

// ─────────────────────────────────────────────────────────────────────────────
// TestInitJobtemplate
// ─────────────────────────────────────────────────────────────────────────────

func TestInitJobtemplate(main *testing.T) {
	main.Run("SetsName", func(t *testing.T) {
		// Arrange
		name := "my-job"

		// Act
		tmpl := InitJobtemplate(name)

		// Assert
		assert.Equal(t, name, tmpl.Name)
	})

	main.Run("OpsMapInitialized", func(t *testing.T) {
		// Act
		tmpl := InitJobtemplate("x")

		// Assert
		assert.NotNil(t, tmpl.Ops)
	})

	main.Run("EdgesSliceInitialized", func(t *testing.T) {
		// Act
		tmpl := InitJobtemplate("x")

		// Assert
		assert.NotNil(t, tmpl.Edges)
	})

	main.Run("EmptyName", func(t *testing.T) {
		// Act
		tmpl := InitJobtemplate("")

		// Assert
		assert.Equal(t, "", tmpl.Name)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// TestAddOpTemplate
// ─────────────────────────────────────────────────────────────────────────────

func TestAddOpTemplate(main *testing.T) {
	main.Run("StoresOp", func(t *testing.T) {
		// Arrange
		tmpl := InitJobtemplate("t")
		op := OpTemplate{ID: 1, Name: "cut", MachineType: 2, ProcessingTime: 5 * time.Second}

		// Act
		tmpl.AddOpTemplate(op)

		// Assert
		got, ok := tmpl.Ops[1]
		require.True(t, ok, "op with ID 1 should be in Ops map")
		assert.Equal(t, op, got)
	})

	main.Run("OverwritesSameID", func(t *testing.T) {
		// Arrange
		tmpl := InitJobtemplate("t")
		tmpl.AddOpTemplate(OpTemplate{ID: 1, Name: "first"})

		// Act
		tmpl.AddOpTemplate(OpTemplate{ID: 1, Name: "second"})

		// Assert
		assert.Equal(t, "second", tmpl.Ops[1].Name)
	})

	main.Run("MultipleOps", func(t *testing.T) {
		// Arrange
		tmpl := InitJobtemplate("t")

		// Act
		for i := OpTemplateID(1); i <= 5; i++ {
			tmpl.AddOpTemplate(OpTemplate{ID: i})
		}

		// Assert
		assert.Len(t, tmpl.Ops, 5)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// TestAddEdge
// ─────────────────────────────────────────────────────────────────────────────

func TestAddEdge(main *testing.T) {
	main.Run("AddsEdge", func(t *testing.T) {
		// Arrange
		tmpl := InitJobtemplate("t")

		// Act
		tmpl.AddEdge(1, 2)

		// Assert
		require.Len(t, tmpl.Edges, 1)
		assert.Equal(t, OpTemplateID(1), tmpl.Edges[0].From)
		assert.Equal(t, OpTemplateID(2), tmpl.Edges[0].To)
	})

	main.Run("MultipleEdges", func(t *testing.T) {
		// Arrange
		tmpl := InitJobtemplate("t")

		// Act
		tmpl.AddEdge(1, 2)
		tmpl.AddEdge(2, 3)
		tmpl.AddEdge(1, 3)

		// Assert
		assert.Len(t, tmpl.Edges, 3)
	})

	main.Run("SelfLoopStored", func(t *testing.T) {
		// Arrange
		tmpl := InitJobtemplate("t")

		// Act — AddEdge itself does not validate; self-loop should be stored
		tmpl.AddEdge(1, 1)

		// Assert
		assert.Len(t, tmpl.Edges, 1)
	})

	main.Run("DuplicateEdgesStored", func(t *testing.T) {
		// Arrange
		tmpl := InitJobtemplate("t")

		// Act — AddEdge is append-only; duplicates are allowed
		tmpl.AddEdge(1, 2)
		tmpl.AddEdge(1, 2)

		// Assert
		assert.Len(t, tmpl.Edges, 2)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// TestValidateTemplate
// ─────────────────────────────────────────────────────────────────────────────

func TestValidateTemplate(main *testing.T) {
	main.Run("Valid", func(t *testing.T) {
		t.Run("EmptyTemplate", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("empty")

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("SingleOpNoEdges", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("LinearChain", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(4)

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("DiamondDAG", func(t *testing.T) {
			// Arrange — diamond: 1 → {2,3} → 4
			tmpl := InitJobtemplate("diamond")
			for i := OpTemplateID(1); i <= 4; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(1, 3)
			tmpl.AddEdge(2, 4)
			tmpl.AddEdge(3, 4)

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("MultipleRoots", func(t *testing.T) {
			// Arrange — two sources converge: 1 → 3, 2 → 3
			tmpl := InitJobtemplate("multi-root")
			for i := OpTemplateID(1); i <= 3; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			tmpl.AddEdge(1, 3)
			tmpl.AddEdge(2, 3)

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("LongChain100Nodes", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(100)

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			assert.NoError(t, err)
		})
	})

	main.Run("InvalidEdges", func(t *testing.T) {
		t.Run("MissingFromNode", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 2})
			tmpl.AddEdge(1, 2) // ID 1 does not exist

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			require.Error(t, err)
		})

		t.Run("MissingToNode", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddEdge(1, 99) // ID 99 does not exist

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			require.Error(t, err)
		})

		t.Run("BothNodesMissing", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("t")
			tmpl.AddEdge(10, 20) // neither 10 nor 20 exists

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			require.Error(t, err)
		})
	})

	main.Run("Cycles", func(t *testing.T) {
		t.Run("SelfLoop", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddEdge(1, 1)

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			require.Error(t, err)
			assert.Equal(t, "cycle presents in topology", err.Error())
		})

		t.Run("TwoNodeCycle", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddOpTemplate(OpTemplate{ID: 2})
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(2, 1)

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			require.Error(t, err)
			assert.Equal(t, "cycle presents in topology", err.Error())
		})

		t.Run("ThreeNodeCycle", func(t *testing.T) {
			// Arrange — 1 → 2 → 3 → 1
			tmpl := InitJobtemplate("t")
			for i := OpTemplateID(1); i <= 3; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(2, 3)
			tmpl.AddEdge(3, 1)

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			require.Error(t, err)
			assert.Equal(t, "cycle presents in topology", err.Error())
		})

		t.Run("CycleWithExtraNodes", func(t *testing.T) {
			// Arrange — node 4 is a valid source; nodes 1→2→3→1 form a cycle
			tmpl := InitJobtemplate("t")
			for i := OpTemplateID(1); i <= 4; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			tmpl.AddEdge(4, 1)
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(2, 3)
			tmpl.AddEdge(3, 1)

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			require.Error(t, err)
		})
	})

	main.Run("ErrorMessages", func(t *testing.T) {
		t.Run("MissingToNodeReportsWrongID", func(t *testing.T) {
			// Arrange
			// NOTE: known bug — when the "To" node is missing the error
			// message reports edge.From instead of edge.To (line 58 of jobs.go).
			tmpl := InitJobtemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddEdge(1, 99) // 99 is missing

			// Act
			err := ValidateTemplate(tmpl)

			// Assert
			require.Error(t, err)
			// Bug: message says "1 does not exist" instead of "99 does not exist"
			assert.Contains(t, err.Error(), "1")
			assert.NotContains(t, err.Error(), "99")
		})
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// TestCreateJob
// ─────────────────────────────────────────────────────────────────────────────

func TestCreateJob(main *testing.T) {
	main.Run("Success", func(t *testing.T) {
		t.Run("ReturnsNonNilJob", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(3)
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, job)
		})

		t.Run("JobIDSet", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(1)
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(42), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, JobID(42), job.ID)
		})

		t.Run("JobNameFromTemplate", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("assembly")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, "assembly", job.Name)
		})

		t.Run("CorrectNumberOfOps", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(5)
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			assert.Len(t, job.Ops, 5)
		})

		t.Run("EmptyTemplate", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("empty")
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			assert.Empty(t, job.Ops)
			assert.Equal(t, OpID(0), lastID)
		})

		t.Run("SingleOp", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("single")
			tmpl.AddOpTemplate(OpTemplate{ID: 1, Name: "solo", ProcessingTime: time.Minute})
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(5), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			require.Len(t, job.Ops, 1)
			op := job.Ops[OpID(1)]
			assert.Equal(t, "solo", op.Name)
			assert.Empty(t, op.Predecessors)
			assert.Empty(t, op.Successors)
		})
	})

	main.Run("InvalidTemplate", func(t *testing.T) {
		t.Run("CyclicTemplateReturnsError", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("bad")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddOpTemplate(OpTemplate{ID: 2})
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(2, 1)
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.Error(t, err)
			assert.Nil(t, job)
		})

		t.Run("ErrorHasCreateJobPrefix", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("bad")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddEdge(1, 1)
			lastID := OpID(0)

			// Act
			_, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), "create job:")
		})

		t.Run("LastIDUnchangedOnError", func(t *testing.T) {
			// Arrange — validation fails before lastID is touched
			tmpl := InitJobtemplate("bad")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddEdge(1, 1)
			lastID := OpID(7)

			// Act
			CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			assert.Equal(t, OpID(7), lastID)
		})
	})

	main.Run("OpCounter", func(t *testing.T) {
		t.Run("AdvancesFromZero", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(3)
			lastID := OpID(0)

			// Act
			CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			assert.Equal(t, OpID(3), lastID)
		})

		t.Run("AdvancesFromNonZero", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(2)
			lastID := OpID(10)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, OpID(12), lastID)
			_, has11 := job.Ops[OpID(11)]
			_, has12 := job.Ops[OpID(12)]
			assert.True(t, has11, "expected op with ID 11")
			assert.True(t, has12, "expected op with ID 12")
		})

		t.Run("SharedAcrossMultipleJobs", func(t *testing.T) {
			// Arrange
			tmpl1 := buildLinearTemplate(2)
			tmpl2 := buildLinearTemplate(3)
			lastID := OpID(0)

			// Act
			CreateJob(JobID(1), tmpl1, &lastID)
			CreateJob(JobID(2), tmpl2, &lastID)

			// Assert — 2 + 3 = 5 total ops
			assert.Equal(t, OpID(5), lastID)
		})
	})

	main.Run("OpFieldsCopied", func(t *testing.T) {
		t.Run("AllFieldsPreserved", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("t")
			tmpl.AddOpTemplate(OpTemplate{
				ID:             1,
				Name:           "weld",
				MachineType:    MachineType(3),
				ProcessingTime: 7 * time.Second,
			})
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			op := job.Ops[OpID(1)]
			assert.Equal(t, "weld", op.Name)
			assert.Equal(t, MachineType(3), op.MachineType)
			assert.Equal(t, 7*time.Second, op.ProcessingTime)
		})
	})

	main.Run("OpWiring", func(t *testing.T) {
		t.Run("LinearChainHasOneRootAndOneLeaf", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(3)
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			var roots, leaves, middles int
			for _, op := range job.Ops {
				switch {
				case len(op.Predecessors) == 0 && len(op.Successors) > 0:
					roots++
				case len(op.Predecessors) > 0 && len(op.Successors) == 0:
					leaves++
				default:
					middles++
				}
			}
			assert.Equal(t, 1, roots, "expected exactly 1 root")
			assert.Equal(t, 1, leaves, "expected exactly 1 leaf")
			assert.Equal(t, 1, middles, "expected exactly 1 middle node")
		})

		t.Run("SuccessorPredecessorConsistency", func(t *testing.T) {
			// Arrange — for every op A with successor B, B must list A as predecessor
			tmpl := buildLinearTemplate(4)
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			for _, op := range job.Ops {
				for _, sucID := range op.Successors {
					suc := job.Ops[sucID]
					found := slices.Contains(suc.Predecessors, op.ID)
					assert.True(t, found, "op %d is successor of %d but %d is not in its Predecessors", sucID, op.ID, op.ID)
				}
			}
		})

		t.Run("DiamondDAGWiring", func(t *testing.T) {
			// Arrange — diamond: 1 → {2,3} → 4
			tmpl := InitJobtemplate("diamond")
			for i := OpTemplateID(1); i <= 4; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i, ProcessingTime: time.Second})
			}
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(1, 3)
			tmpl.AddEdge(2, 4)
			tmpl.AddEdge(3, 4)
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			var root, sink *Operation
			for _, op := range job.Ops {
				if len(op.Predecessors) == 0 {
					root = op
				}
				if len(op.Successors) == 0 {
					sink = op
				}
			}
			require.NotNil(t, root)
			require.NotNil(t, sink)
			assert.Len(t, root.Successors, 2)
			assert.Len(t, sink.Predecessors, 2)
		})

		t.Run("NoEdgesAllOpsIsolated", func(t *testing.T) {
			// Arrange
			tmpl := InitJobtemplate("isolated")
			for i := OpTemplateID(1); i <= 3; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			for _, op := range job.Ops {
				assert.Empty(t, op.Predecessors, "op %d should have no predecessors", op.ID)
				assert.Empty(t, op.Successors, "op %d should have no successors", op.ID)
			}
		})
	})

	main.Run("EdgeCases", func(t *testing.T) {
		t.Run("LargeLinearChain200Nodes", func(t *testing.T) {
			// Arrange
			const n = 200
			tmpl := buildLinearTemplate(n)
			lastID := OpID(0)

			// Act
			job, err := CreateJob(JobID(1), tmpl, &lastID)

			// Assert
			require.NoError(t, err)
			assert.Len(t, job.Ops, n)
			assert.Equal(t, OpID(n), lastID)
		})
	})
}
