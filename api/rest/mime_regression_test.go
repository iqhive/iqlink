package rest

import (
	"errors"
	"reflect"
	"testing"
)

func TestPlainTextPointerCycle(t *testing.T) {
	var cycle any
	cycle = &cycle
	if _, err := plainText(cycle); !errors.Is(err, errValueTooDeep) {
		t.Errorf("plainText returned %v", err)
	}
	if hasPlainForm(reflect.ValueOf(cycle)) {
		t.Error("a pointer cycle has no plain text form")
	}
	if kind := resultKind(reflect.ValueOf(cycle)); kind != reflect.Invalid {
		t.Errorf("resultKind returned %v", kind)
	}
}

func TestPlainTextRecursivePointer(t *testing.T) {
	type pointer *pointer
	var cycle pointer
	cycle = &cycle
	if _, err := plainText(cycle); !errors.Is(err, errValueTooDeep) {
		t.Errorf("plainText returned %v", err)
	}
	if hasPlainForm(reflect.ValueOf(cycle)) {
		t.Error("a pointer cycle has no plain text form")
	}
	if hasPlainForm(reflect.ValueOf([]pointer{cycle})) {
		t.Error("a slice of recursive pointers has no plain text form")
	}
	if kind := resultKind(reflect.ValueOf(cycle)); kind != reflect.Invalid {
		t.Errorf("resultKind returned %v", kind)
	}
}
