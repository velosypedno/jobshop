package report

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/velosypedno/jobshop/internal/core"
	"github.com/velosypedno/jobshop/internal/engine"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const renderItemJS = `function (params, api) {
    var categoryIndex = api.value(0); 
    var start = api.coord([api.value(1), categoryIndex]); 
    var end = api.coord([api.value(2), categoryIndex]); 
    var height = api.size([0, 1])[1] * 0.6; 
    var width = Math.max(end[0] - start[0], 1);

    var rectShape = echarts.graphic.clipRectByRect({ 
        x: start[0], 
        y: start[1] - height / 2, 
        width: width, 
        height: height 
    }, { 
        x: params.coordSys.x, 
        y: params.coordSys.y, 
        width: params.coordSys.width, 
        height: params.coordSys.height 
    });

    var operationName = api.value(3);
    
    var isLargeEnough = width > 90;

    var rect = { 
        type: 'rect', 
        transition: ['shape'], 
        shape: rectShape, 
        style: api.style({
            stroke: '#ffffff',
            lineWidth: 0.5
        })
    };

    if (isLargeEnough) {
        rect.textContent = {
            style: {
                text: operationName,
                fill: '#fff',
                fontSize: 10,
                fontFamily: 'sans-serif',
                textShadowColor: 'rgba(0,0,0,0.4)',
                textShadowBlur: 2,
                overflow: 'truncate',
                ellipsis: '..'
            }
        };
        rect.textConfig = {
            position: 'inside'
        };
    }

    return rectShape && rect; 
}`

const renderTooltipJS = `function(p){
	var dateStart = new Date(p.value[1]);
	var dateEnd = new Date(p.value[2]);
	var timeStart = dateStart.toLocaleTimeString('uk-UA', {hour12: false});
	var timeEnd = dateEnd.toLocaleTimeString('uk-UA', {hour12: false});
	return '<b>' + p.value[4] + '</b><br/>' + 
			'Operation: ' + p.value[3] + '<br/>' + 
			timeStart + ' - ' + timeEnd;
}`

func sortMachines(machines []*core.Machine) {
	sort.Slice(machines, func(i, j int) bool {
		if machines[i].Type != machines[j].Type {
			return machines[i].Type < machines[j].Type
		}
		return machines[i].ID < machines[j].ID
	})
}

func generateMachineIndexMap(machines []*core.Machine) map[core.MachineID]int {
	mMap := make(map[core.MachineID]int)
	for i, m := range machines {
		mMap[m.ID] = i
	}
	return mMap
}

func generateYAxisCategories(machines []*core.Machine) []string {
	categories := make([]string, 0, len(machines))
	for _, m := range machines {
		categories = append(categories, fmt.Sprintf("%s [ID: %d]", m.Name, m.ID))
	}
	return categories
}

func createBaseCustomChart(machines []*core.Machine, period core.Period, description string) *charts.Custom {
	chart := charts.NewCustom()

	lineCount := strings.Count(description, "\n") + 1
	descriptionHeight := lineCount * 27
	topOffset := descriptionHeight + 50

	baseHeight := len(machines)*75 + 120
	totalHeight := baseHeight + descriptionHeight

	chart.Initialization.Width = "90%"
	chart.Initialization.Height = fmt.Sprintf("%dpx", totalHeight)

	yAxisCategories := generateYAxisCategories(machines)

	chart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Production Plan",
			Subtitle: description,
			Left:     "12%",
			TitleStyle: &opts.TextStyle{
				FontSize:   24,
				FontWeight: "bold",
				Color:      "#1a1a1a",
			},
			SubtitleStyle: &opts.TextStyle{
				Color:      "#333",
				FontSize:   16,
				FontWeight: "500",
				LineHeight: 25,
				FontFamily: "monospace, Courier New",
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      opts.Bool(true),
			Trigger:   "item",
			Formatter: opts.FuncOpts(renderTooltipJS),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "time",
			Min:       period.Start.UnixMilli(),
			Max:       period.End.UnixMilli(),
			SplitLine: &opts.SplitLine{Show: opts.Bool(true)},
			AxisLabel: &opts.AxisLabel{
				Show:       opts.Bool(true),
				Formatter:  "{HH}:{mm}:{ss}",
				FontSize:   14,
				FontWeight: "bold",
				Color:      "#333",
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "category",
			Data:      yAxisCategories,
			SplitLine: &opts.SplitLine{Show: opts.Bool(true)},
			AxisLabel: &opts.AxisLabel{
				Show:       opts.Bool(true),
				FontSize:   14,
				FontWeight: "600",
				Color:      "#222",
				Margin:     15,
			},
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  opts.Bool(true),
			Right: "20",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  opts.Bool(true),
					Title: "Download PNG",
					Type:  "png",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  opts.Bool(true),
					Title: "Data View",
					Lang:  []string{"Data View", "Close", "Refresh"},
				},
			},
		}),
		charts.WithGridOpts(opts.Grid{
			Top:          fmt.Sprintf("%dpx", topOffset),
			Left:         "5%",
			Right:        "15%",
			ContainLabel: opts.Bool(true),
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:   opts.Bool(true),
			Orient: "vertical",
			Right:  "5px",
			Top:    fmt.Sprintf("%dpx", topOffset),
			Type:   "scroll",
		}),
	)

	return chart
}

func addSolutionSeries(chart *charts.Custom, solution *core.Solution, problemCtx *core.ProblemContext) {
	mMap := generateMachineIndexMap(problemCtx.Problem.Machines)

	for _, job := range problemCtx.Problem.Jobs {
		var seriesData []opts.CustomData
		fullJobName := fmt.Sprintf("%s [%d]", job.Name, job.ID)

		for _, opID := range problemCtx.GetOperationsByJob(job.ID) {
			opSolution := solution.OperationMap[opID]
			op, _ := problemCtx.GetOperation(opID)
			seriesData = append(seriesData, opts.CustomData{
				Value: []interface{}{
					mMap[opSolution.MachineID],
					problemCtx.Problem.StartTime.Add(opSolution.Offset).UnixMilli(),
					problemCtx.Problem.StartTime.Add(opSolution.Offset).Add(opSolution.Duration).UnixMilli(),
					op.Name,
					fullJobName,
				},
			})
		}

		chart.AddSeries(fullJobName, seriesData).
			SetSeriesOptions(
				charts.WithCustomChartOpts(opts.CustomChart{
					RenderItem: opts.FuncOpts(renderItemJS),
				}),
			)
	}
}

func formatStrategyDescription(report engine.Report) string {
	execTime := report.StrategyMetrics.SchedulingTime.Round(time.Millisecond).String()
	makespan := report.SolutionMetrics.MakeSpan.String()
	utilization := fmt.Sprintf("%.1f%%", report.SolutionMetrics.UtilizationLevel*100)

	line1 := fmt.Sprintf("STRATEGY: %s", strings.ToUpper(report.StrategyMetrics.StrategyType))

	line2 := fmt.Sprintf("TIME: %s  │  MAKESPAN: %s  │  UTILIZATION: %s",
		execTime, makespan, utilization)

	return fmt.Sprintf(
		"%s\n%s\n"+
			"─────────────────────────────────────────────────────────────────\n"+
			"%s",
		line1,
		line2,
		report.StrategyMetrics.StrategyDescription,
	)
}

func GenerateFromSolutions(
	problem *core.Problem,
	reports []engine.Report,
) *components.Page {
	problemCtx := core.NewProblemContext(problem)

	sortMachines(problem.Machines)

	page := components.NewPage()
	page.SetLayout(components.PageNoneLayout)
	page.PageTitle = "Multi-Strategy Comparison"

	for _, r := range reports {
		period := r.Solution.GetPeriod(problem.StartTime)
		description := formatStrategyDescription(r)

		chart := createBaseCustomChart(problem.Machines, period, description)
		addSolutionSeries(chart, r.Solution, problemCtx)

		page.AddCharts(chart)
	}

	return page
}

type GanttChartsReporter struct {
	writer io.Writer
}

func NewGanttCharts(w io.Writer) *GanttChartsReporter {
	return &GanttChartsReporter{writer: w}
}

func (r *GanttChartsReporter) Report(problem *core.Problem, results []engine.Report) error {
	page := GenerateFromSolutions(problem, results)
	return page.Render(r.writer)

}
