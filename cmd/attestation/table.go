package attestation

import (
	"os"
	"strings"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/jedib0t/go-pretty/v6/table"
)

func convPredSubtoRows(psMap map[string][]intoto.Statement) []table.Row {
	items := make([]table.Row, 0, len(psMap))

	for pred, statements := range psMap {
		for _, statement := range statements {
			var subjects []string
			for _, s := range statement.Subject {
				subjects = append(subjects, s.Name)
			}
			subjectsString := strings.Join(subjects, ", ")
			items = append(items, table.Row{
				pred,
				subjectsString,
			})
		}
	}
	return items
}

func printPredSubjTable(psMap map[string][]intoto.Statement) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Predicate", "Subjects"})
	rows := convPredSubtoRows(psMap)
	t.AppendRows(rows)
	t.Render()

}
