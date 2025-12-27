package global

import "errors"

var ErrEpochMismatch = errors.New("epoch mismatch across regions")

type ConsistencyMode uint8

const (
	Fast ConsistencyMode = iota
	Strict
)

func CheckEpochAlignment(epochs []uint64, mode ConsistencyMode) error {

	if mode == Fast {
		return nil
	}

	if len(epochs) == 0 {
		return nil
	}

	base := epochs[0]
	for _, e := range epochs[1:] {
		if e != base {
			return ErrEpochMismatch
		}
	}

	return nil
}
