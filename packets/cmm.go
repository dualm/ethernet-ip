package packets

import (
	"github.com/dualm/ethernet-ip/bufferEip"
	"github.com/dualm/ethernet-ip/types"
)

func NewConnectedMessage(connectionID types.UDINT, squenceNumber types.UINT, messageRouterRequest *MessageRouterRequest) (*CommandPacketFormat, error) {
	buffer := bufferEip.New(nil)
	buffer.WriteLittle(connectionID)

	buffer1 := bufferEip.New(nil)
	buffer1.WriteLittle(squenceNumber)

	mr, err := messageRouterRequest.Encode()
	if err != nil {
		return nil, err
	}

	buffer1.WriteLittle(mr)

	return NewCommandPacketFormat([]CommandPacketFormatItem{
		{
			TypeID: ItemIDConnectionBased,
			Data:   buffer.Bytes(),
		},
		{
			TypeID: ItemIDConnectedTransportPacket,
			Data:   buffer1.Bytes(),
		},
	}), nil
}
