package sendunitdata

import (
	"gitee.com/ziIoT/ethernet-ip/packets"
	"gitee.com/ziIoT/ethernet-ip/packets/command"
	"gitee.com/ziIoT/ethernet-ip/types"
)

func New(session types.UDINT, context types.ULINT, cpf *packets.CommandPacketFormat) (*packets.EncapsulationMessagePackets, error) {
	specificData := packets.SpecificData{
		InterfaceHandle: 0,
		Timeout:         0,
		Packet:          cpf,
	}

	specificDataBytes, err := specificData.Encode()
	if err != nil {
		return nil, err
	}

	return &packets.EncapsulationMessagePackets{
		Header: packets.EncapsulationHeader{
			Command:       command.SendUnitData,
			Length:        types.UINT(len(specificDataBytes)),
			SessionHandle: session,
			Status:        0,
			SenderContext: context,
			Options:       0,
		},
		SpecificData: specificDataBytes,
	}, nil
}

func Decode(raw *packets.EncapsulationMessagePackets) (*packets.SpecificData, error) {
	result := new(packets.SpecificData)

	err := result.Decode(raw.SpecificData)
	if err != nil {
		return nil, err
	}

	return result, nil
}
