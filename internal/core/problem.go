package core

import "time"

type Problem struct {
	Jobs      []*Job
	Machines  []*Machine
	StartTime time.Time
}

type ProblemContext struct {
	Problem *Problem

	jobsByID            map[JobID]*Job
	operationsByID      map[OperationID]*Operation
	jobIDByOperationId  map[OperationID]JobID
	operationsIDByJobID map[JobID][]OperationID
	machineByID         map[MachineID]*Machine
}

func NewProblemContext(problem *Problem) *ProblemContext {
	ctx := &ProblemContext{
		Problem:             problem,
		jobsByID:            make(map[JobID]*Job),
		operationsByID:      make(map[OperationID]*Operation),
		jobIDByOperationId:  make(map[OperationID]JobID),
		operationsIDByJobID: make(map[JobID][]OperationID),
		machineByID:         make(map[MachineID]*Machine),
	}

	for _, m := range problem.Machines {
		ctx.machineByID[m.ID] = m
	}

	for _, job := range problem.Jobs {
		ctx.jobsByID[job.ID] = job

		for _, rootOp := range job.Operations {
			ctx.indexOperation(job.ID, rootOp)
		}
	}

	return ctx
}

func (ctx *ProblemContext) indexOperation(jobID JobID, op *Operation) {
	if op == nil {
		return
	}

	ctx.operationsByID[op.ID] = op
	ctx.jobIDByOperationId[op.ID] = jobID
	ctx.operationsIDByJobID[jobID] = append(ctx.operationsIDByJobID[jobID], op.ID)

	for _, child := range op.ChildOperations {
		if _, exists := ctx.operationsByID[child.ID]; !exists {
			ctx.indexOperation(jobID, child)
		}
	}
}

func (ctx *ProblemContext) GetOperation(id OperationID) (*Operation, bool) {
	op, ok := ctx.operationsByID[id]
	return op, ok
}

func (ctx *ProblemContext) GetJob(id JobID) (*Job, bool) {
	job, ok := ctx.jobsByID[id]
	return job, ok
}

func (ctx *ProblemContext) GetJobIDByOperation(id OperationID) (JobID, bool) {
	jobID, ok := ctx.jobIDByOperationId[id]
	return jobID, ok
}

func (ctx *ProblemContext) GetOperationsByJob(id JobID) []OperationID {
	return ctx.operationsIDByJobID[id]
}

func (ctx *ProblemContext) GetMachine(id MachineID) (*Machine, bool) {
	m, ok := ctx.machineByID[id]
	return m, ok
}
