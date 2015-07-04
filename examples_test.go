package optigo

import (
	"fmt"
)

func ExampleOptionParser() {
	op := NewParser([]string{
		"v|verbose+",
	})
	args := []string{
		"-v",
		"--verbose",
		"-v",
		"extra",
	}
	if err := op.ProcessAll(args); err != nil {
		panic(err)
	}

	fmt.Printf("verbose: %d\n", op.Results["verbose"])
	fmt.Printf("unparsed args: %v\n", op.Args)

	// Output:
	// verbose: 3
	// unparsed args: [extra]
}

func ExampleNewParser() {
	// Note that all values will be stored in OptionParser.Results after a Process function
	// is called.  The Result key will be stored as the last alias.
	op := NewParser([]string{
		// Allow for repeated `--inc` or `--increment` options.  Each one of
		// the aliases is repeated the value is incrased by one.
		"inc|increment+",

		// Allow for `-S string` or `--string-list string` options.  The string
		// values will be stored in a slice in order of appearance.
		"S|string-list=s@",

		// Allow for `-I 123` or `--int-list 123` options.  The int values will
		// be stored in a slice in order of appearance
		"I|int-list=i@",

		// Allow for `-F 1.23` or `--float-list 1.23` options.  The float values
		// will be stored in a slice in order of appearance
		"F|float-list=f@",

		// Allow for `-s string` or `--string-value string`.
		"s|string-value=s",

		// Allow for `-i 123` or `--int-value 123`.
		"i|int-value=i",

		// Allow for `-f 1.23` or `--float-value 1.23`.
		"f|float-value=f",

		// Allow for `-b` or `--bool`.
		"b|bool",
	})

	args := []string{
		"--inc",
		"--increment",
		"-S", "A",
		"--string-list", "B",
		"-I", "1",
		"--int-list", "2",
		"-F", ".1",
		"--float-list", ".2",
		"-s", "hey",
		"-i", "42",
		"-f", "3.141593",
		"--bool",
	}

	if err := op.ProcessAll(args); err != nil {
		panic(err)
	}

	fmt.Printf("increment: %d\n", op.Results["increment"])
	fmt.Printf("string-list: %v\n", op.Results["string-list"])
	fmt.Printf("int-list: %v\n", op.Results["int-list"])
	fmt.Printf("float-list: %v\n", op.Results["float-list"])
	fmt.Printf("string-value: %s\n", op.Results["string-value"])
	fmt.Printf("int-value: %d\n", op.Results["int-value"])
	fmt.Printf("float-value: %f\n", op.Results["float-value"])
	fmt.Printf("bool: %t\n", op.Results["bool"])

	// Output:
	// increment: 2
	// string-list: [A B]
	// int-list: [1 2]
	// float-list: [0.1 0.2]
	// string-value: hey
	// int-value: 42
	// float-value: 3.141593
	// bool: true
}

func ExampleNewParser_nonUnique() {
	defer func () {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	NewParser([]string{
		"i|inc|increment+",
		"i|int=i",
	})

	// Output:
	// invalid option spec: -i is not unique from i|int
}

func ExampleNewDirectAssignParser() {
	var increment, intValue int64
	var stringList = make([]string, 0)
	var intList = make([]int64, 0)
	var floatList = make([]float32, 0)
	var stringValue string
	var floatValue float64
	var bool bool

	// After calling one of the Process routines the variable references passed in
	// will have the parsed option value directly assigned.
	op := NewDirectAssignParser(map[string]interface{}{
		// Allow for repeated `-i` or `--inc` or `--increment` options.  Each one of
		// the aliases is repeated the value is incrased by one.
		"inc|increment+": &increment,

		// Allow for `-S string` or `--string-list string` options.  The string
		// values will be stored in a slice in order of appearance.
		"S|string-list=s@": &stringList,

		// Allow for `-I 123` or `--int-list 123` options.  The int values will
		// be stored in a slice in order of appearance
		"I|int-list=i@": &intList,

		// Allow for `-F 1.23` or `--float-list 1.23` options.  The float values
		// will be stored in a slice in order of appearance
		"F|float-list=f@": &floatList,

		// Allow for `-s string` or `--string-value string`.
		"s|string-value=s": &stringValue,

		// Allow for `-i 123` or `--int-value 123`.
		"i|int-value=i": &intValue,

		// Allow for `-f 1.23` or `--float-value 1.23`.
		"f|float-value=f": &floatValue,

		// Allow for `-b` or `--bool`.
		"b|bool": &bool,
	})

	args := []string{
		"--inc",
		"--increment",
		"-S", "A",
		"--string-list", "B",
		"-I", "1",
		"--int-list", "2",
		"-F", ".1",
		"--float-list", ".2",
		"-s", "hey",
		"-i", "42",
		"-f", "3.141593",
		"--bool",
	}

	if err := op.ProcessAll(args); err != nil {
		panic(err)
	}

	fmt.Printf("increment: %d\n", increment)
	fmt.Printf("string-list: %v\n", stringList)
	fmt.Printf("int-list: %v\n", intList)
	fmt.Printf("float-list: %v\n", floatList)
	fmt.Printf("string-value: %s\n", stringValue)
	fmt.Printf("int-value: %d\n", intValue)
	fmt.Printf("float-value: %f\n", floatValue)
	fmt.Printf("bool: %t\n", bool)

	// Output:
	// increment: 2
	// string-list: [A B]
	// int-list: [1 2]
	// float-list: [0.1 0.2]
	// string-value: hey
	// int-value: 42
	// float-value: 3.141593
	// bool: true
}

func ExampleNewDirectAssignParser_callbacks() {

	usage := func() {
		fmt.Println(`
Usage: <appname> --help ...
`)
	}

	stuff := make(map[string]interface{})
	mapper := func(name string, value interface{}) {
		stuff[name] = value
	}

	list := make([]interface{}, 0)
	appender := func(value interface{}) {
		list = append(list, value)
	}

	op := NewDirectAssignParser(map[string]interface{}{
		"h|help":   usage,
		"o|opt=s":  mapper,
		"i|item=i": appender,
		"f|flag":   mapper,
		"m|more=s": appender,
	})

	args := []string{
		"-h",
		"--opt", "value",
		"-i", "123",
		"--flag",
		"-m", "more",
		"--item", "42",
	}

	if err := op.ProcessAll(args); err != nil {
		panic(err)
	}

	fmt.Printf("stuff[opt] = %s\n", stuff["opt"])
	fmt.Printf("stuff[flag] = %t\n", stuff["flag"])
	fmt.Printf("list: %v\n", list)

	// Output:
	// Usage: <appname> --help ...
	//
	// stuff[opt] = value
	// stuff[flag] = true
	// list: [123 more 42]
}

func ExampleOptionParser_ProcessAll() {
	op := NewParser([]string{
		"v|verbose+",
	})

	args := []string{
		"-v",
		"extra",
	}

	// No extranous options, just an unparsed argument, so no error
	if err := op.ProcessAll(args); err != nil {
		panic(err)
	}

	fmt.Printf("verbose: %d\n", op.Results["verbose"])
	fmt.Printf("unparsed args: %v\n", op.Args)

	op = NewParser([]string{
		"v|verbose+",
	})

	args = []string{
		"-v",
		"extra",
		"--bogus",
	}

	// This will error and panic with `Unknown option: --bogus`
	if err := op.ProcessAll(args); err != nil {
		fmt.Println(err)
	}

	// Output:
	// verbose: 1
	// unparsed args: [extra]
	// Unknown option: --bogus
}

func ExampleOptionParser_ProcessSome() {
	op := NewParser([]string{
		"v|verbose+",
	})

	args := []string{
		"-v",
		"--bogus",
		"extra",
	}

	// No error on unknown --bogus option
	if err := op.ProcessSome(args); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("verbose: %d\n", op.Results["verbose"])
	fmt.Printf("unparsed args: %v\n", op.Args)

	// Output:
	// verbose: 1
	// unparsed args: [--bogus extra]
}
