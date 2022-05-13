package listinterfaces

import (
	"github.com/dualm/common"
	"github.com/dualm/ethernet-ip/packets"
	"github.com/dualm/ethernet-ip/packets/command"
	"github.com/dualm/ethernet-ip/types"
)

func New(_context types.ULINT) (*packets.EncapsulationMessagePackets, error) {
	return &packets.EncapsulationMessagePackets{
		Header: packets.EncapsulationHeader{
			Command:       command.ListInterfaces,
			Length:        0,
			SessionHandle: 0,
			Status:        0,
			SenderContext: _context,
			Options:       0,
		},
		SpecificData: nil,
	}, nil
}

type ListInterfaceItems struct {
	ItemCount int
	Items     []CIPIdentityItem
}

type CIPIdentityItem struct {
	ItemTypeCode types.UINT
	ItemLength   types.UINT
	ItemData     []byte
}

func Decode(packet *packets.EncapsulationMessagePackets) (*ListInterfaceItems, error) {
	result := new(ListInterfaceItems)
	buffer := common.NewBuffer(packet.SpecificData)
	buffer.ReadLittle(&result.ItemCount)

	for i := types.UINT(0); i < types.UINT(result.ItemCount); i++ {
		item := CIPIdentityItem{}

		buffer.ReadLittle(&item.ItemTypeCode)
		buffer.ReadLittle(&item.ItemLength)
		item.ItemData = make([]byte, item.ItemLength)
		buffer.ReadLittle(&item.ItemData)

		result.Items = append(result.Items, item)
	}

	return result, nil
}
