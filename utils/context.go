package utils

import (
	"math/rand"
	"time"

	"github.com/dualm/ethernet-ip/types"
)

func GetNewContext() types.ULINT {
	time.Sleep(time.Nanosecond)
	rand.Seed(time.Now().UnixNano())

	return types.ULINT(rand.Int63())
}
