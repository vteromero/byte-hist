// Copyright (c) 2017 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bytehist

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
func (bh *ByteHistogram) Update(b []byte, l int) {
	for i := 0; i < l; i++ {
		bh.Count[b[i]]++
	}
	bh.DataSize += uint64(l)
}

// ByteList returns two values: an array of the bytes that have been counted
// once at least and an array with the actual number of times that every byte
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
