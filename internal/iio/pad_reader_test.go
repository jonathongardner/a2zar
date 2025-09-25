package iio

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/jonathongardner/a2zar/internal/test"
)

type result struct {
	data   string
	err    error
	padErr error
}
type exp struct {
	name    string
	limit   int64
	padding int64
	data    string
	results []result
	remaing result
}

func TestPadReader(te *testing.T) {
	exps := []exp{
		{
			name: "read-exact-no-pad", limit: 4, padding: 2, data: "foo-lala", results: []result{
				{data: "foo-"}, {data: "lala"},
			},
		},
		{
			name: "read-exact-with-pad", limit: 4, padding: 3, data: "foo-lala-che", results: []result{
				{data: "foo-"}, {data: "la-c"},
			},
		},
		{
			name: "read-2nd-really-short", limit: 4, padding: 2, data: "foo-", results: []result{
				{data: "foo-"}, {err: io.ErrUnexpectedEOF},
			},
		},
		{
			name: "read-2nd-short", limit: 4, padding: 2, data: "foo-lal", results: []result{
				{data: "foo-"}, {data: "lal", err: io.ErrUnexpectedEOF},
			},
		},
		{
			name: "read-2nd-pad-really-short", limit: 4, padding: 3, data: "foo-lala-c", results: []result{
				{data: "foo-"}, {data: "la-c", padErr: io.ErrUnexpectedEOF},
			},
		},
		{
			name: "read-2nd-pad-short", limit: 4, padding: 3, data: "foo-lala-ch", results: []result{
				{data: "foo-"}, {data: "la-c", padErr: io.ErrUnexpectedEOF},
			},
		},
	}
	for _, exp := range exps {
		te.Run(fmt.Sprintf("%s-%d-%d", exp.name, exp.limit, exp.padding), func(t *testing.T) {
			r := strings.NewReader(exp.data)
			for i, res := range exp.results {
				p := NewLimitPadReader(r, exp.limit, exp.padding)
				assertPad(t, res, p, i)
			}
			assertPad(t, exp.remaing, nopPadReader{r}, "remaing")
		})
	}
}

type nopPadReader struct {
	io.Reader
}

func (n nopPadReader) Pad() error {
	return nil
}

func assertPad(t *testing.T, res result, r PadReader, arg any) {
	if res.err != nil && res.padErr != nil {
		panic("cant do that!")
	}
	data, err := io.ReadAll(r)
	test.AssertEqualF(t, res.data, string(data), "expected same data %v", arg)
	test.AssertEqualF(t, res.err, err, "expected same error %v", arg)
	if err == nil {
		test.AssertEqualF(t, res.padErr, r.Pad(), "expected same pad error %v", arg)
	}
	test.CheckError(t)
}
