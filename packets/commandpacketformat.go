package packets

import (
	"gitee.com/ziIoT/common"
	"gitee.com/ziIoT/ethernet-ip/types"
)

type ItemID types.UINT

const (
	ItemIDUCMM                     ItemID = 0x0000
	ItemIDListIdentityResponse     ItemID = 0x000C
	ItemIDConnectionBased          ItemID = 0x00A1
	ItemIDConnectedTransportPacket ItemID = 0x00B1
	ItemIDUnconnectedMessage       ItemID = 0x00B2
	ItemIDListServicesResponse     ItemID = 0x0100
	ItemIDSockaddrInfoOToT         ItemID = 0x8000
	ItemIDSockaddrInfoTToO         ItemID = 0x8001
	ItemIDSequencedAddressItem     ItemID = 0x8002
)

type CommandPacketFormat struct {
	ItemCount types.UINT
	Items     []CommandPacketFormatItem
}

func NewCommandPacketFormat(items []CommandPacketFormatItem) *CommandPacketFormat {
	return &CommandPacketFormat{
		ItemCount: types.UINT(len(items)),
		Items:     items,
	}
}

func (cpf *CommandPacketFormat) Encode() ([]byte, error) {
	if cpf.ItemCount == 0 {
		cpf.ItemCount = types.UINT(len(cpf.Items))
	}

	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle(cpf.ItemCount)

	for _, item := range cpf.Items {
		itemBytes, err := item.Encode()
		if err != nil {
			return nil, err
		}

		buffer.WriteLittle(itemBytes)
	}

	return buffer.Bytes(), nil
}

func (cpf *CommandPacketFormat) Decode(buffer *common.Buffer) error {
	buffer.ReadLittle(&cpf.ItemCount)
	if err := buffer.Error(); err != nil {
		return err
	}

	cpf.Items = make([]CommandPacketFormatItem, cpf.ItemCount)
	for i := types.UINT(0); i < cpf.ItemCount; i++ {
		item := CommandPacketFormatItem{}
		item.Decode(buffer)
		if err := buffer.Error(); err != nil {
			return err
		}

		cpf.Items[i] = item
	}

	return nil
}

type CommandPacketFormatItem struct {
	TypeID ItemID
	Length types.UINT
	Data   []byte
}

func (item *CommandPacketFormatItem) Encode() ([]byte, error) {
	if item.Length == 0 {
		item.Length = types.UINT(len(item.Data))
	}

	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle(item.TypeID)
	buffer.WriteLittle(item.Length)
	buffer.WriteLittle(item.Data)

	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (item *CommandPacketFormatItem) Decode(raw *common.Buffer) {
	raw.ReadLittle(&item.TypeID)
	raw.ReadLittle(&item.Length)
	item.Data = make([]byte, item.Length)
	raw.ReadLittle(&item.Data)
}
