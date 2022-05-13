package packets

import (
	"fmt"

	"github.com/dualm/common"
	"github.com/dualm/ethernet-ip/types"
)

type SpecificData struct {
	InterfaceHandle types.UDINT
	Timeout         types.UINT
	Packet          *CommandPacketFormat
}

func (data SpecificData) Encode() ([]byte, error) {
	buffer := common.NewEmptyBuffer()

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

func (data *SpecificData) Decode(raw []byte) error {
	buffer := common.NewBuffer(raw)
	buffer.ReadLittle(&data.InterfaceHandle)
	buffer.ReadLittle(&data.Timeout)
	data.Packet = new(CommandPacketFormat)

	if err := data.Packet.Decode(buffer); err != nil {
		return fmt.Errorf("decode error, Error: %w", err)
	}

	if err := buffer.Error(); err != nil {
		return err
	}

	return nil
}
