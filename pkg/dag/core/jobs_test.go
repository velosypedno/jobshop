package core

import (
	"errors"
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
	t := InitJobTemplate("linear")
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
		tmpl := InitJobTemplate(name)

		// Assert
		assert.Equal(t, name, tmpl.Name)
	})

	main.Run("OpsMapInitialized", func(t *testing.T) {
		// Act
		tmpl := InitJobTemplate("x")

		// Assert
		assert.NotNil(t, tmpl.Ops)
	})

	main.Run("EdgesSliceInitialized", func(t *testing.T) {
		// Act
		tmpl := InitJobTemplate("x")

		// Assert
		assert.NotNil(t, tmpl.Edges)
	})

	main.Run("EmptyName", func(t *testing.T) {
		// Act
		tmpl := InitJobTemplate("")

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
		tmpl := InitJobTemplate("t")
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
		tmpl := InitJobTemplate("t")
		tmpl.AddOpTemplate(OpTemplate{ID: 1, Name: "first"})

		// Act
		tmpl.AddOpTemplate(OpTemplate{ID: 1, Name: "second"})

		// Assert
		assert.Equal(t, "second", tmpl.Ops[1].Name)
	})

	main.Run("MultipleOps", func(t *testing.T) {
		// Arrange
		tmpl := InitJobTemplate("t")

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
		tmpl := InitJobTemplate("t")

		// Act
		tmpl.AddEdge(1, 2)

		// Assert
		require.Len(t, tmpl.Edges, 1)
		assert.Equal(t, OpTemplateID(1), tmpl.Edges[0].From)
		assert.Equal(t, OpTemplateID(2), tmpl.Edges[0].To)
	})

	main.Run("MultipleEdges", func(t *testing.T) {
		// Arrange
		tmpl := InitJobTemplate("t")

		// Act
		tmpl.AddEdge(1, 2)
		tmpl.AddEdge(2, 3)
		tmpl.AddEdge(1, 3)

		// Assert
		assert.Len(t, tmpl.Edges, 3)
	})

	main.Run("SelfLoopStored", func(t *testing.T) {
		// Arrange
		tmpl := InitJobTemplate("t")

		// Act — AddEdge itself does not validate; self-loop should be stored
		tmpl.AddEdge(1, 1)

		// Assert
		assert.Len(t, tmpl.Edges, 1)
	})

	main.Run("DuplicateEdgesStored", func(t *testing.T) {
		// Arrange
		tmpl := InitJobTemplate("t")

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
			tmpl := InitJobTemplate("empty")

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("SingleOpNoEdges", func(t *testing.T) {
			// Arrange
			tmpl := InitJobTemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("LinearChain", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(4)

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("DiamondDAG", func(t *testing.T) {
			// Arrange — diamond: 1 → {2,3} → 4
			tmpl := InitJobTemplate("diamond")
			for i := OpTemplateID(1); i <= 4; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(1, 3)
			tmpl.AddEdge(2, 4)
			tmpl.AddEdge(3, 4)

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("MultipleRoots", func(t *testing.T) {
			// Arrange — two sources converge: 1 → 3, 2 → 3
			tmpl := InitJobTemplate("multi-root")
			for i := OpTemplateID(1); i <= 3; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			tmpl.AddEdge(1, 3)
			tmpl.AddEdge(2, 3)

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			assert.NoError(t, err)
		})

		t.Run("LongChain100Nodes", func(t *testing.T) {
			// Arrange
			tmpl := buildLinearTemplate(100)

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			assert.NoError(t, err)
		})
	})

	main.Run("InvalidEdges", func(t *testing.T) {
		t.Run("MissingFromNode", func(t *testing.T) {
			// Arrange
			tmpl := InitJobTemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 2})
			tmpl.AddEdge(1, 2) // ID 1 does not exist

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			require.Error(t, err)
		})

		t.Run("MissingToNode", func(t *testing.T) {
			// Arrange
			tmpl := InitJobTemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddEdge(1, 99) // ID 99 does not exist

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			require.Error(t, err)
		})

		t.Run("BothNodesMissing", func(t *testing.T) {
			// Arrange
			tmpl := InitJobTemplate("t")
			tmpl.AddEdge(10, 20) // neither 10 nor 20 exists

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			require.Error(t, err)
		})
	})

	main.Run("Cycles", func(t *testing.T) {
		t.Run("SelfLoop", func(t *testing.T) {
			// Arrange
			tmpl := InitJobTemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddEdge(1, 1)

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrCycleInGraph))
		})

		t.Run("TwoNodeCycle", func(t *testing.T) {
			// Arrange
			tmpl := InitJobTemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddOpTemplate(OpTemplate{ID: 2})
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(2, 1)

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrCycleInGraph))
		})

		t.Run("ThreeNodeCycle", func(t *testing.T) {
			// Arrange — 1 → 2 → 3 → 1
			tmpl := InitJobTemplate("t")
			for i := OpTemplateID(1); i <= 3; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(2, 3)
			tmpl.AddEdge(3, 1)

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrCycleInGraph))
		})

		t.Run("CycleWithExtraNodes", func(t *testing.T) {
			// Arrange — node 4 is a valid source; nodes 1→2→3→1 form a cycle
			tmpl := InitJobTemplate("t")
			for i := OpTemplateID(1); i <= 4; i++ {
				tmpl.AddOpTemplate(OpTemplate{ID: i})
			}
			tmpl.AddEdge(4, 1)
			tmpl.AddEdge(1, 2)
			tmpl.AddEdge(2, 3)
			tmpl.AddEdge(3, 1)

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert
			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrCycleInGraph))
		})
	})

	main.Run("ErrorMessages", func(t *testing.T) {
		t.Run("MissingToNodeCarriesCorrectID", func(t *testing.T) {
			// Arrange
			tmpl := InitJobTemplate("t")
			tmpl.AddOpTemplate(OpTemplate{ID: 1})
			tmpl.AddEdge(1, 99) // 99 is missing

			// Act
			err := validateOperationsTopology(tmpl)

			// Assert — typed error carries the correct missing ID (99)
			require.Error(t, err)
			assert.True(t, errors.Is(err, ErrUnknownOpTemplate))
			detail, ok := AsUnknownOpTemplate(err)
			require.True(t, ok)
			assert.Equal(t, OpTemplateID(99), detail.ID)
		})
	})
}
