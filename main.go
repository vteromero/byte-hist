// Copyright (c) 2017 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/vteromero/byte-hist/bytehist"
)

// VERSION holds the application version
const VERSION = "0.2"

const (
	formatBinary = iota
	formatDecimal
	formatHexadecimal
)

const (
	sortOrderNone = iota
	sortOrderAscending
	sortOrderDescending
)

type config struct {
	File       *os.File
	ByteFormat int
	SortOrder  int
}

func newConfig() config {
	return config{nil, formatDecimal, sortOrderNone}
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

func sortOrder(orderStr string) (int, error) {
	switch orderStr {
	case "":
		return sortOrderNone, nil
	case "asc":
		return sortOrderAscending, nil
	case "desc":
		return sortOrderDescending, nil
	default:
		return 0, errors.New("invalid sort order")
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

		bhist.Update(buf[0:n])
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

func usage() {
	fmt.Println(`
usage: byte-hist [-help] [-version] [-format={d|x|b}] [-sort={asc|desc}] [FILE]

options:`)
	flag.PrintDefaults()
	fmt.Println()
}

func main() {
	log.SetPrefix("byte-hist: ")

	flag.Usage = usage

	helpPtr := flag.Bool("help", false, "print this message")
	versionPtr := flag.Bool("version", false, "print the version")
	formatPtr := flag.String("format", "d", "byte format {\"d\"ecimal | he\"x\"adecimal | \"b\"inary}")
	sortPtr := flag.String("sort", "", "sort by count {\"asc\" | \"desc\"}")

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
		log.Fatal(err)
	}

	sortorder, err := sortOrder(*sortPtr)
	if err != nil {
		log.Fatal(err)
	}

	if flag.NArg() > 1 {
		log.Fatal("too many arguments")
	}

	file, err := openFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	defer closeFile(file)

	cfg := newConfig()
	cfg.File = file
	cfg.ByteFormat = bytefmt
	cfg.SortOrder = sortorder

	if err := doByteHistogram(cfg); err != nil {
		log.Fatal(err)
	}
}
