package eip

import (
	"gitee.com/ziIoT/ethernet-ip/types"
)

const (
	defaultPort        uint16      = 0xAF12
	defaultTimeTick    types.USINT = 3
	defaultTimeTickOut types.USINT = 250
)

type Config struct {
	TCPPort     uint16
	UDPPort     uint16
	Slot        uint8
	TimeTick    types.USINT
	TimeTickOut types.USINT
}

func DefaultConfig() *Config {
	return &Config{
		TCPPort:     defaultPort,
		UDPPort:     defaultPort,
		Slot:        0,
		TimeTick:    defaultTimeTick,
		TimeTickOut: defaultTimeTickOut,
	}
}
