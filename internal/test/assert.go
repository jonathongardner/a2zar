package test

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"testing"
)

func AssertEqual[T comparable](t *testing.T, want T, got T, msg string) {
	if got != want {
		t.Errorf("%s: want %v, got %v", msg, want, got)
	}
}

func AssertErrors(t *testing.T, want []error, got error, msg string) {
	t.Helper()

	if len(want) == 0 {
		if got != nil {
			t.Errorf("%s: want no error, got %v", msg, got)
		}
		return
	}

	for i, err := range want {
		if !errors.Is(got, err) {
			t.Errorf("%s: %v want error %v, got %v", msg, i, err, got)
		}
	}
}

func AssertSha1(t *testing.T, expected string, r io.Reader, msg string) {
	t.Helper()

	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		t.Fatal(fmt.Errorf("failed to compute sha1: %v", err))
	}

	actual := fmt.Sprintf("%x", h.Sum(nil))
	AssertEqual(t, expected, actual, msg)
}
