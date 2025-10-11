package xar

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/jonathongardner/a2zar/internal/test"
)

type res struct {
	sha1 string
	exp  expXar
}

var emptySha1 = "da39a3ee5e6b4b0d3255bfef95601890afd80709"

func TestReadingXar(t *testing.T) {
	// 2025-07-02 22:11:16 +0000 UTC
	expTime := []time.Time{
		time.Date(2025, 7, 2, 22, 11, 16, 0, time.UTC),
		time.Date(2025, 7, 2, 22, 11, 29, 0, time.UTC),
	}
	expTar := []res{
		{
			sha1: "39822d840f7301585cd35561838c9a600baf3ef6",
			exp: expXar{
				name:    "Distribution",
				path:    "Distribution",
				mode:    0,
				size:    1812,
				modTime: BadTime,
				errs:    []error{ErrInvalidTime, ErrInvalidMode},
			},
		},
		{
			sha1: emptySha1,
			exp: expXar{
				name:    "Resources",
				path:    "Resources",
				mode:    0700 | os.ModeDir,
				size:    0,
				modTime: BadTime,
			},
		},
		{
			sha1: "f06425976709d5ccb1ccfa1acc8ce201d7395544",
			exp: expXar{
				name:    "background.png",
				path:    "Resources/background.png",
				mode:    0,
				size:    10293,
				modTime: BadTime,
				errs:    []error{ErrInvalidTime, ErrInvalidMode},
			},
		},
		{
			sha1: emptySha1,
			exp: expXar{
				name:    "org.golang.go.pkg",
				path:    "org.golang.go.pkg",
				mode:    0700 | os.ModeDir,
				size:    0,
				modTime: BadTime,
			},
		},
		{
			sha1: "1666a57b5bc7061b853416c1b7788c3a7c883861",
			exp: expXar{
				name:    "Bom",
				path:    "org.golang.go.pkg/Bom",
				mode:    0644,
				size:    4124144,
				modTime: expTime[0],
			},
		},
		{
			sha1: "789bea9d8af228dffe193de53ce2179e99893ccf",
			exp: expXar{
				name:    "Payload",
				path:    "org.golang.go.pkg/Payload",
				mode:    0644,
				size:    76144916,
				modTime: expTime[1],
			},
		},
		{
			sha1: "63915d57722c0ab98eaef75cd778edc9e8257fe1",
			exp: expXar{
				name:    "Scripts",
				path:    "org.golang.go.pkg/Scripts",
				mode:    0644,
				size:    300,
				modTime: expTime[1],
			},
		},
		{
			sha1: "2753eeeb8b5e480657834804212883a6e8d47825",
			exp: expXar{
				name:    "PackageInfo",
				path:    "org.golang.go.pkg/PackageInfo",
				mode:    0,
				size:    569,
				modTime: BadTime,
				errs:    []error{ErrInvalidTime},
			},
		},
	}
	file, err := os.Open(test.LargeFile("go1.24.5.darwin-arm64.pkg"))
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			t.Fatalf("failed to close file: %v", err)
		}
	}()

	xr, err := NewReader(file)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to create xar reader: %v", err))
	}

	count := 0
	for {
		fi, err := xr.Next()
		if err == io.EOF {
			break
		}

		msg := fmt.Sprintf("index: %d ", count)
		assertFileInfo(t, fi, expTar[count].exp, msg, err)

		expSha1 := expTar[count].sha1
		if expSha1 != emptySha1 {
			ex := fi.file.Data.ExtractedChecksum
			test.AssertEqual(t, "sha1", ex.Style, "should have sha1 in TOC")
			test.AssertEqual(t, expTar[count].sha1, ex.Text, "should have expected sha1 in TOC")
		}
		test.AssertSha1(t, expSha1, xr, msg)

		count++
	}
	test.AssertEqual(t, count, len(expTar), "should have read all expected files")
}
