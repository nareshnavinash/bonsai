package ui

import (
	"fmt"
	"strings"
	"time"
)

func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(b)/float64(div), units[exp])
}

func FormatRelativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case d < 30*24*time.Hour:
		weeks := int(d.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	default:
		months := int(d.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
}

func FormatRelativeFuture(t time.Time) string {
	d := time.Until(t)
	switch {
	case d < time.Minute:
		return "less than a minute"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute from now"
		}
		return fmt.Sprintf("%d minutes from now", m)
	default:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour from now"
		}
		return fmt.Sprintf("%d hours from now", h)
	}
}

func TruncateID(digest string) string {
	s := strings.TrimPrefix(digest, "sha256:")
	if len(s) > 12 {
		return s[:12]
	}
	return s
}
