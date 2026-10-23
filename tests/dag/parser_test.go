package dag_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/velosypedno/jobshop/internal/dag/parser"
	"github.com/velosypedno/jobshop/pkg/dag/core"
)

func TestParseFactoryConfig(main *testing.T) {
	main.Run("SequentialJobs", func(t *testing.T) {
		// Arrange
		fixturePath := filepath.Join("fixtures", "sequential_jobs.json")

		// Act
		machines, templates, err := parser.ParseFactoryConfig(fixturePath)

		// Assert
		require.NoError(t, err)

		t.Run("MachinesParsedCorrectly", func(t *testing.T) {
			// Assert
			require.Len(t, machines, 3)
			assert.Equal(t, core.MachineConfigEntry{Type: 1, Name: "Milling", Count: 2}, machines[0])
			assert.Equal(t, core.MachineConfigEntry{Type: 2, Name: "Lathe", Count: 1}, machines[1])
			assert.Equal(t, core.MachineConfigEntry{Type: 3, Name: "Polisher", Count: 1}, machines[2])
		})

		t.Run("JobTemplateAndLinearEdgesParsedCorrectly", func(t *testing.T) {
			// Assert
			require.Len(t, templates, 1)
			tmpl := templates[0]

			assert.Equal(t, core.JobTemplateID(101), tmpl.ID)
			assert.Equal(t, "Sequential Shaft Processing", tmpl.Name)
			require.Len(t, tmpl.Ops, 3)

			assert.Equal(t, core.OpTemplate{
				ID:             1,
				Name:           "Rough Turn",
				MachineType:    2,
				ProcessingTime: 30 * time.Minute,
			}, tmpl.Ops[1])

			assert.Equal(t, core.OpTemplate{
				ID:             2,
				Name:           "Keyway Milling",
				MachineType:    1,
				ProcessingTime: 15 * time.Minute,
			}, tmpl.Ops[2])

			assert.Equal(t, core.OpTemplate{
				ID:             3,
				Name:           "Surface Polishing",
				MachineType:    3,
				ProcessingTime: 20 * time.Minute,
			}, tmpl.Ops[3])

			require.Len(t, tmpl.Edges, 2)
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 1, To: 2})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 2, To: 3})
		})
	})

	main.Run("TreeJobs", func(t *testing.T) {
		// Arrange
		fixturePath := filepath.Join("fixtures", "tree_jobs.json")

		// Act
		machines, templates, err := parser.ParseFactoryConfig(fixturePath)

		// Assert
		require.NoError(t, err)

		t.Run("MachinesParsedCorrectly", func(t *testing.T) {
			// Assert
			require.Len(t, machines, 3)
			assert.Equal(t, core.MachineConfigEntry{Type: 1, Name: "Stamping", Count: 3}, machines[0])
			assert.Equal(t, core.MachineConfigEntry{Type: 2, Name: "Welding", Count: 2}, machines[1])
			assert.Equal(t, core.MachineConfigEntry{Type: 3, Name: "Assembly", Count: 1}, machines[2])
		})

		t.Run("JobTemplateAndTreeEdgesParsedCorrectly", func(t *testing.T) {
			// Assert
			require.Len(t, templates, 1)
			tmpl := templates[0]

			assert.Equal(t, core.JobTemplateID(201), tmpl.ID)
			assert.Equal(t, "Tree Assembly Gearbox", tmpl.Name)
			require.Len(t, tmpl.Ops, 5)

			assert.Equal(t, core.OpTemplate{
				ID:             1,
				Name:           "Stamp Gear A",
				MachineType:    1,
				ProcessingTime: 10 * time.Minute,
			}, tmpl.Ops[1])

			assert.Equal(t, core.OpTemplate{
				ID:             2,
				Name:           "Stamp Gear B",
				MachineType:    1,
				ProcessingTime: 10 * time.Minute,
			}, tmpl.Ops[2])

			assert.Equal(t, core.OpTemplate{
				ID:             3,
				Name:           "Weld Sub-Assembly",
				MachineType:    2,
				ProcessingTime: 25 * time.Minute,
			}, tmpl.Ops[3])

			assert.Equal(t, core.OpTemplate{
				ID:             4,
				Name:           "Stamp Casing",
				MachineType:    1,
				ProcessingTime: 15 * time.Minute,
			}, tmpl.Ops[4])

			assert.Equal(t, core.OpTemplate{
				ID:             5,
				Name:           "Final Assembly",
				MachineType:    3,
				ProcessingTime: 45 * time.Minute,
			}, tmpl.Ops[5])

			require.Len(t, tmpl.Edges, 4)
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 1, To: 3})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 2, To: 3})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 3, To: 5})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 4, To: 5})
		})
	})

	main.Run("GraphJobs", func(t *testing.T) {
		// Arrange
		fixturePath := filepath.Join("fixtures", "graph_jobs.json")

		// Act
		machines, templates, err := parser.ParseFactoryConfig(fixturePath)

		// Assert
		require.NoError(t, err)

		t.Run("MachinesParsedCorrectly", func(t *testing.T) {
			// Assert
			require.Len(t, machines, 4)
			assert.Equal(t, core.MachineConfigEntry{Type: 10, Name: "LaserCutter", Count: 2}, machines[0])
			assert.Equal(t, core.MachineConfigEntry{Type: 20, Name: "BendingMachine", Count: 1}, machines[1])
			assert.Equal(t, core.MachineConfigEntry{Type: 30, Name: "HeatTreatment", Count: 1}, machines[2])
			assert.Equal(t, core.MachineConfigEntry{Type: 40, Name: "QualityControl", Count: 2}, machines[3])
		})

		t.Run("JobTemplateAndGraphEdgesParsedCorrectly", func(t *testing.T) {
			// Assert
			require.Len(t, templates, 1)
			tmpl := templates[0]

			assert.Equal(t, core.JobTemplateID(301), tmpl.ID)
			assert.Equal(t, "Complex Graph Bracket", tmpl.Name)
			require.Len(t, tmpl.Ops, 5)

			assert.Equal(t, core.OpTemplate{
				ID:             1,
				Name:           "Cut Sheet Plate",
				MachineType:    10,
				ProcessingTime: 12 * time.Minute,
			}, tmpl.Ops[1])

			assert.Equal(t, core.OpTemplate{
				ID:             2,
				Name:           "Form Flange A",
				MachineType:    20,
				ProcessingTime: 18 * time.Minute,
			}, tmpl.Ops[2])

			assert.Equal(t, core.OpTemplate{
				ID:             3,
				Name:           "Form Flange B",
				MachineType:    20,
				ProcessingTime: 18 * time.Minute,
			}, tmpl.Ops[3])

			assert.Equal(t, core.OpTemplate{
				ID:             4,
				Name:           "Harden Joints",
				MachineType:    30,
				ProcessingTime: 60 * time.Minute,
			}, tmpl.Ops[4])

			assert.Equal(t, core.OpTemplate{
				ID:             5,
				Name:           "Coordinate Metrology",
				MachineType:    40,
				ProcessingTime: 30 * time.Minute,
			}, tmpl.Ops[5])

			require.Len(t, tmpl.Edges, 6)
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 1, To: 2})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 1, To: 3})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 2, To: 4})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 3, To: 4})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 4, To: 5})
			assert.Contains(t, tmpl.Edges, core.OpTemplateDependencyEdge{From: 2, To: 5})
		})
	})
}
