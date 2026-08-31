package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile"
)

type CompileConfigEntries struct {
	RowCount int    `json:"rowCount,omitempty"`
	Seed     int64  `json:"seed,omitempty"`
	Schema   string `json:"schema,omitempty"`
}

type CompileConfig struct {
	InputPath  string               `json:"inputPath"`
	OutputPath string               `json:"outputPath"`
	Config     CompileConfigEntries `json:"config"`
	Verbose    bool                 `json:"verbose,omitempty"`
}

const (
	ErrMissingConfig      = 3
	ErrInvalidConfig      = 4
	ErrInputPathRequired  = 12
	ErrOutputPathRequired = 13
	ErrCompileFailed      = 1
)

func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-seed-data <config>")
		fmt.Fprintln(os.Stderr, "  config: JSON string with inputPath, outputPath, and optional config parameters")
		os.Exit(ErrMissingConfig)
	}

	rawConfig := os.Args[1]
	var compileConfig CompileConfig
	if err := json.Unmarshal([]byte(rawConfig), &compileConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ErrInvalidConfig)
	}

	if compileConfig.InputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required")
		os.Exit(ErrInputPathRequired)
	}

	if compileConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Output path is required")
		os.Exit(ErrOutputPathRequired)
	}

	inputAbs, err := filepath.Abs(compileConfig.InputPath)
	if err == nil {
		compileConfig.InputPath = inputAbs
	}

	outputAbs, err := filepath.Abs(compileConfig.OutputPath)
	if err == nil {
		compileConfig.OutputPath = outputAbs
	}

	logInfo(compileConfig.Verbose, "Processing Morphe registry from: '%s'", compileConfig.InputPath)
	logInfo(compileConfig.Verbose, "Output seed SQL to: '%s'", compileConfig.OutputPath)

	morpheConfig := compile.DefaultMorpheCompileConfig(
		compileConfig.InputPath,
		compileConfig.OutputPath,
	)

	if compileConfig.Config.RowCount > 0 {
		morpheConfig.SeedConfig.RowCount = compileConfig.Config.RowCount
	}
	if compileConfig.Config.Seed != 0 {
		morpheConfig.SeedConfig.Seed = compileConfig.Config.Seed
	}
	if compileConfig.Config.Schema != "" {
		morpheConfig.SeedConfig.Schema = compileConfig.Config.Schema
	}

	logInfo(compileConfig.Verbose, "Starting seed data generation (rowCount=%d, seed=%d)...",
		morpheConfig.SeedConfig.RowCount, morpheConfig.SeedConfig.Seed)

	compileErr := compile.MorpheToSeedSQL(morpheConfig)
	if compileErr != nil {
		fmt.Fprintln(os.Stderr, "Seed data generation failed:", compileErr)
		os.Exit(ErrCompileFailed)
	}

	logInfo(compileConfig.Verbose, "Seed data generation completed successfully")
	os.Exit(0)
}
