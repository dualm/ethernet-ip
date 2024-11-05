package utils

import (
	"gitee.com/ziIoT/ethernet-ip/types"
)

func Len(raw []byte) types.USINT {
	l := len(raw)
	if l%2 == 1 {
		l += 1
	}

	return types.USINT(l / 2)
}
