package utils

// DeltaUint64 is a fuction return |i-j| on a ring mod 2^64
// due to the type of i and j is unsined int64,
// so it's important to test whether one parameter is greater than the other one.
// 2.1 inpHID = 1, nodeHID = 5 => delta = 5 - 1 = 4
// 2.2 inpHID = 1, nodeHID = 2^64 - 2 => delta = 2^64 - 2 - 1 = 2^64 - 3 = 18446744073709551613
// 2.3. inpHID = 2^64 - 2, nodeHID = 1 => delta = (2^64-1) - ((2^64-2) - 1) + 1 = 3
// NOTICE: why "+1" ? because UINT64MAX = 2^64 - 1, instead of 2^64
// 2.4 inpHID = 5, nodeHID = 1 => delta = (2^64-1) - (5 - 1) + 1 = 2^64 - 4 = 18446744073709551612
func DeltaUint64(inpHID, nodeHID uint64) (uint64, error) {
	const (
		UINT64MAX     = ^uint64(0)
		HALFUINT64MAX = UINT64MAX >> 1
	)

	if inpHID < nodeHID {
		return nodeHID - inpHID, nil
	}
	return UINT64MAX - (inpHID - nodeHID) + 1, nil
}
