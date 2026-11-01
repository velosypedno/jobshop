package core

import (
	"encoding/json"
	"fmt"
	"time"
)

type problemDTO struct {
	Jobs     []jobDTO     `json:"jobs"`
	Machines []machineDTO `json:"machines"`
}

type jobDTO struct {
	ID   int                     `json:"id"`
	Name string                  `json:"name"`
	Ops  map[string]operationDTO `json:"ops"`
}

type operationDTO struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ProcessingTime string `json:"processing_time"`
	MachineType    uint   `json:"machine_type"`
	Predecessors   []int  `json:"predecessors"`
	Successors     []int  `json:"successors"`
}

type machineDTO struct {
	ID   int    `json:"id"`
	Type uint   `json:"type"`
	Name string `json:"name"`
}

func (p *Problem) ToJSON() ([]byte, error) {
	dto := problemDTO{
		Jobs:     make([]jobDTO, len(p.Jobs)),
		Machines: make([]machineDTO, len(p.Machines)),
	}

	for i, job := range p.Jobs {
		opsMap := make(map[string]operationDTO, len(job.Ops))
		for opID, op := range job.Ops {
			preds := make([]int, len(op.Predecessors))
			for j, id := range op.Predecessors {
				preds[j] = int(id)
			}

			succs := make([]int, len(op.Successors))
			for j, id := range op.Successors {
				succs[j] = int(id)
			}

			opsMap[fmt.Sprintf("%d", opID)] = operationDTO{
				ID:             int(op.ID),
				Name:           op.Name,
				ProcessingTime: op.ProcessingTime.String(),
				MachineType:    uint(op.MachineType),
				Predecessors:   preds,
				Successors:     succs,
			}
		}

		dto.Jobs[i] = jobDTO{
			ID:   int(job.ID),
			Name: job.Name,
			Ops:  opsMap,
		}
	}

	for i, m := range p.Machines {
		dto.Machines[i] = machineDTO{
			ID:   int(m.ID),
			Type: uint(m.Type),
			Name: m.Name,
		}
	}

	return json.Marshal(dto)
}

func ProblemFromJSON(data []byte) (Problem, error) {
	var dto problemDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return Problem{}, err
	}

	jobs := make([]*Job, len(dto.Jobs))
	for i, jobDTO := range dto.Jobs {
		ops := make(map[OpID]*Operation, len(jobDTO.Ops))
		for _, opDTO := range jobDTO.Ops {
			dur, err := time.ParseDuration(opDTO.ProcessingTime)
			if err != nil {
				return Problem{}, err
			}

			preds := make([]OpID, len(opDTO.Predecessors))
			for j, id := range opDTO.Predecessors {
				preds[j] = OpID(id)
			}

			succs := make([]OpID, len(opDTO.Successors))
			for j, id := range opDTO.Successors {
				succs[j] = OpID(id)
			}

			opID := OpID(opDTO.ID)
			ops[opID] = &Operation{
				ID:             opID,
				Name:           opDTO.Name,
				ProcessingTime: dur,
				MachineType:    MachineType(opDTO.MachineType),
				Predecessors:   preds,
				Successors:     succs,
			}
		}

		jobs[i] = &Job{
			ID:   JobID(jobDTO.ID),
			Name: jobDTO.Name,
			Ops:  ops,
		}
	}

	machines := make([]*Machine, len(dto.Machines))
	for i, mDTO := range dto.Machines {
		m := NewMachine(MachineID(mDTO.ID), MachineType(mDTO.Type), mDTO.Name)
		machines[i] = &m
	}

	return Problem{
		Jobs:     jobs,
		Machines: machines,
	}, nil
}
