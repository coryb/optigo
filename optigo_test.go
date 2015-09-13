package optigo

import (
	"reflect"
	"testing"
)

func TestParseValue(t *testing.T) {
	// test fail to parse int from string
	o := option{"test", false, reflect.ValueOf(nil), atASSIGN, dtINTEGER}
	if _, err := o.parseValue("abc"); err == nil {
		t.Fail()
	}

	// test fail to parse float from string
	o.dataType = dtFLOAT
	if _, err := o.parseValue("abc"); err == nil {
		t.Fail()
	}

	// test unknown parse type
	o.dataType = 42
	if _, err := o.parseValue("abc"); err == nil {
		t.Fail()
	}
}

func TestBogusUsage(t *testing.T) {
	// this test will panic, so expect that
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	// when using @ you need to specify a type as well
	NewParser([]string{
		"many@",
	})
}

func TestBogusDirectAssignUsage(t *testing.T) {
	// this test will panic, so expect that
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	// when using @ you need to specify a type as well
	many := make([]string, 0)
	NewDirectAssignParser(map[string]interface{}{
		"many@": &many,
	})
}

func TestProcesFailure(t *testing.T) {
	op := NewParser([]string{
		"opt=i",
	})
	// bogus opt with no value
	args := []string{"--opt"}

	if err := op.ProcessAll(args); err == nil {
		t.Fail()
	}

	// bogus opt with string value instead of int
	args = []string{"--opt", "abc"}

	if err := op.ProcessAll(args); err == nil {
		t.Fail()
	}
}

func TestDashDash(t *testing.T) {
	foobar := false
	op := NewDirectAssignParser(map[string]interface{}{
		"foobar": &foobar,
	})

	args := []string{"--", "--foobar"}

	// this should process fine
	if err := op.ProcessSome(args); err != nil {
		t.Fail()
	}

	// foobar should not have been set
	if foobar {
		t.Fail()
	}

	// --foobar should be in unparsed args
	if len(op.Args) != 1 {
		t.Fail()
	}

	if op.Args[0] != "--foobar" {
		t.Fail()
	}

}

func TestSingleOptionNoArg(t *testing.T) {
	var foobar int64
	op := NewDirectAssignParser(map[string]interface{}{
		"foobar=i": &foobar,
	})

	args := []string{"--foobar="}

	// this should fail to process, missing argument
	if err := op.ProcessAll(args); err == nil {
		t.Fail()
	}
}

func TestSingleOptionInvalidArg(t *testing.T) {
	var foobar int64
	op := NewDirectAssignParser(map[string]interface{}{
		"foobar=i": &foobar,
	})

	args := []string{"--foobar=abc"}

	// this should fail to process, invalid argument
	if err := op.ProcessAll(args); err == nil {
		t.Fail()
	}
}
