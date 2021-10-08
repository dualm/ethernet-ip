package path

import (
	"github.com/dualm/ethernet-ip/bufferEip"
	"github.com/dualm/ethernet-ip/types"
	"github.com/dualm/ethernet-ip/utils"
)

type SegmentType types.USINT

const (
	PortSegment                SegmentType = 0 << 5
	LogicalSegment             SegmentType = 1 << 5
	NetworkSegmaent            SegmentType = 2 << 5
	SymbolicSegment            SegmentType = 3 << 5
	DataSegment                SegmentType = 4 << 5
	DataTypeConstructedSegment SegmentType = 5 << 5
	DataTypeElementarySegment  SegmentType = 6 << 5
)

type LogicalType types.USINT

const (
	LogicalClassID         LogicalType = 0 << 2
	LogicalInstaceID       LogicalType = 1 << 2
	LogicalMemberID        LogicalType = 2 << 2
	LogicalConnectionPoint LogicalType = 3 << 2
	LogicalAttributeID     LogicalType = 4 << 2
	LogicalSpecial         LogicalType = 5 << 2
	LogicalServiceID       LogicalType = 6 << 2
)

type DataSegmentSubType types.USINT

const (
	SimpleDataSegment DataSegmentSubType = 0x80
	SymbolSegment     DataSegmentSubType = 0x91
)

// format
// 0: 8bit
// 1: 16bit
// 2: 32bit
func LogicalBuild(logicalType LogicalType, value types.UDINT, format uint8, padded bool) ([]byte, error) {
	buffer := bufferEip.New(nil)
	firstByte := uint8(LogicalSegment) | uint8(logicalType) | uint8(format)

	buffer.WriteLittle(firstByte)

	if format == 1 && padded {
		buffer.WriteLittle(uint8(0))
	}

	switch format {
	case 0:
		buffer.WriteLittle(uint8(value))
	case 1:
		buffer.WriteLittle(uint16(value))
	case 2:
		buffer.WriteLittle(uint32(value))
	}

	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func PortBuild(link []byte, portID uint16) ([]byte, error) {
	extentLinkAddressSizebit := len(link) > 1
	extentPortIdentifier := portID > 14

	buffer := bufferEip.New(nil)
	firstByte := uint8(PortSegment)

	if extentPortIdentifier {
		firstByte = firstByte | 0x0f
	} else {
		firstByte = firstByte | uint8(portID)
	}

	if extentLinkAddressSizebit {
		firstByte = firstByte | 0x10
		buffer.WriteLittle(firstByte)
		buffer.WriteLittle(uint8(len(link)))
	} else {
		buffer.WriteLittle(firstByte)
	}

	if extentPortIdentifier {
		buffer.WriteLittle(portID)
	}

	buffer.WriteLittle(link)

	if buffer.Len()%2 == 1 {
		buffer.WriteLittle(uint8(0))
	}

	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Join(args ...[]byte) []byte {
	buffer := bufferEip.New(nil)

	for i := 0; i < len(args); i++ {
		buffer.WriteLittle(args[i])
	}

	return buffer.Bytes()
}

func DataBuild(datatype DataSegmentSubType, raw []byte) ([]byte, error) {
	buffer := bufferEip.New(nil)

	buffer.WriteLittle(datatype)
	
	switch datatype{
	case SimpleDataSegment:
		buffer.WriteLittle(utils.Len(raw))
		buffer.WriteLittle(raw)
	case SymbolSegment:
		l := len(raw)
		buffer.WriteLittle(uint8(l))
		buffer.WriteLittle(raw)

		if l %2 ==1 {
			buffer.WriteLittle(uint8(0))
		}
	}

	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
