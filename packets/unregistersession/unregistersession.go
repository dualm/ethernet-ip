package unregistersession

import (
	"github.com/dualm/ethernet-ip/packets"
	"github.com/dualm/ethernet-ip/packets/command"
	"github.com/dualm/ethernet-ip/types"
)

func New(session types.UDINT, context types.ULINT) (*packets.EncapsulationMessagePackets, error) {
	return &packets.EncapsulationMessagePackets{
		Header: packets.EncapsulationHeader{
			Command:       command.UnRegisterSession,
			Length:        0,
			SessionHandle: session,
			Status:        0,
			SenderContext: context,
			Options:       0,
		},
		SpecificData: nil,
	}, nil
}
