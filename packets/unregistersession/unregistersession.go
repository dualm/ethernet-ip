package unregistersession

import (
	"gitee.com/ziIoT/ethernet-ip/packets"
	"gitee.com/ziIoT/ethernet-ip/packets/command"
	"gitee.com/ziIoT/ethernet-ip/types"
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
