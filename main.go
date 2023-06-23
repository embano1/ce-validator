package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	ce "github.com/cloudevents/sdk-go/v2"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(flags []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("ce-validator", flag.ContinueOnError)
	fs.SetOutput(stdout)
	stdin := fs.Bool("i", false, "Read event from standard input")
	file := fs.String("f", "", "Read event from file")

	if err := fs.Parse(flags); err != nil {
		return err
	}

	input, err := getInput(*stdin, *file)
	if err != nil {
		return fmt.Errorf("failed to read event: %w", err)
	}

	version, err := validate(input)
	if err != nil {
		return fmt.Errorf("failed to validate event: %w", err)
	}

	_, err = fmt.Fprintf(stdout, "Event is a valid CloudEvent (spec version: %s)\n", version)
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}

	return nil
}

func getInput(stdin bool, filePath string) (io.Reader, error) {
	if stdin {
		return bufio.NewReader(os.Stdin), nil
	}

	if filePath == "" {
		return nil, fmt.Errorf("input file not specified (use -i to read from standard input)")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}

// returns event spec version and any validation errors
func validate(r io.Reader) (string, error) {
	var event ce.Event

	b, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to read event: %w", err)
	}

	if err := json.Unmarshal(b, &event); err != nil {
		return "", fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return event.SpecVersion(), event.Validate()
}
