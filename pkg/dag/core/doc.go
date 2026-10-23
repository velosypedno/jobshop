// Package core defines domain entities and topology representation for DAG-based
// Job Shop scheduling.
//
// # Key Principles & Immutability
//
//   - Templates are populated by the consumer: [JobGraphTemplate] instances and
//     their operations/edges are constructed manually (e.g., via parsers or code).
//   - Runtime entities are read-only: [Job] and [Operation] represent instantiated
//     work orders. They must not be created or mutated manually.
//   - Factory workflow: Templates and machine configs are fed into [Manufactory],
//     which validates topologies and materializes an immutable, read-only [Problem]
//     from incoming orders.
//
// # Reference
//
// ## Type Aliases
//
//	type (
//		MachineType   uint
//		MachineID     int
//		OpID          int
//		JobID         int
//		OpTemplateID  int
//		JobTemplateID int
//	)
//
// ## Structs
//
//	// Templates
//	OpTemplate
//	OpTemplateDependencyEdge
//	JobGraphTemplate
//
//	// Configuration & Machines
//	Machine
//	MachineConfigEntry
//
//	// Factory & Orders
//	OrderEntry
//	Manufactory
//
//	// Runtime Entities (Read-Only)
//	Operation
//	Job
//	Problem
//
//	// Error Details
//	ErrUnknownOpTemplateDetail
//	ErrMissingMachineTypeDetail
//	ErrUnknownJobTemplateDetail
//
// ## Main Functions & Methods
//
//	// Templates
//	InitJobTemplate(name string) JobGraphTemplate
//	(*JobGraphTemplate).AddOpTemplate(opTemplate OpTemplate)
//	(*JobGraphTemplate).AddEdge(from, to OpTemplateID)
//
//	// Factory Construction
//	NewMachine(id MachineID, machineType MachineType, machineName string) Machine
//	NewManufactory(machineConfigs []MachineConfigEntry, jobTemplates []JobGraphTemplate) (*Manufactory, error)
//
//	// Problem Generation
//	(*Manufactory).NewProblem(orders []OrderEntry) (*Problem, error)
package core
