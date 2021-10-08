package packets

import (
	"github.com/dualm/ethernet-ip/bufferEip"
	"github.com/dualm/ethernet-ip/types"
)

type SpecificData struct {
	InterfaceHandle types.UDINT
	Timeout         types.UINT
	Packet          *CommandPacketFormat
}

func (data SpecificData)Encode() ([]byte, error) {
	buffer := bufferEip.New(nil)

	buffer.WriteLittle(data.InterfaceHandle)
	buffer.WriteLittle(data.Timeout)
	packet, err := data.Packet.Encode()
	if err != nil {
		return nil, err
	}
	buffer.WriteLittle(packet)

	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
	
}

func (data *SpecificData)Decode(raw []byte) error {
	buffer := bufferEip.New(raw)
	buffer.ReadLittle(&data.InterfaceHandle)
	buffer.ReadLittle(&data.Timeout)
	data.Packet = new(CommandPacketFormat)
	data.Packet.Decode(buffer)

	if err := buffer.Error(); err != nil {
		return err 
	}

	return nil
}