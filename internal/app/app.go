package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/chart"
	"github.com/velosypedno/resource-allocation/internal/parser"
	"github.com/velosypedno/resource-allocation/internal/reporter"
	"github.com/velosypedno/resource-allocation/internal/reporter/formatter"
	"github.com/velosypedno/resource-allocation/internal/scheduler"
	"go.uber.org/zap"
)

type App struct {
	Scheduler *scheduler.Scheduler
}

func New(machinesConfig []parser.MachineConfig, templates []base.JobTemplate, strategies []base.Strategy) *App {
	s := &scheduler.Scheduler{}
	s.Configure(machinesConfig, templates)
	s.SetPlanners(strategies...)

	return &App{
		Scheduler: s,
	}
}
func (a *App) Run(startTime time.Time, orders []parser.OrderDTO, customName string) error {
	outputDir := "results"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating directory: %v", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	baseName := fmt.Sprintf("plan_%s%s", timestamp, customName)

	logPath := filepath.Join(outputDir, baseName+".log")

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logPath}

	logger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	a.Scheduler.SetLogger(logger)

	logger.Info("Starting application run",
		zap.Time("start_time", startTime),
		zap.Int("orders_count", len(orders)),
	)

	results, err := a.Scheduler.Plan(orders, startTime)
	if err != nil {
		logger.Error("Planning failed", zap.Error(err))
		return fmt.Errorf("during planning: %v", err)
	}
	rep := reporter.NewReporter(os.Stdout, &formatter.TableFormatter{})

	if err := rep.Generate(results); err != nil {
		logger.Warn("Could not generate text report", zap.Error(err))
	}

	solutionsChart := chart.GenerateFromSolutions(results, a.Scheduler.Machines)

	chartPath := filepath.Join(outputDir, baseName+".html")
	outputFile, err := os.Create(chartPath)
	if err != nil {
		return fmt.Errorf("creating output file: %v", err)
	}
	defer outputFile.Close()

	err = solutionsChart.Render(outputFile)
	if err != nil {
		return fmt.Errorf("error rendering chart: %v", err)
	}

	fmt.Printf("Successfully generated chart: %s\n", chartPath)
	fmt.Printf("Log file created: %s\n", logPath)

	logger.Info("Run completed successfully", zap.String("chart_path", chartPath))

	return nil
}
