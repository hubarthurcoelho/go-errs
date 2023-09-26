package errs

import (
	"bytes"
	"fmt"
	"log"
)

type Code int
type Message string
type Params interface{}

type HTTPError struct {
	kind    kind
	message Message
	err     error
	params  Params
	stack
}

func (e *HTTPError) Log() {
	b := new(bytes.Buffer)

	e.printStack(b)

	pad(b, ": ")
	b.WriteString("\n")

	b.WriteString(e.Error())

	if e.params != nil {
		pad(b, ": ")
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("%+v", e.params))
	}

	log.Println(b.String())
}

func (e *HTTPError) Error() string {
	b := new(bytes.Buffer)
	if e.kind != 0 {
		pad(b, ": ")
		b.WriteString(e.kind.String())
	}
	if e.message != "" {
		pad(b, ": ")
		b.WriteString(string(e.message))
	}
	if e.err != nil {
		// Indent on new line if we are cascading non-empty Upspin errors.
		if prevErr, ok := e.err.(*HTTPError); ok {
			if !prevErr.isZero() {
				pad(b, separator)
				b.WriteString(e.err.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.err.Error())
		}
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

func (e *HTTPError) isZero() bool {
	return e.kind == 0 && e.message == "" && e.err == nil
}

func (he *HTTPError) Status() int {
	if he.kind != 0 {
		return he.kind.HttpStatus()
	}

	copy := *he
	for copy.err != nil {
		if err, ok := copy.err.(*HTTPError); ok {
			if err.kind != 0 {
				return err.kind.HttpStatus()
			} else {
				copy = *err
			}
		} else {
			break
		}
	}

	return he.kind.HttpStatus()
}

func (he *HTTPError) Message() string {
	b := new(bytes.Buffer)
	b.WriteString(string(he.message))

	copy := *he
	for copy.err != nil {
		if err, ok := copy.err.(*HTTPError); ok {
			pad(b, ": ")
			b.WriteString(string(err.message))
			copy = *err
		} else {
			pad(b, ": ")
			b.WriteString(copy.err.Error())
			break
		}

	}
	return b.String()
}

// E creates an error from arguments. The supported arguments are:
//
//   - Kind: the kind of the error (see herrs/kinds.go).
//   - Message: Message type, represents the error message.
//   - error: error, represents an inner error to wrap.
//
// The function returns an error that can be used to handle and propagate
// errors in a structured way. If no inner error is provided, the function
// returns a new error with the provided arguments. If an inner error is
// provided, the function returns a new error that wraps the inner error.
//
// Example usage:
//
//	err := errs.E(Op("addUser"), Code(400), Message("something went wrong"))
//
//	innerErr := someFunction()
//	err = errs.E(Op("outerErr"), Code(500), innerErr)
//
// Note: The function will panic if no arguments are provided or if
// unsupported argument types are provided.
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &HTTPError{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Message:
			e.message = arg
		case kind:
			e.kind = arg
		case *HTTPError:
			// Make a copy
			copy := *arg
			e.err = &copy
		case error:
			e.err = arg
		case Params:
			e.params = arg
		default:
			panic(fmt.Sprintf("unknown type %T, value %v in error call", arg, arg))
		}
	}

	// Populate stack information (only in debug mode).
	e.populateStack()

	prev, ok := e.err.(*HTTPError)
	if !ok {
		return e
	}

	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind, file name or user name
	// twice.
	if prev.message == e.message {
		prev.message = ""
	}
	if prev.kind == e.kind {
		prev.kind = 0
	}
	return e
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
