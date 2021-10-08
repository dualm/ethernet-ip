package utils

import (
	"github.com/dualm/ethernet-ip/types"
	"math/rand"
	"time"
)

func GetNewContext() types.ULINT {
	time.Sleep(time.Nanosecond)
	rand.Seed(time.Now().UnixNano())
	return types.ULINT(rand.Int63())
}
