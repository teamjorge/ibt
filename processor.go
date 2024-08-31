package ibt

import (
	"context"
	"errors"
	"sort"

	"github.com/teamjorge/ibt/headers"
	"github.com/teamjorge/ibt/utilities"
)

type Processor interface {
	Process(input Tick, hasNext bool, session *headers.Session) error
	Whitelist() []string
}

func Process(ctx context.Context, stubs StubGroup, processors ...Processor) error {
	sort.Sort(stubs)

	for _, stub := range stubs {
		if err := process(ctx, stub, processors...); err != nil {
			return err
		}
	}

	return nil
}

func process(ctx context.Context, stub Stub, processors ...Processor) error {
	header := stub.header

	whitelist := buildWhitelist(header.VarHeader, processors...)

	parser := NewParser(stub.r, header, whitelist...)
	for {
		select {
		case <-ctx.Done():
			return errors.New("context cancelled")
		default:
		}

		tick, hasNext := parser.Next()
		for _, proc := range processors {
			if err := proc.Process(tick.Filter(proc.Whitelist()...), hasNext, header.SessionInfo); err != nil {
				return err
			}
		}

		if !hasNext {
			break
		}
	}

	return nil
}

// getcinoketeWhitelist compiles the whitelists from all processors and removes overlap
func buildWhitelist(vars map[string]headers.VarHeader, processors ...Processor) []string {
	whitelist := make([]string, 0)

	for _, proc := range processors {
		whitelist = append(whitelist, parseAndValidateWhitelist(vars, proc)...)
	}

	return utilities.GetDistinct(whitelist)
}

// parseWhitelist will retrieve vars when * is used and ensure a unique list
//
// Variables that are not found in the VarHeader will automatically be excluded.
func parseAndValidateWhitelist(vars map[string]headers.VarHeader, processor Processor) []string {
	whitelist := processor.Whitelist()

	if len(whitelist) == 0 {
		return headers.AvailableVars(vars)
	}

	for _, col := range whitelist {
		if col == "*" {
			return headers.AvailableVars(vars)
		}
	}

	columns := make([]string, 0)

	// Ensure only valid columns are added
	for _, col := range whitelist {
		if _, ok := vars[col]; ok {
			columns = append(columns, col)
		}
	}

	return columns
}
