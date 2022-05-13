package listidentity

import (
	"github.com/dualm/common"
	"github.com/dualm/ethernet-ip/packets"
	"github.com/dualm/ethernet-ip/packets/command"
	"github.com/dualm/ethernet-ip/types"
)

func New(_context types.ULINT) (*packets.EncapsulationMessagePackets, error) {
	return &packets.EncapsulationMessagePackets{
		Header: packets.EncapsulationHeader{
			Command:       command.ListIdentity,
			Length:        0,
			SessionHandle: 0,
			Status:        0,
			SenderContext: 0,
			Options:       0,
		},
		SpecificData: nil,
	}, nil
}

type ListIdentityItems struct {
	ItemCount int
	Items     []CIPIdentityItem
}

type CIPIdentityItem struct {
	ItemTypeCode    types.UINT
	ItemLength      types.UINT
	ProtocolVersion types.UINT
	SinFamily       types.INT
	SinPort         types.UINT
	SinAddr         types.UDINT
	SinZero         types.ULINT
	VendorID        types.UINT
	DeviceType      types.UINT
	ProductCode     types.UINT
	Major           types.USINT
	Minor           types.USINT
	Status          types.WORD
	SerialNumber    types.UDINT
	NameLength      types.USINT
	ProductName     types.STRING
	State           types.USINT
}

func Decode(packet *packets.EncapsulationMessagePackets) (*ListIdentityItems, error) {
	result := new(ListIdentityItems)
	io := common.NewBuffer(packet.SpecificData)

	io.ReadLittle(&result.ItemCount)

	for i := types.UINT(0); i < types.UINT(result.ItemCount); i++ {
		item := new(CIPIdentityItem)

		io.ReadLittle(&item.ItemTypeCode)
		io.ReadLittle(&item.ItemLength)
		io.ReadLittle(&item.ProtocolVersion)
		io.ReadLittle(&item.SinFamily)
		io.ReadLittle(&item.SinPort)
		io.ReadLittle(&item.SinAddr)
		io.ReadLittle(&item.SinZero)
		io.ReadLittle(&item.VendorID)
		io.ReadLittle(&item.DeviceType)
		io.ReadLittle(&item.ProductCode)
		io.ReadLittle(&item.Major)
		io.ReadLittle(&item.Minor)
		io.ReadLittle(&item.Status)
		io.ReadLittle(&item.SerialNumber)

		io.ReadLittle(&item.NameLength)
		item.ProductName = types.STRING(make([]byte, item.NameLength))
		io.ReadLittle(&item.ProductName)
		io.ReadLittle(&item.State)

		result.Items = append(result.Items, *item)
	}

	return result, nil
}
