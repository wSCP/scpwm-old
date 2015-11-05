// Package flarg parses command line arguments using the fields from a struct.
// derived from: https://github.com/alexflint/go-arg
package flarg

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// MustParse processes command line arguments and exits upon failure
func MustParse(dest interface{}, exclude ...string) {
	p, err := New(dest, exclude...)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	help, err := p.Parse(os.Args[1:])
	if help {
		p.Help(os.Stdout)
		os.Exit(0)
	}
	if err != nil {
		p.Fail(err.Error())
	}
}

// Parse processes command line arguments and stores them in dest
func Parse(dest interface{}, exclude ...string) error {
	p, err := New(dest, exclude...)
	if err != nil {
		return err
	}
	_, err = p.Parse(os.Args[1:])
	return err
}

// Parser represents a set of command line options with destination values
type Parser struct {
	spec    []*spec
	invalid []string
}

type spec struct {
	dest       reflect.Value
	long       string
	short      string
	multiple   bool
	required   bool
	positional bool
	help       string
	wasPresent bool
}

var defaultInvalid []string = []string{"--"}

var (
	NoPointer          = Xrror("%s is not a pointer (did you forget an ampersand?)").Out
	NoStructPointer    = Xrror("%T is not a struct pointer").Out
	UnsupportedField   = Xrror("%s.%s: %s fields are not supported").Out
	LongShortArguments = Xrror("%s.%s: short arguments must be one character only").Out
	UnrecognizedTag    = Xrror("unrecognized tag '%s' on field %s").Out
)

// NewParser constructs a parser from a list of destination structs
func New(dest interface{}, exclude ...string) (*Parser, error) {
	var specs []*spec
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		panic(NoPointer(v.Type()).Error())
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic(NoStructPointer(dest).Error())
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		// Check for the ignore switch in the tag
		field := t.Field(i)

		tag := field.Tag.Get("arg")
		if tag == "-" {
			continue
		}

		spec := spec{
			long: strings.ToLower(field.Name),
			dest: v.Field(i),
		}

		// Get the scalar type for this field
		scalarType := field.Type
		if scalarType.Kind() == reflect.Slice {
			spec.multiple = true
			scalarType = scalarType.Elem()
			if scalarType.Kind() == reflect.Ptr {
				scalarType = scalarType.Elem()
			}
		}

		// Check for unsupported types
		switch scalarType.Kind() {
		case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Struct,
			reflect.Complex64, reflect.Complex128:
			return nil, UnsupportedField(t.Name(), field.Name, scalarType.Kind())
		}

		// Look at the tag
		if tag != "" {
			for _, key := range strings.Split(tag, ",") {
				var value string
				if pos := strings.Index(key, ":"); pos != -1 {
					value = key[pos+1:]
					key = key[:pos]
				}

				switch {
				case strings.HasPrefix(key, "--"):
					spec.long = key[2:]
				case strings.HasPrefix(key, "-"):
					if len(key) != 2 {
						return nil, LongShortArguments(t.Name(), field.Name)
					}
					spec.short = key[1:]
				case key == "required":
					spec.required = true
				case key == "positional":
					spec.positional = true
				case key == "help":
					spec.help = value
				default:
					return nil, UnrecognizedTag(key, tag)
				}
			}
		}
		specs = append(specs, &spec)
	}
	invalid := append(defaultInvalid, exclude...)
	return &Parser{specs, invalid}, nil
}

var InvalidArg = Xrror("%s is not a valid argument").Out

func contains(ss []string, s string) bool {
	for _, st := range ss {
		if st == s {
			return true
		}
	}
	return false
}

// Parse processes the given command line option, storing the results in the field
// of the structs from which NewParser was constructed
func (p *Parser) Parse(args []string) (bool, error) {
	// If -h or --help were specified then print usage
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true, nil
		}
		if contains(p.invalid, arg) {
			return false, InvalidArg(arg)
		}
	}

	// Process all command line arguments
	err := process(p.spec, args)
	if err != nil {
		return false, err
	}

	// Validate
	return false, validate(p.spec)
}

var (
	ProcessError    = Xrror("error processing %s: %v").Out
	UnknownArgument = Xrror("unknown argument %s").Out
	MissingValue    = Xrror("missing value for %s").Out
	Required        = Xrror("%s is required").Out
	PositionalError = Xrror("too many positional arguments at '%s'").Out
)

func process(specs []*spec, args []string) error {
	// construct a map from --option to spec
	optionMap := make(map[string]*spec)
	for _, spec := range specs {
		if spec.positional {
			continue
		}
		if spec.long != "" {
			optionMap[spec.long] = spec
		}
		if spec.short != "" {
			optionMap[spec.short] = spec
		}
	}

	// process each string from the command line
	var positioned bool
	var positionals []string

	// must use explicit for loop, not range, because we manipulate i inside the loop
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positioned = true
			continue
		}

		if !strings.HasPrefix(arg, "-") || positioned {
			positionals = append(positionals, arg)
			continue
		}

		// check for an equals sign, as in "--foo=bar"
		var value string
		opt := strings.TrimLeft(arg, "-")
		if pos := strings.Index(opt, "="); pos != -1 {
			value = opt[pos+1:]
			opt = opt[:pos]
		}

		// lookup the spec for this option
		spec, ok := optionMap[opt]
		if !ok {
			return UnknownArgument(arg)
		}
		spec.wasPresent = true

		// deal with the case of multiple values
		if spec.multiple {
			var values []string
			if value == "" {
				for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					values = append(values, args[i+1])
					i++
				}
			} else {
				values = append(values, value)
			}
			err := setSlice(spec.dest, values)
			if err != nil {
				return ProcessError(arg, err)
			}
			continue
		}

		// if it's a flag and it has no value then set the value to true
		if spec.dest.Kind() == reflect.Bool && value == "" {
			value = "true"
		}

		// if we have something like "--foo" then the value is the next argument
		if value == "" {
			if i+1 == len(args) || strings.HasPrefix(args[i+1], "-") {
				return MissingValue(arg)
			}
			value = args[i+1]
			i++
		}

		err := setScalar(spec.dest, value)
		if err != nil {
			return ProcessError(arg, err)
		}
	}

	// process positionals
	for _, spec := range specs {
		if spec.positional {
			if spec.multiple {
				err := setSlice(spec.dest, positionals)
				if err != nil {
					return ProcessError(spec.long, err)
				}
				positionals = nil
			} else if len(positionals) > 0 {
				err := setScalar(spec.dest, positionals[0])
				if err != nil {
					return ProcessError(spec.long, err)
				}
				positionals = positionals[1:]
			} else if spec.required {
				return Required(spec.long)
			}
		}
	}

	if len(positionals) > 0 {
		return PositionalError(positionals[0])
	}

	return nil
}

func validate(spec []*spec) error {
	for _, arg := range spec {
		if !arg.positional && arg.required && !arg.wasPresent {
			return Required(arg.long)
		}
	}
	return nil
}

var Unwritable = Xrror("%v field is not writable").Out

func setSlice(dest reflect.Value, values []string) error {
	if !dest.CanSet() {
		return Unwritable(dest)
	}

	var ptr bool
	elem := dest.Type().Elem()
	if elem.Kind() == reflect.Ptr {
		ptr = true
		elem = elem.Elem()
	}

	for _, s := range values {
		v := reflect.New(elem)
		if err := setScalar(v.Elem(), s); err != nil {
			return err
		}
		if ptr {
			v = v.Addr()
		}
		dest.Set(reflect.Append(dest, v.Elem()))
	}
	return nil
}

var (
	NotExported = Xrror("%v field is not exported").Out
	NotScalar   = Xrror("not a scalar type: %s").Out
)

func setScalar(v reflect.Value, s string) error {
	if !v.CanSet() {
		return NotExported(v)
	}

	switch v.Kind() {
	case reflect.String:
		v.Set(reflect.ValueOf(s))
	case reflect.Bool:
		x, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(x))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(x).Convert(v.Type()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(x).Convert(v.Type()))
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(x).Convert(v.Type()))
	default:
		return NotScalar(v.Kind())
	}
	return nil
}
