package cpumanager

import "fmt"

const (
	// ErrorSMTAlignment represents the type of an SMTAlignmentError
	ErrorSMTAlignment = "SMTAlignmentError"
)

// SMTAlignmentError represents an error due to SMT alignment
type SMTAlignmentError struct {
	RequestedCPUs int
	CpusPerCore   int
}

func (e SMTAlignmentError) Error() string {
	return fmt.Sprintf("SMT Alignment Error: requested %d cpus not multiple cpus per core = %d", e.RequestedCPUs, e.CpusPerCore)
}

func (e SMTAlignmentError) Type() string {
	return ErrorSMTAlignment
}
