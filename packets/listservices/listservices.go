package listservices

import (
	"github.com/dualm/common"
	"github.com/dualm/ethernet-ip/packets"
	"github.com/dualm/ethernet-ip/packets/command"
	"github.com/dualm/ethernet-ip/types"
)

func New(context types.ULINT) (*packets.EncapsulationMessagePackets, error) {
	return &packets.EncapsulationMessagePackets{
		Header: packets.EncapsulationHeader{
			Command:       command.ListServices,
			Length:        0,
			SessionHandle: 0,
			Status:        0,
			SenderContext: context,
			Options:       0,
		},
		SpecificData: nil,
	}, nil
}

type ListServicesItems struct {
	ItemCount int
	Items     []CIPIdentityItem
}

type CIPIdentityItem struct {
	ItemTypeCode    types.UINT
	ItemLength      types.UINT
	Version         types.UINT
	CapabilityFlags types.UINT
	ServicesName    []byte
}

func Decode(packet *packets.EncapsulationMessagePackets) (*ListServicesItems, error) {
	result := new(ListServicesItems)
	buffer := common.NewBuffer(packet.SpecificData)

	buffer.ReadLittle(&result.ItemCount)

	for i := types.UINT(0); i < types.UINT(result.ItemCount); i++ {
		item := CIPIdentityItem{}

		buffer.ReadLittle(&item.ItemTypeCode)
		buffer.ReadLittle(&item.ItemLength)
		buffer.ReadLittle(&item.Version)
		buffer.ReadLittle(&item.CapabilityFlags)
		item.ServicesName = make([]byte, 16)
		buffer.ReadLittle(&item.ServicesName)

		result.Items = append(result.Items, item)
	}

	return result, nil
}
