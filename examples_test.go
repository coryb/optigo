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

	// Prints: verbose: 3
	fmt.Printf("verbose: %d\n", op.Results["verbose"])

	// Prints: unparsed args: [extra]
	fmt.Printf("unparsed args: %v\n", op.Args)
}

func ExampleNewParser() {
	// Note that all values will be stored in OptionParser.Results after a Process function
	// is called.  The Result key will be stored as the last alias.
	op := NewParser([]string{
		// Allow for repeated `-i` or `--inc` or `--increment` options.  Each one of
		// the aliases is repeated the value is incrased by one.
		"i|inc|increment+",

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
		"-i",
		"--increment",
		"-S", "A",
		"--string-list", "B",
		"-I", "1",
		"--int-list", "2",
		"-F", ".1",
		"--float-list", ".2",
		"-s", "hey",
		"-i", "42",
		"-f", "3.141592653589793238462643",
		"--bool",
	}

	if err := op.ProcessAll(args); err != nil {
		panic(err)
	}

	// Prints: increment: 2
	fmt.Printf("increment: %d\n", op.Results["increment"])

	// Prints: string-list: [A B]
	fmt.Printf("string-list: %v\n", op.Results["string-list"])

	// Prints: int-list: [1 2]
	fmt.Printf("int-list: %v\n", op.Results["int-list"])

	// Prints: float-list: [0.1 0.2]
	fmt.Printf("float-list: %v\n", op.Results["float-list"])

	// Prints: string-value: hey
	fmt.Printf("string-value: %s\n", op.Results["string-value"])

	// Prints: int-value: 42
	fmt.Printf("int-value: %d\n", op.Results["int-value"])

	// Prints: float-value: 3.141592653589793238462643
	fmt.Printf("float-value: %f\n", op.Results["float-value"])

	// Prints: bool: true
	fmt.Printf("bool: %t\n", op.Results["bool"])
}

func ExampleNewDirectAssignParser() {
	var increment, intValue int
	var stringList = make([]string, 0)
	var intList = make([]int, 0)
	var floatList = make([]float32, 0)
	var stringValue string
	var floatValue float64
	var bool bool

	// After calling one of the Process routines the variable references passed in
	// will have the parsed option value directly assigned.
	op := NewDirectAssignParser(map[string]interface{}{
		// Allow for repeated `-i` or `--inc` or `--increment` options.  Each one of
		// the aliases is repeated the value is incrased by one.
		"i|inc|increment+": &increment,

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
		"-i",
		"--increment",
		"-S", "A",
		"--string-list", "B",
		"-I", "1",
		"--int-list", "2",
		"-F", ".1",
		"--float-list", ".2",
		"-s", "hey",
		"-i", "42",
		"-f", "3.141592653589793238462643",
		"--bool",
	}

	if err := op.ProcessAll(args); err != nil {
		panic(err)
	}

	// Prints: increment: 2
	fmt.Printf("increment: %d\n", increment)

	// Prints: string-list: [A B]
	fmt.Printf("string-list: %v\n", stringList)

	// Prints: int-list: [1 2]
	fmt.Printf("int-list: %v\n", intList)

	// Prints: float-list: [0.1 0.2]
	fmt.Printf("float-list: %v\n", floatList)

	// Prints: string-value: hey
	fmt.Printf("string-value: %s\n", stringValue)

	// Prints: int-value: 42
	fmt.Printf("int-value: %d\n", intValue)

	// Prints: float-value: 3.141592653589793238462643
	fmt.Printf("float-value: %f\n", floatValue)

	// Prints: bool: true
	fmt.Printf("bool: %t\n", bool)
}

func ExampleProcessAll() {
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

	// Prints: verbose: 1
	fmt.Printf("verbose: %d\n", op.Results["verbose"])

	// Prints: unparsed args: [extra]
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
		panic(err)
	}
}

func ExampleProcessSome() {
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
		panic(err)
	}

	// Prints: verbose: 1
	fmt.Printf("verbose: %d\n", op.Results["verbose"])

	// Prints: unparsed args: [--bogus extra]
	fmt.Printf("unparsed args: %v\n", op.Args)
}
