package registersession

import (
	"gitee.com/ziIoT/common"
	"gitee.com/ziIoT/ethernet-ip/packets"
	"gitee.com/ziIoT/ethernet-ip/packets/command"
	"gitee.com/ziIoT/ethernet-ip/types"
)

func New(context types.ULINT) (*packets.EncapsulationMessagePackets, error) {
	data := RegisterSessionSpecificData{
		ProtocolVersion: 1,
		OptionsFlags:    0,
	}

	specificDataBytes, err := data.Encode()
	if err != nil {
		return nil, err
	}

	return &packets.EncapsulationMessagePackets{
		Header: packets.EncapsulationHeader{
			Command:       command.RegisterSession,
			Length:        4,
			SessionHandle: 0,
			Status:        0,
			SenderContext: context,
			Options:       0,
		},
		SpecificData: specificDataBytes,
	}, nil
}

type RegisterSessionSpecificData struct {
	ProtocolVersion types.UINT
	OptionsFlags    types.UINT
}

func (data *RegisterSessionSpecificData) Encode() ([]byte, error) {
	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle(data.ProtocolVersion)
	buffer.WriteLittle(data.OptionsFlags)
	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
