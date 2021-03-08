package opalparser2

type errType string

const (
	errUnexpectedTerm     = "Unexpected terminator"
	errUnexpectedEOF      = "Unexpected end of file"
	errInlineTagNoName    = "Inline tag does not have a name"
	errInlineTagNoContent = "Inline tag does not have any content"
)
