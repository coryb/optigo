package optigo

import "testing"

func TestArg(t *testing.T) {
	op := NewParser([]string{})
	if err := op.ProcessAll([]string{"arg"}); err != nil {
		t.Error(err)
	}
	if len(op.Args) != 1 {
		t.Fail()
	}
	if op.Args[0] != "arg" {
		t.Fail()
	}
}

func TestBool(t *testing.T) {
	op := NewParser([]string{"b|bool"})
	if err := op.ProcessAll([]string{}); err != nil {
		t.Error(err)
	}
	if op.Results["bool"] != false {
		t.Fail()
	}

	op = NewParser([]string{"b|bool"})
	if err := op.ProcessAll([]string{"--bool"}); err != nil {
		t.Error(err)
	}
	if op.Results["bool"] != true {
		t.Fail()
	}

	op = NewParser([]string{"b|bool"})
	if err := op.ProcessAll([]string{"-b"}); err != nil {
		t.Error(err)
	}
	if op.Results["bool"] != true {
		t.Fail()
	}
}

func TestString(t *testing.T) {
	op := NewParser([]string{"s|string=s"})
	if err := op.ProcessAll([]string{}); err != nil {
		t.Error(err)
	}

	if op.Results["string"] != "" {
		t.Fail()
	}

	op = NewParser([]string{"s|string=s"})
	if err := op.ProcessAll([]string{"--string", "strval"}); err != nil {
		t.Error(err)
	}
	if op.Results["string"] != "strval" {
		t.Fail()
	}

	op = NewParser([]string{"s|string=s"})
	if err := op.ProcessAll([]string{"-s", "strval"}); err != nil {
		t.Error(err)
	}
	if op.Results["string"] != "strval" {
		t.Fail()
	}
}
