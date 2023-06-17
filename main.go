package main

import (
	"github.com/cloudflare/circl/xof/k12"

	"bufio"
	"crypto/subtle"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

var (
	check = flag.Bool("c", false, "read checksum from FILE and check them")

	h   k12.State
	buf []byte
)

// Hash file. If in check mode, expected contains the desired hash.
func process(file string, expected []byte) bool {
	var (
		rd  io.ReadCloser
		err error
	)

	h.Reset()

	if file == "-" {
		rd = os.Stdin
	} else {
		rd, err = os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: open(): %v", file, err)
			return false
		}

		defer func() {
			rd.Close()
		}()
	}

	for {
		n := h.NextWriteSize()
		n, err = rd.Read(buf[:n])
		if err == io.EOF {
			h.Write(buf[:n])
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: read(): %v", file, err)
			return false
		}
	}

	var hash [32]byte
	h.Read(hash[:])
	if expected == nil {
		fmt.Printf("%x  %s\n", hash, file)
		return true
	}

	if subtle.ConstantTimeCompare(hash[:], expected) == 1 {
		fmt.Printf("%s: OK\n", file)
		return true
	}

	fmt.Printf("%s: FAILED\n", file)
	return false
}

// Parses and checks a checkSums file.
func checkSums(file string) bool {
	var (
		rd       io.ReadCloser
		err      error
		badLines int
	)
	ok := true

	if file == "-" {
		rd = os.Stdin
	} else {
		rd, err = os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: open(): %v", file, err)
			return false
		}

		defer func() {
			rd.Close()
		}()
	}

	scanner := bufio.NewScanner(rd)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}
		bits := strings.SplitN(line, "  ", 2)
		if len(bits) != 2 {
			badLines++
			continue
		}
		if len(bits[0]) != 64 {
			badLines++
			continue
		}
		expected, err := hex.DecodeString(bits[0])
		if err != nil {
			badLines++
			continue
		}
		ret := process(bits[1], expected)
		ok = ok && ret
	}

	if badLines > 0 {
		fmt.Fprintf(os.Stderr, "%s: %d line is improperly formatted\n",
			file, badLines)
		ok = false
	}

	if err = scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v", file, err)
		ok = false
	}

	return ok
}

func run() int {
	h = k12.NewDraft10(
		k12.WithWorkers(runtime.NumCPU()),
	)

	buf = make([]byte, h.MaxWriteSize())
	ok := true

	files := flag.Args()
	if len(files) == 0 {
		files = []string{"-"}
	}
	for _, file := range files {
		if *check {
			ret := checkSums(file)
			ok = ok && ret
			continue
		}

		ret := process(file, nil)
		ok = ok && ret
	}

	if !ok {
		return 1
	}
	return 0
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: [FILE]...\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	code := run()

	if code != 0 {
		os.Exit(code)
	}
}
