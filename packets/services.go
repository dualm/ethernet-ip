package packets

import (
	"github.com/dualm/ethernet-ip/types"
)

const (
	ServiceGetAttributesAll          types.USINT = 0x01
	ServiceGetAttributeSignle        types.USINT = 0x0E
	ServiceSetAttributeSingle        types.USINT = 0x10
	ServiceForwardOpen               types.USINT = 0x4E
	ServiceUnconnectedSend           types.USINT = 0x52
	ServiceForwardClose              types.USINT = 0x54
	ServiceGetConnectionData         types.USINT = 0x56
	ServiceSearchConnectionData      types.USINT = 0x57
	ServiceGetConnectionOwner        types.USINT = 0x5A
	ServiceLargeForwardOpen          types.USINT = 0x5B
	ServiceReadTag                   types.USINT = 0x4C
	ServiceReadTagFragmented         types.USINT = 0x52
	ServiceWriteTag                  types.USINT = 0x4D
	ServiceWriteTagFragmentedService types.USINT = 0x53
	ServiceReadModifyWriteTagService types.USINT = 0x4E
	ServiceMultipleServicePacket     types.USINT = 0x0a
	ServiceGetInstanceAttributeList  types.USINT = 0x55
)
