// Copyright (c) 2017 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bytehist

import "sort"

// ByteHistogram is an structure that holds the information of a byte histogram.
type ByteHistogram struct {
	Count    [256]uint64
	DataSize uint64
}

// NewByteHistogram creates a new ByteHistogram.
func NewByteHistogram() *ByteHistogram {
	return &ByteHistogram{}
}

// Init initializes a ByteHistogram.
func (bh *ByteHistogram) Init() {
	for i := 0; i < 256; i++ {
		bh.Count[i] = 0
	}
	bh.DataSize = 0
}

// Update updates a ByteHistogram with an array of bytes.
func (bh *ByteHistogram) Update(bytes []byte) {
	for _, b := range bytes {
		bh.Count[b]++
	}
	bh.DataSize += uint64(len(bytes))
}

// ByteList returns two values: a slice of the bytes that have been counted
// once at least and a slice with the actual number of times that every byte
// appears on the processed data.
func (bh *ByteHistogram) ByteList() ([]byte, []uint64) {
	bytelist := make([]byte, 256)
	bytecount := make([]uint64, 256)
	listlen := 0

	for i, c := range bh.Count {
		if c > 0 {
			bytelist[listlen] = byte(i)
			bytecount[listlen] = uint64(c)
			listlen++
		}
	}

	return bytelist[0:listlen], bytecount[0:listlen]
}

type byteCountPair struct {
	b byte
	c uint64
}

type byCountAsc []byteCountPair

func (a byCountAsc) Len() int      { return len(a) }
func (a byCountAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCountAsc) Less(i, j int) bool {
	return (a[i].c < a[j].c) || (a[i].c == a[j].c && a[i].b < a[j].b)
}

type byCountDesc []byteCountPair

func (a byCountDesc) Len() int      { return len(a) }
func (a byCountDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCountDesc) Less(i, j int) bool {
	return (a[i].c > a[j].c) || (a[i].c == a[j].c && a[i].b < a[j].b)
}

// SortedByteList returns two values as the ByteList function does, but the
// resulting slices are sorted by the number of bytes.
// The sorting order is specified by ascOrder, that will be ascending if
// the param is true or descending if it is false.
func (bh *ByteHistogram) SortedByteList(ascOrder bool) ([]byte, []uint64) {
	pairs := make([]byteCountPair, 256)

	for i, count := range bh.Count {
		pairs[i] = byteCountPair{b: byte(i), c: count}
	}

	if ascOrder {
		sort.Sort(byCountAsc(pairs))
	} else {
		sort.Sort(byCountDesc(pairs))
	}

	bytelist := make([]byte, 256)
	bytecount := make([]uint64, 256)
	listlen := 0

	for _, pair := range pairs {
		if pair.c > 0 {
			bytelist[listlen] = pair.b
			bytecount[listlen] = pair.c
			listlen++
		}
	}

	return bytelist[0:listlen], bytecount[0:listlen]
}
