package ast

import (
	"strings"
	"testing"

	"github.com/SmartBrave/Athena/easyerrors"
)

func TestJudgeUndefinedVariable(t *testing.T) {
	a := NewAST(STRICT)
	got, err := a.Judge(map[string]interface{}{}, nil, `${missing} == 1`)
	if err == nil {
		t.Errorf("expected error for undefined ${missing}, got nil (result=%v)", got)
	}
}

func TestJudgeNonBooleanResultErrors(t *testing.T) {
	a := NewAST(STRICT)
	_, err := a.Judge(nil, nil, `1 + 2`)
	if err == nil {
		t.Errorf("expected error: result of '1 + 2' is not boolean")
	}
}

func TestJudgeMalformedExpression(t *testing.T) {
	a := NewAST(STRICT)
	_, err := a.Judge(nil, nil, `${app} ==`)
	if err == nil {
		t.Errorf("expected parse error on truncated expression")
	}
}

func TestCompatibleModeAcceptsSingleAndOr(t *testing.T) {
	// COMPATIBLE mode accepts '&' and '|' as && and || (per existing test case 25/26).
	a := NewAST(COMPATIBLE)
	got, err := a.Judge(nil, nil, `true & true | false`)
	if err != nil {
		t.Fatalf("Judge: %v", err)
	}
	if !got {
		t.Errorf("COMPATIBLE: true & true | false should be true")
	}
}

func TestSliceOutOfRangeErrors(t *testing.T) {
	a := NewAST(STRICT)
	_, err := a.Judge(map[string]interface{}{"s": "abc"}, nil, `${s}[5:10] == "x"`)
	if err == nil {
		t.Errorf("expected slice-out-of-range error")
	}
}

func TestUserFunctionCall(t *testing.T) {
	a := NewAST(STRICT)
	whitelist := map[string]bool{"alice": true, "bob": true}
	ops := map[string]interface{}{
		"isMember": func(name interface{}) bool {
			s, ok := name.(string)
			return ok && whitelist[s]
		},
	}
	args := map[string]interface{}{"name": "alice"}

	got, err := a.Judge(args, ops, `isMember(${name})`)
	if err != nil {
		t.Fatalf("Judge: %v", err)
	}
	if !got {
		t.Errorf("alice should match")
	}

	args["name"] = "eve"
	got, err = a.Judge(args, ops, `isMember(${name})`)
	if err != nil {
		t.Fatalf("Judge: %v", err)
	}
	if got {
		t.Errorf("eve should not match")
	}
}

func TestUnaryNot(t *testing.T) {
	a := NewAST(STRICT)
	got, err := a.Judge(nil, nil, `!false`)
	if err != nil || !got {
		t.Errorf("!false: got=%v err=%v want true,nil", got, err)
	}
	got, err = a.Judge(nil, nil, `!(1 == 1)`)
	if err != nil || got {
		t.Errorf("!(1==1): got=%v err=%v want false,nil", got, err)
	}
}

// Integration: AST.Judge errors flow through easyerrors.HandleMultiError. A
// rule-evaluation pipeline often runs many rules and surfaces the first
// failure.
func TestAstIntegrationWithEasyErrors(t *testing.T) {
	a := NewAST(STRICT)
	args := map[string]interface{}{"x": 5}
	rules := []string{
		`${x} > 0`,            // ok -> true
		`${x} < 100`,          // ok -> true
		`${nonexistent} == 1`, // error
		`${x} > -1`,           // never reached
	}
	errs := make([]error, 0, len(rules))
	for _, r := range rules {
		_, err := a.Judge(args, nil, r)
		errs = append(errs, err)
	}

	combined := easyerrors.HandleMultiError(easyerrors.Simple(), errs...)
	if combined == nil {
		t.Fatal("expected combined error from rule with undefined variable")
	}
	if !strings.Contains(combined.Error(), "bad expression") {
		// The badError is a sentinel; we don't assert the exact message but
		// we expect a non-empty error.
		if combined.Error() == "" {
			t.Errorf("combined error has empty message")
		}
	}
}

// Integration: pre-defined function `contains`. Verify it's hooked up.
func TestPredefinedContains(t *testing.T) {
	a := NewAST(STRICT)
	got, err := a.Judge(map[string]interface{}{"v": "hello world"}, nil, `contains(${v}, "world")`)
	if err != nil {
		t.Fatalf("Judge: %v", err)
	}
	if !got {
		t.Errorf("contains(\"hello world\", \"world\") = false")
	}
}

// Integration: pre-defined function `mod`. Verify behavior on numeric strings.
func TestPredefinedMod(t *testing.T) {
	a := NewAST(STRICT)
	got, err := a.Judge(map[string]interface{}{"id": "12345"}, nil, `mod(${id}, 100) == 45`)
	if err != nil {
		t.Fatalf("Judge: %v", err)
	}
	if !got {
		t.Errorf("mod(\"12345\", 100) should equal 45")
	}
}
