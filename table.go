package commons

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

func RenderTable(header []string, rows [][]string) {
	alignments := []tw.Align{
		tw.AlignDefault,
	}
	for len(alignments) < len(header) {
		alignments = append(alignments, tw.AlignRight)
	}
	tableConfig := tablewriter.WithConfig(tablewriter.Config{
		Header: tw.CellConfig{
			Formatting: tw.CellFormatting{AutoFormat: tw.Off},
			Alignment: tw.CellAlignment{Global: tw.AlignLeft},
		}},
	)
	alignmentConfig := tablewriter.WithAlignment(alignments)
	table := tablewriter.NewTable(os.Stdout, tableConfig, alignmentConfig)
	table.Header(header)
	table.Bulk(rows)
	table.Render()
}