package optigo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type actionType int

const (
	atINCREMENT actionType = iota
	atAPPEND
	atASSIGN
)

type dataType int

const (
	dtSTRING dataType = iota
	dtINTEGER
	dtFLOAT
	dtBOOLEAN
)

type option struct {
	name     string
	unary    bool
	dest     reflect.Value
	action   actionType
	dataType dataType
}

func (o *option) parseValue(val string) (interface{}, error) {
	switch o.dataType {
	case dtSTRING:
		return val, nil
	case dtINTEGER:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i, nil
		} else {
			return nil, err
		}
	case dtFLOAT:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, nil
		} else {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unable to parse value: %s", val)
	}
}

type actions map[string]option

func parseAction(spec string, dest interface{}, actions map[string]option) error {
	unary := false
	var a actionType
	var t dataType
	if spec[len(spec)-1] == '+' {
		unary = true
		a = atINCREMENT
		spec = spec[0 : len(spec)-1]
	} else if spec[len(spec)-1] == '@' {
		a = atAPPEND
		spec = spec[0 : len(spec)-1]
	} else {
		a = atASSIGN
	}

	switch spec[len(spec)-2:] {
	case "=s":
		t = dtSTRING
		spec = spec[0 : len(spec)-2]
	case "=i":
		t = dtINTEGER
		spec = spec[0 : len(spec)-2]
	case "=f":
		t = dtFLOAT
		spec = spec[0 : len(spec)-2]
	default:
		if a == atINCREMENT {
			t = dtINTEGER
		} else {
			t = dtBOOLEAN
		}
		unary = true
	}

	if unary && a == atAPPEND {
		return fmt.Errorf("invalid spec, using @ to parse repeated options, but not specifying type with either =i =s or =f: %s", spec)
	}

	optionNames := strings.Split(spec, "|")
	name := optionNames[len(optionNames)-1]
	for _, opt := range optionNames {
		var dashName string
		if len(opt) == 1 {
			dashName = "-" + opt
		} else {
			dashName = "--" + opt
		}
		actions[dashName] = option{name, unary, reflect.ValueOf(dest), a, t}
	}
	return nil
}

func increment(val reflect.Value) reflect.Value {
	return reflect.ValueOf(int(val.Int()) + 1)
}

func push(arr reflect.Value, val interface{}) reflect.Value {
	return reflect.Append(arr, reflect.ValueOf(val))
}

func initResultMap(actions actions) map[string]interface{} {
	results := make(map[string]interface{})
	for _, opt := range actions {
		if opt.unary {
			if opt.action == atINCREMENT {
				results[opt.name] = int64(0)
			} else {
				results[opt.name] = false
			}
		} else {
			if opt.action == atAPPEND {
				switch opt.dataType {
				case dtSTRING:
					results[opt.name] = make([]string, 0)
				case dtINTEGER:
					results[opt.name] = make([]int64, 0)
				case dtFLOAT:
					results[opt.name] = make([]float64, 0)
				}
			} else {
				switch opt.dataType {
				case dtSTRING:
					results[opt.name] = ""
				case dtINTEGER:
					results[opt.name] = int64(0)
				case dtFLOAT:
					results[opt.name] = float64(0)
				case dtBOOLEAN:
					results[opt.name] = false
				}
			}
		}
	}
	return results
}

type OptionParser struct {
	actions actions
	Results map[string]interface{}
	Args []string
}

// NewParser generates an OptionParser object from the opts passed in.
// The opts strings should be in the form of:
//  alias|alias(+|=s|=i|=f)[@]
//  Example:
//  v|verbose+  this will set a "verbose" key in OptionParser.Results with an integer key for how many times the verbose option was repeated on the command line
func NewParser(opts []string) OptionParser {
	actions := make(actions)
	for _, spec := range opts {
		if err := parseAction(spec, nil, actions); err != nil {
			panic(err)
		}
	}
	results := initResultMap(actions)
	return OptionParser{actions,results,nil}
}

func NewInlineParser(opts map[string]interface{}) OptionParser {
	actions := make(actions)
	for spec, ref := range opts {
		if err := parseAction(spec, ref, actions); err != nil {
			panic(err)
		}
	}
	return OptionParser{actions, nil, nil}
}

func (o *OptionParser) ProcessAll(args []string) error {
	err := o.ProcessSome(args)
	if err != nil {
		return err
	} else {
		for _, opt := range o.Args {
			if opt[0] == '-' {
				return fmt.Errorf("Unknown option: %s", opt)
			}
		}
	}
	return nil
}

func (o *OptionParser) ProcessSome(args []string) error {
	o.Args = make([]string,0)
	for len(args) > 0 {
		if args[0] == "--" {
			o.Args = append(o.Args, args[1:]...)
			return nil
		}
		if opt, ok := o.actions[args[0]]; ok {
			var value interface{}
			var err error
			if opt.unary {
				value = true
				args = args[1:]
			} else {
				if len(args) < 2 {
					return fmt.Errorf("missing argument value for option: --%s", opt.name)
				} else {
					if value, err = opt.parseValue(args[1]); err != nil {
						return err
					}
				}
				args = args[2:]
			}

			if opt.dest.IsValid() {
				switch opt.action {
				case atINCREMENT:
					opt.dest.Elem().Set(increment(opt.dest.Elem()))
				case atAPPEND:
					opt.dest.Elem().Set(push(opt.dest.Elem(), value))
				case atASSIGN:
					if opt.dest.Kind() == reflect.Func {
						t := reflect.TypeOf(opt.dest.Interface())
						if t.NumIn() == 1 { 
							callbackArgs := make([]reflect.Value,1)
							callbackArgs[0] = reflect.ValueOf(value)
							opt.dest.Call(callbackArgs)
						} else {
							opt.dest.Call(nil)
						}
					} else {
						opt.dest.Elem().Set(reflect.ValueOf(value))
					}
				}
			} else {
				switch opt.action {
				case atINCREMENT:
					o.Results[opt.name] = increment(reflect.ValueOf(o.Results[opt.name])).Interface()
				case atAPPEND:
					o.Results[opt.name] = push(reflect.ValueOf(o.Results[opt.name]), value).Interface()
				case atASSIGN:
					o.Results[opt.name] = reflect.ValueOf(value).Interface()
				}
			}
		} else {
			o.Args = append(o.Args, args[0])
			args = args[1:]
		}
	}
	return nil
}
