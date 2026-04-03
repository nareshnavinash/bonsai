package ui

import (
	"fmt"
	"strings"
)

func PrintTable(headers []string, rows [][]string) {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	gap := "    "
	parts := make([]string, len(headers))
	for i, h := range headers {
		parts[i] = fmt.Sprintf("%-*s", widths[i], h)
	}
	fmt.Println(strings.Join(parts, gap))

	for _, row := range rows {
		for i := range parts {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			parts[i] = fmt.Sprintf("%-*s", widths[i], cell)
		}
		fmt.Println(strings.Join(parts, gap))
	}
}
