package packets

import (
	"github.com/dualm/common"
	"github.com/dualm/ethernet-ip/types"
)

func NewConnectedMessage(connectionID types.UDINT, sequenceNumber types.UINT, messageRouterRequest *MessageRouterRequest) (*CommandPacketFormat, error) {
	buffer := common.NewEmptyBuffer()
	buffer.WriteLittle(connectionID)

	if err := buffer.Error(); err != nil {
		return nil, err
	}

	buffer1 := common.NewEmptyBuffer()
	buffer1.WriteLittle(sequenceNumber)

	mr, err := messageRouterRequest.Encode()
	if err != nil {
		return nil, err
	}

	buffer1.WriteLittle(mr)

	if err := buffer1.Error(); err != nil {
		return nil, err
	}

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
