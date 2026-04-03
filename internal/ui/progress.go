package ui

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	total int64
	width int
	label string
}

func NewProgressBar(total int64, label string) *ProgressBar {
	return &ProgressBar{
		total: total,
		width: 30,
		label: label,
	}
}

func (p *ProgressBar) Update(completed int64) {
	if p.total <= 0 {
		return
	}
	pct := float64(completed) / float64(p.total)
	filled := int(float64(p.width) * pct)
	if filled > p.width {
		filled = p.width
	}
	empty := p.width - filled
	bar := strings.Repeat("=", filled)
	if empty > 0 {
		bar += ">"
		bar += strings.Repeat(" ", empty-1)
	}
	fmt.Printf("\r%s[%s] %d%% %s/%s",
		p.label, bar, int(pct*100),
		FormatBytes(completed), FormatBytes(p.total))
}

func (p *ProgressBar) Done() {
	fmt.Println()
}
