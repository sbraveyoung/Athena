package easyerrors

import (
	"errors"
	"testing"
)

func TestSimpleHandlesNil(t *testing.T) {
	if !Simple()(nil) {
		t.Errorf("Simple()(nil) = false, want true")
	}
}

func TestSimpleHandlesNonNil(t *testing.T) {
	if Simple()(errors.New("boom")) {
		t.Errorf("Simple()(non-nil) = true, want false")
	}
}

func TestHandleMultiErrorAllNil(t *testing.T) {
	if err := HandleMultiError(Simple(), nil, nil, nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestHandleMultiErrorReturnsFirstNonNil(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")
	if err := HandleMultiError(Simple(), nil, first, second); err != first {
		t.Errorf("expected first error returned, got %v", err)
	}
}

func TestHandleMultiErrorEmpty(t *testing.T) {
	if err := HandleMultiError(Simple()); err != nil {
		t.Errorf("expected nil for empty errs, got %v", err)
	}
}

func TestHandleMultiErrorShortCircuits(t *testing.T) {
	calls := 0
	first := errors.New("stop")
	second := errors.New("never")
	handler := func(e error) bool {
		calls++
		return e == nil
	}
	got := HandleMultiError(handler, nil, first, second)
	if got != first {
		t.Errorf("expected first, got %v", got)
	}
	// handler should be invoked for nil + first, then short-circuit before second.
	if calls != 2 {
		t.Errorf("handler invoked %d times, want 2 (short-circuit)", calls)
	}
}

func TestHandleMultiErrorCustomHandler(t *testing.T) {
	// A handler that swallows a specific sentinel and treats everything else as
	// fatal. Demonstrates the "filter and short-circuit" use case.
	skip := errors.New("skip me")
	handler := func(e error) bool { return e == nil || errors.Is(e, skip) }
	bad := errors.New("real")
	if err := HandleMultiError(handler, nil, skip, skip, bad, errors.New("after")); err != bad {
		t.Errorf("expected %v, got %v", bad, err)
	}
	if err := HandleMultiError(handler, nil, skip, skip); err != nil {
		t.Errorf("expected all-skipped to return nil, got %v", err)
	}
}
