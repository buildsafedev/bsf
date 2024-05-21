package attestation

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func convPredSubtoRows(psMap map[string][]string) []table.Row {
	items := make([]table.Row, 0, len(psMap))

	for pred, subjects := range psMap {
		for _, sub := range subjects {
			items = append(items, table.Row{
				pred,
				sub,
			})
		}
	}
	return items
}

func printPredSubjTable(psMap map[string][]string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Predicate", "Subject"})
	rows := convPredSubtoRows(psMap)
	t.AppendRows(rows)
	t.Render()

}
