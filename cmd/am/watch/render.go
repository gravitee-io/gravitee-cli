package watch

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gravitee-io/gio-cli/internal/cmdutil"
)

type AuditEvent struct {
	ID        string
	EventType string
	Status    string
	Actor     string
	Target    string
	Timestamp string
	RawTs     int64
}

type TypeCount struct {
	Type  string
	Count int
}

type DashboardStats struct {
	Total     int
	Successes int
	Failures  int
	TopTypes  []TypeCount
	TopErrors []TypeCount
}

type DashboardData struct {
	DomainName  string
	Workspace   string
	RefreshedAt string
	Events      []AuditEvent
	Stats       DashboardStats
}

func buildDashboardData(rawEvents []map[string]interface{}, domainName, workspace string) DashboardData {
	events := make([]AuditEvent, 0, len(rawEvents))
	for _, e := range rawEvents {
		ev := AuditEvent{
			ID:        cmdutil.StringField(e, "id"),
			EventType: cmdutil.StringField(e, "type"),
		}
		if outcome, ok := e["outcome"].(map[string]interface{}); ok {
			ev.Status = cmdutil.StringField(outcome, "status")
		}
		if actor, ok := e["actor"].(map[string]interface{}); ok {
			ev.Actor = cmdutil.StringField(actor, "displayName")
		}
		if target, ok := e["target"].(map[string]interface{}); ok {
			ev.Target = cmdutil.StringField(target, "displayName")
		}
		if ts, ok := e["timestamp"].(float64); ok {
			ev.RawTs = int64(ts)
			ev.Timestamp = time.UnixMilli(int64(ts)).UTC().Format("2006-01-02 15:04:05")
		}
		events = append(events, ev)
	}
	sortDesc(events)

	successes, failures := 0, 0
	typeCounts := make(map[string]int)
	errorCounts := make(map[string]int)
	for _, ev := range events {
		if ev.Status == "SUCCESS" {
			successes++
		} else {
			failures++
			errorCounts[ev.EventType]++
		}
		typeCounts[ev.EventType]++
	}

	return DashboardData{
		DomainName:  domainName,
		Workspace:   workspace,
		RefreshedAt: time.Now().UTC().Format("2006-01-02 15:04:05"),
		Events:      events,
		Stats: DashboardStats{
			Total:     len(events),
			Successes: successes,
			Failures:  failures,
			TopTypes:  topN(typeCounts, 5),
			TopErrors: topN(errorCounts, 5),
		},
	}
}

func render(data DashboardData, intervalSec int) string {
	var sb strings.Builder
	hr := strings.Repeat("─", 80)

	sb.WriteString(fmt.Sprintf("  Gravitee AM — %s (%s)\n", data.DomainName, data.Workspace))
	sb.WriteString(fmt.Sprintf("  %s\n\n", hr))

	successRate := 0
	if data.Stats.Total > 0 {
		successRate = int(math.Round(float64(data.Stats.Successes) / float64(data.Stats.Total) * 100))
	}
	sb.WriteString(fmt.Sprintf("  Events: %d    Success: %d    Failure: %d    Rate: %d%%\n\n",
		data.Stats.Total, data.Stats.Successes, data.Stats.Failures, successRate))

	if len(data.Stats.TopTypes) > 0 {
		sb.WriteString("  Event types:\n")
		for _, t := range data.Stats.TopTypes {
			barLen := 0
			if data.Stats.Total > 0 {
				barLen = int(math.Round(float64(t.Count) / float64(data.Stats.Total) * 30))
			}
			bar := strings.Repeat("█", barLen)
			sb.WriteString(fmt.Sprintf("    %-25s %s %d\n", t.Type, bar, t.Count))
		}
		sb.WriteString("\n")
	}

	if len(data.Stats.TopErrors) > 0 {
		sb.WriteString("  Top errors:\n")
		for _, e := range data.Stats.TopErrors {
			sb.WriteString(fmt.Sprintf("    %-30s %d\n", e.Type, e.Count))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("  Recent events:\n")
	maxEvents := 15
	if len(data.Events) < maxEvents {
		maxEvents = len(data.Events)
	}
	for _, ev := range data.Events[:maxEvents] {
		icon := "+"
		if ev.Status != "SUCCESS" {
			icon = "!"
		}
		target := ""
		if ev.Target != "" {
			target = " -> " + ev.Target
		}
		sb.WriteString(fmt.Sprintf("    %s %s %-20s %s%s\n",
			ev.Timestamp, icon, ev.EventType, ev.Actor, target))
	}

	sb.WriteString(fmt.Sprintf("\n  %s\n", hr))
	sb.WriteString(fmt.Sprintf("  Last refresh: %s    Interval: %ds    Ctrl+C to stop\n", data.RefreshedAt, intervalSec))
	return sb.String()
}

func sortDesc(events []AuditEvent) {
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].RawTs > events[j].RawTs
	})
}

func topN(counts map[string]int, n int) []TypeCount {
	result := make([]TypeCount, 0, len(counts))
	for k, v := range counts {
		result = append(result, TypeCount{Type: k, Count: v})
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	if len(result) > n {
		result = result[:n]
	}
	return result
}
