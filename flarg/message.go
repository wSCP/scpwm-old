package flarg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// Fail prints usage information to stdout and exits with non-zero status
func (p *Parser) Fail(msg string) {
	p.Usage(os.Stdout)
	fmt.Println("error:", msg)
	os.Exit(-1)
}

// Usage writes usage information of the provided []*spec to the given writer
func (p *Parser) Usage(w io.Writer) {
	var positionals, options []*spec
	for _, spec := range p.spec {
		if spec.positional {
			positionals = append(positionals, spec)
		} else {
			options = append(options, spec)
		}
	}

	fmt.Fprintf(w, "usage: %s ", filepath.Base(os.Args[0]))

	// write the option component of the usage message
	for _, spec := range options {
		if !spec.required {
			fmt.Fprint(w, "[")
		}
		fmt.Fprint(w, synopsis(spec, "--"+spec.long))
		if !spec.required {
			fmt.Fprint(w, "]")
		}
		fmt.Fprint(w, " ")
	}

	// write the positional component of the usage message
	for _, spec := range positionals {
		up := strings.ToUpper(spec.long)
		if spec.multiple {
			fmt.Fprintf(w, "[%s [%s ...]]", up, up)
		} else {
			fmt.Fprint(w, up)
		}
		fmt.Fprint(w, " ")
	}
	fmt.Fprint(w, "\n")
}

// Help writes the usage string of the provided []*spec followed by the full
// help string for each option
func (p *Parser) Help(w io.Writer) {
	var positionals, options []*spec
	for _, spec := range p.spec {
		if spec.positional {
			positionals = append(positionals, spec)
		} else {
			options = append(options, spec)
		}
	}

	p.Usage(w)

	// write the list of positionals
	if len(positionals) > 0 {
		fmt.Fprint(w, "\npositional arguments:\n")
		for _, spec := range positionals {
			fmt.Fprintf(w, "  %s\n", spec.long)
		}
	}

	// write the list of options
	if len(options) > 0 {
		fmt.Fprint(w, "\noptions:\n")
		const colWidth = 25
		for _, spec := range options {
			left := "  " + synopsis(spec, "--"+spec.long)
			if spec.short != "" {
				left += ", " + synopsis(spec, "-"+spec.short)
			}
			fmt.Fprint(w, left)
			if spec.help != "" {
				if len(left)+2 < colWidth {
					fmt.Fprint(w, strings.Repeat(" ", colWidth-len(left)))
				} else {
					fmt.Fprint(w, "\n"+strings.Repeat(" ", colWidth))
				}
				fmt.Fprint(w, spec.help)
			}
			fmt.Fprint(w, "\n")
		}
	}
}

func synopsis(spec *spec, form string) string {
	if spec.dest.Kind() == reflect.Bool {
		return form
	} else {
		return fmt.Sprintf("%s %s", form, strings.ToUpper(spec.long))
	}
}
