package nop

import (
	"github.com/dualm/ethernet-ip/packets"
	"github.com/dualm/ethernet-ip/packets/command"
	"github.com/dualm/ethernet-ip/types"
)

func New(data []byte) (*packets.EncapsulationMessagePackets, error) {
	return &packets.EncapsulationMessagePackets{
		Header: packets.EncapsulationHeader{
			Command:       command.NOP,
			Length:        types.UINT(len(data)),
			SessionHandle: 0,
			Status:        0,
			SenderContext: 0,
			Options:       0,
		},
		// unused data
		SpecificData: data,
	}, nil
}
