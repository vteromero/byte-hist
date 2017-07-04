package bytehist

import "testing"

func checkInitialized(bh *ByteHistogram, t *testing.T) {
	if bh.DataSize != 0 {
		t.Errorf("expected DataSize to be 0; got %v instead", bh.DataSize)
	}

	for _, c := range bh.Count {
		if c > 0 {
			t.Errorf("expected Count value to be 0; got %v instead", c)
		}
	}
}

func TestNewByteHistogram(t *testing.T) {
	bh := NewByteHistogram()

	checkInitialized(bh, t)
}

func TestByteHistogram_Init(t *testing.T) {
	bh := NewByteHistogram()
	bh.DataSize = 10000
	bh.Count[0] = 10
	bh.Count[1] = 20
	bh.Count[2] = 30

	bh.Init()

	checkInitialized(bh, t)
}

func TestByteHistogram_Update(t *testing.T) {
	var (
		testbytes = [][]byte{
			{0, 0, 0, 0, 0, 0, 0},
			{1, 1, 2, 2, 3, 3, 4},
			{100, 101, 102, 103, 104, 105, 106},
			{111, 111, 111, 112, 112, 112, 112},
		}
		expectedCount = [256]uint64{
			7, 2, 2, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 3,
			4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}
		expectedDataSize uint64 = 28
	)

	bh := NewByteHistogram()

	for _, arr := range testbytes {
		bh.Update(arr)
	}

	if bh.DataSize != expectedDataSize {
		t.Errorf("DataSize value %v does not match the expected %v",
			bh.DataSize, expectedDataSize)
	}

	if bh.Count != expectedCount {
		t.Errorf("Count array does not match the expected")
	}
}

func compareByteSlices(slice, expected []byte, t *testing.T) {
	if len(slice) != len(expected) {
		t.Fatalf("the slice length (%v) does not match the expected (%v)",
			len(slice), len(expected))
	}

	for i := range slice {
		if slice[i] != expected[i] {
			t.Errorf("the value %v does not match the expected %v at position %v",
				slice[i], expected[i], i)
		}
	}
}

func compareUint64Slices(slice, expected []uint64, t *testing.T) {
	if len(slice) != len(expected) {
		t.Fatalf("the slice length (%v) does not match the expected (%v)",
			len(slice), len(expected))
	}

	for i := range slice {
		if slice[i] != expected[i] {
			t.Errorf("the value %v does not match the expected %v at position %v",
				slice[i], expected[i], i)
		}
	}
}

func TestByteHistogram_ByteList(t *testing.T) {
	var (
		testbytes = []byte{
			10, 10, 10, 2, 2, 99, 99, 100, 67, 203,
			2, 99, 1, 207, 228, 13, 99, 2, 100, 177,
		}
		expectedByteList  = []byte{1, 2, 10, 13, 67, 99, 100, 177, 203, 207, 228}
		expectedByteCount = []uint64{1, 4, 3, 1, 1, 4, 2, 1, 1, 1, 1}
	)

	bh := NewByteHistogram()
	bh.Update(testbytes)
	bytelist, bytecount := bh.ByteList()

	compareByteSlices(bytelist, expectedByteList, t)
	compareUint64Slices(bytecount, expectedByteCount, t)
}

func TestByteHistogram_SortedByteList_ascOrder(t *testing.T) {
	var (
		testbytes = []byte{
			10, 10, 10, 2, 2, 99, 99, 100, 67, 203,
			2, 99, 1, 207, 228, 13, 99, 2, 100, 177,
		}
		expectedByteList  = []byte{1, 13, 67, 177, 203, 207, 228, 100, 10, 2, 99}
		expectedByteCount = []uint64{1, 1, 1, 1, 1, 1, 1, 2, 3, 4, 4}
	)

	bh := NewByteHistogram()
	bh.Update(testbytes)
	bytelist, bytecount := bh.SortedByteList(true)

	compareByteSlices(bytelist, expectedByteList, t)
	compareUint64Slices(bytecount, expectedByteCount, t)
}

func TestByteHistogram_SortedByteList_descOrder(t *testing.T) {
	var (
		testbytes = []byte{
			10, 10, 10, 2, 2, 99, 99, 100, 67, 203,
			2, 99, 1, 207, 228, 13, 99, 2, 100, 177,
		}
		expectedByteList  = []byte{2, 99, 10, 100, 1, 13, 67, 177, 203, 207, 228}
		expectedByteCount = []uint64{4, 4, 3, 2, 1, 1, 1, 1, 1, 1, 1}
	)

	bh := NewByteHistogram()
	bh.Update(testbytes)
	bytelist, bytecount := bh.SortedByteList(false)

	compareByteSlices(bytelist, expectedByteList, t)
	compareUint64Slices(bytecount, expectedByteCount, t)
}
