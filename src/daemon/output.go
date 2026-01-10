package daemon

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/sirupsen/logrus"
)

func formatContainerTable(containers []entity.Container) string {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	defer func() { _ = writer.Flush() }()

	header := "CONTAINER ID\tIMAGE\tCOMMAND\tCREATED\tSTATUS\tNAMES"
	_, err := fmt.Fprintln(writer, header)
	if err != nil {
		return ""
	}

	var tableContent strings.Builder
	tableContent.WriteString(header + "\n")

	for _, c := range containers {
		shortID := c.Id

		formattedCmd := fmt.Sprintf("\"%s\"", c.Command)

		createdAt := time.UnixMilli(c.CreatedAt)
		createdStr := formatTimeAgo(createdAt)

		statusStr := formatContainerStatus(c)

		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
			shortID,
			c.Image,
			formattedCmd,
			createdStr,
			statusStr,
			c.Name,
		)

		// 写入终端和字符串缓存
		_, err = fmt.Fprintln(writer, line)
		if err != nil {
			logrus.Error("error writing container table: %v", err)
			continue
		}
		tableContent.WriteString(line + "\n")
	}

	return tableContent.String()
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Hour:
		minutes := int(duration.Minutes())
		if minutes == 0 {
			return "a few seconds ago"
		}
		return fmt.Sprintf("%d %s ago", minutes, pluralize(minutes, "minute", "minutes"))
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		return fmt.Sprintf("%d %s ago", hours, pluralize(hours, "hour", "hours"))
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d %s ago", days, pluralize(days, "day", "days"))
	default:
		return t.Format("2006-01-02 15:04")
	}
}

func formatContainerStatus(c entity.Container) string {
	if c.Status == entity.ContainerRunning {
		return "Up " + formatTimeAgo(time.UnixMilli(c.CreatedAt))
	}

	exitAt := time.UnixMilli(c.ExitAt)
	exitAgo := formatTimeAgo(exitAt)
	exitCode := 0
	return fmt.Sprintf("Exited (%d) %s", exitCode, exitAgo)
}

func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}
