package main

import (
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/resource-allocation/base"
	"github.com/velosypedno/resource-allocation/chart"
	"github.com/velosypedno/resource-allocation/factory"
	"github.com/velosypedno/resource-allocation/strategy/naive"
)

func main() {
	factory := &factory.Factory{}
	factory.AddJob(base.BicycleJobTemplate)
	factory.AddJob(base.BicycleJobTemplate)
	factory.AddJob(base.ScooterJobTemplate)
	factory.AddJob(base.SkateboardJobTemplate)
	factory.AddJob(base.SkateboardJobTemplate)
	factory.AddJob(base.BicycleJobTemplate)
	factory.AddJob(base.SkateboardJobTemplate)
	factory.AddJob(base.SkateboardJobTemplate)
	factory.AddJob(base.BicycleJobTemplate)
	factory.AddMachine(base.Assembler)
	factory.AddMachine(base.Assembler)
	factory.AddMachine(base.Smeltery)
	factory.AddMachine(base.Smeltery)
	factory.AddMachine(base.Extruder)
	factory.AddMachine(base.Sawmill)
	factory.AddMachine(base.Sawmill)
	factory.AddMachine(base.Assembler)
	factory.AddMachine(base.Assembler)
	factory.SetPlanner(&naive.Strategy{})
	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)
	solution, metaInfo, err := factory.Plan(startTime)
	if err != nil {
		panic(err)
	}
	fmt.Println(solution)

	solutionChart := chart.GenerateFromSolution(solution, factory.Machines, metaInfo)
	f, err := os.Create("bar.html")
	if err != nil {
		panic(err)
	}
	err = solutionChart.Render(f)
	if err != nil {
		panic(err)
	}
}
