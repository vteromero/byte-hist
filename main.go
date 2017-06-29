// Copyright (c) 2017 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/vteromero/byte-hist/bytehist"
)

const VERSION = "0.1"

const (
	formatBinary = iota
	formatDecimal
	formatHexadecimal
)

type config struct {
	File       *os.File
	ByteFormat int
}

func newConfig() config {
	return config{nil, formatDecimal}
}

func byteFormat(formatStr string) (int, error) {
	switch formatStr {
	case "b":
		return formatBinary, nil
	case "d":
		return formatDecimal, nil
	case "x":
		return formatHexadecimal, nil
	default:
		return 0, errors.New("invalid byte format")
	}
}

func openFile(fileName string) (*os.File, error) {
	// if there is no fileName, we read from stdin
	if len(fileName) == 0 {
		// checking that stdin is not a terminal
		fi, err := os.Stdin.Stat()
		if err != nil {
			return nil, err
		}

		if fi.Mode()&os.ModeCharDevice != 0 {
			return nil, errors.New("the data won't be read from a terminal")
		}

		return os.Stdin, nil
	}

	// only read regular files
	fi, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}

	mode := fi.Mode()
	if !mode.IsRegular() {
		return nil, fmt.Errorf("'%s' is not a regular file", fileName)
	}

	return os.Open(fileName)
}

func closeFile(file *os.File) {
	// do not close stdin
	if file == os.Stdin {
		return
	}

	if err := file.Close(); err != nil {
		panic(err)
	}
}

func readBytes(file *os.File, bhist *bytehist.ByteHistogram) {
	buf := make([]byte, 1024)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		bhist.Update(buf, n)
	}
}

func doByteHistogram(cfg config) error {
	bhist := bytehist.NewByteHistogram()

	readBytes(cfg.File, bhist)

	if bhist.DataSize == 0 {
		return errors.New("no data to process")
	}

	printByteHistogram(cfg, bhist)

	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Println("usage: byte-hist [-help] [-format={d|x|b}] [FILE]")
		flag.PrintDefaults()
	}

	helpPtr := flag.Bool("help", false, "print this message")
	versionPtr := flag.Bool("version", false, "print the version")
	formatPtr := flag.String("format", "d", "byte format {\"d\"ecimal | he\"x\"adecimal | \"b\"inary}")

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		os.Exit(0)
	}

	if *versionPtr {
		fmt.Printf("byte-hist %s\n", VERSION)
		os.Exit(0)
	}

	bytefmt, err := byteFormat(*formatPtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if flag.NArg() > 1 {
		fmt.Println("too many arguments")
		os.Exit(1)
	}

	file, err := openFile(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer closeFile(file)

	cfg := newConfig()
	cfg.File = file
	cfg.ByteFormat = bytefmt

	if err := doByteHistogram(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
