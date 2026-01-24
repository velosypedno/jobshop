package main

import (
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/resource-allocation/chart"
	"github.com/velosypedno/resource-allocation/factory"
	"github.com/velosypedno/resource-allocation/parser"
	"github.com/velosypedno/resource-allocation/strategy/naive"
)

func main() {
	filepath := "./default_factory_config.json"
	machinesConfig, templates, err := parser.ParseFactoryConfig(filepath)
	if err != nil {
		panic(err)
	}
	factory := &factory.Factory{}
	factory.Configure(machinesConfig, templates)
	factory.SetPlanner(&naive.Strategy{})
	factory.AddJobByName("Bicycle")
	factory.AddJobByName("Bicycle")
	factory.AddJobByName("Scooter")
	factory.AddJobByName("Skateboard")
	factory.AddJobByName("Skateboard")
	factory.AddJobByName("Bicycle")
	factory.AddJobByName("Skateboard")
	factory.AddJobByName("Skateboard")
	factory.AddJobByName("Bicycle")
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
