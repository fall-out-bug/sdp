package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

func logExportCmd() *cobra.Command {
	var format string
	c := &cobra.Command{
		Use:   "export",
		Short: "Export events as CSV or JSON",
		Long:  `Export events for analysis. Supports --format=csv or --format=json.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogExport(format)
		},
	}
	c.Flags().StringVar(&format, "format", "csv", "Export format: csv or json")
	return c
}

func runLogExport(format string) error {
	path, err := evidenceLogPath()
	if err != nil {
		return err
	}
	r := evidence.NewReader(path)
	events, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}
	if len(events) == 0 {
		fmt.Println("No events to export.")
		return nil
	}

	exporter := evidence.NewExporter()
	var output string
	if format == "json" {
		output, err = exporter.ToJSON(events)
	} else {
		output, err = exporter.ToCSV(events)
	}
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}
	fmt.Println(output)
	return nil
}
