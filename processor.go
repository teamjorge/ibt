package ibt

import (
	"github.com/teamjorge/ibt/headers"
)

type Processor interface {
	Process(input map[string]headers.VarHeader, hasNext bool, session *headers.Session) error
	Whitelist() []string
}

func Process(reader headers.Reader, header headers.Header, processors ...Processor) error {
	whitelist := make([]string, 0)

	for _, proc := range processors {
		whitelist = append(whitelist, proc.Whitelist()...)
	}

	parser := NewParser(reader, header, whitelist...)

	for {
		tick, hasNext := parser.Next()
		for _, proc := range processors {
			if err := proc.Process(tick, hasNext, header.SessionInfo()); err != nil {
				return err
			}
		}

		if !hasNext {
			break
		}
	}

	return nil
}
