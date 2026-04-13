package reporter

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/velosypedno/jobshop/internal/scheduler"
)

type Reporter struct {
	writer io.Writer
}

func New(w io.Writer) *Reporter {
	return &Reporter{writer: w}
}

func (f *Reporter) format(results []scheduler.PlanResult) (string, error) {
	var buf strings.Builder
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "Strategy\tType\tMakeSpan\tUtil %\tSched Time")
	fmt.Fprintln(w, "--------\t----\t--------\t------\t----------")

	for _, res := range results {
		_, err := fmt.Fprintf(w, "%s\t%s\t%v\t%.2f%%\t%v\n",
			res.Info.StrategyName,
			res.Info.StrategyType,
			res.Info.MakeSpan,
			res.Info.UtilizationLevel*100,
			res.Info.SchedulingTime,
		)
		if err != nil {
			return "", err
		}
	}

	w.Flush()
	return buf.String(), nil
}

func (r *Reporter) Generate(results []scheduler.PlanResult) error {
	content, err := r.format(results)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(r.writer, content)
	return err
}
