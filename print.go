// Copyright (c) 2017 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/vteromero/byte-hist/bytehist"
)

const (
	colIdxByte = iota
	colIdxCount
	colIdxPercentage
	colIdxSumPercentage
	colIdxHistogram
)

const (
	blackSquareChar            = '\u25a0'
	blackVerticalRectangleChar = '\u25ae'
)

const numColumns = 5

var (
	colWidth       = [numColumns]int{6, 15, 6, 6, 32}
	colValueFmt    = [numColumns]string{"%d", "%d", "%.2f", "%.2f", "%-30s"}
	histogramWidth = colWidth[colIdxHistogram] - 2
)

func wrappingFormat(width int, rightAlignment bool) string {
	sign := "-"
	if rightAlignment {
		sign = ""
	}
	return fmt.Sprintf("%%%s%ds", sign, width)
}

func printSummary(filename string, datasize uint64, bytelistlen int) {
	fmt.Println()
	fmt.Printf("%-20s%s\n", "File name:", filename)
	fmt.Printf("%-20s%d\n", "File size:", datasize)
	fmt.Printf("%-20s%d\n", "Different bytes:", bytelistlen)
	fmt.Println()
}

func printTableHeader() {
	fmt.Printf(wrappingFormat(colWidth[colIdxByte], true), "byte")
	fmt.Printf(wrappingFormat(colWidth[colIdxCount], true), "count")
	fmt.Printf(wrappingFormat(colWidth[colIdxPercentage], true), "rate")
	fmt.Printf(wrappingFormat(colWidth[colIdxSumPercentage], true), "sum")
	fmt.Printf(wrappingFormat(colWidth[colIdxHistogram], false), "  hist")
	fmt.Println()

	fmt.Printf("%s\n", strings.Repeat("=", colWidth[colIdxByte]+
		colWidth[colIdxCount]+colWidth[colIdxPercentage]+
		colWidth[colIdxSumPercentage]+colWidth[colIdxHistogram]+1))
}

func printTableBody(bytelist []byte, bytecount []uint64, datasize uint64) {
	maxcount := bytecount[0]
	for _, c := range bytecount[1:] {
		if c > maxcount {
			maxcount = c
		}
	}

	mincount := bytecount[0]
	for _, c := range bytecount[1:] {
		if c < mincount {
			mincount = c
		}
	}

	var sum uint64

	for i, b := range bytelist {
		percentage := float64(bytecount[i]) / float64(datasize)

		sum += bytecount[i]
		sumPercentage := float64(sum) / float64(datasize)

		normCount := 1.0
		if maxcount > mincount {
			normCount = float64(bytecount[i]-mincount) / float64(maxcount-mincount)
		}

		histLenFloat := normCount * float64(histogramWidth)
		histLenFloatRem := histLenFloat - math.Floor(histLenFloat)
		histStr := strings.Repeat(string(blackSquareChar), int(histLenFloat))

		if histLenFloatRem > 0.1 {
			histStr += string(blackVerticalRectangleChar)
		} else if histLenFloatRem > 0.6 {
			histStr += string(blackSquareChar)
		}

		fmt.Printf(wrappingFormat(colWidth[colIdxByte], true),
			fmt.Sprintf(colValueFmt[colIdxByte], b))
		fmt.Printf(wrappingFormat(colWidth[colIdxCount], true),
			fmt.Sprintf(colValueFmt[colIdxCount], bytecount[i]))
		fmt.Printf(wrappingFormat(colWidth[colIdxPercentage], true),
			fmt.Sprintf(colValueFmt[colIdxPercentage], percentage))
		fmt.Printf(wrappingFormat(colWidth[colIdxSumPercentage], true),
			fmt.Sprintf(colValueFmt[colIdxSumPercentage], sumPercentage))
		fmt.Printf(wrappingFormat(colWidth[colIdxHistogram], true),
			fmt.Sprintf(colValueFmt[colIdxHistogram], histStr))
		fmt.Println()
	}
}

func printByteHistogram(cfg config, bhist *bytehist.ByteHistogram) {
	switch cfg.ByteFormat {
	case formatBinary:
		colWidth[colIdxByte] = 10
		colValueFmt[colIdxByte] = "%08b"
	case formatHexadecimal:
		colWidth[colIdxByte] = 6
		colValueFmt[colIdxByte] = "%02x"
	case formatCharacter:
		colWidth[colIdxByte] = 10
		colValueFmt[colIdxByte] = "%q"
	}

	var (
		bytelist  []byte
		bytecount []uint64
	)

	switch cfg.SortOrder {
	case sortOrderNone:
		bytelist, bytecount = bhist.ByteList()
	case sortOrderAscending:
		bytelist, bytecount = bhist.SortedByteList(true)
	case sortOrderDescending:
		bytelist, bytecount = bhist.SortedByteList(false)
	}

	printSummary(cfg.File.Name(), bhist.DataSize, len(bytelist))

	printTableHeader()

	printTableBody(bytelist, bytecount, bhist.DataSize)
}
