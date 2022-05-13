package packets

import (
	"github.com/dualm/common"
	"github.com/dualm/ethernet-ip/path"
	"github.com/dualm/ethernet-ip/types"
	"github.com/dualm/ethernet-ip/utils"
)

type UnConnectedSendServiceParameters struct {
	PriorityTimeTick   types.USINT
	TimeoutTicks       types.USINT
	MessageRequestSize types.UINT
	MessageRequest     *MessageRouterRequest
	Pad                types.USINT
	RoutePathSize      types.USINT
	Reserved           types.USINT
	RoutePath          []byte
}

func (u *UnConnectedSendServiceParameters) Encode() ([]byte, error) {
	messageRouterRequest, err := u.MessageRequest.Encode()
	if err != nil {
		return nil, err
	}

	buffer := common.NewEmptyBuffer()
	buffer.WriteLittle(u.PriorityTimeTick)
	buffer.WriteLittle(u.TimeoutTicks)

	l := len(messageRouterRequest)

	buffer.WriteLittle(types.UINT(l))
	buffer.WriteLittle(messageRouterRequest)

	if l%2 == 1 {
		buffer.WriteLittle(types.USINT(0))
	}

	buffer.WriteLittle(utils.Len(u.RoutePath))
	buffer.WriteLittle(types.USINT(0))
	buffer.WriteLittle(u.RoutePath)

	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func UnConnectedMessageRouterRequest(slot uint8, timeTick types.USINT, timeoutTicks types.USINT, mr *MessageRouterRequest) (*MessageRouterRequest, error) {
	port, err := path.PortBuild([]byte{slot}, 1)
	if err != nil {
		return nil, err
	}

	ucmr := UnConnectedSendServiceParameters{
		PriorityTimeTick: timeTick,
		TimeoutTicks:     timeoutTicks,
		MessageRequest:   mr,
		RoutePath:        port,
	}

	data, err := ucmr.Encode()
	if err != nil {
		return nil, err
	}

	classID, err := path.LogicalBuild(path.LogicalClassID, 06, 0, true)
	if err != nil {
		return nil, err
	}

	instanceID, err := path.LogicalBuild(path.LogicalInstaceID, 01, 0, true)
	if err != nil {
		return nil, err
	}

	return NewMessageRouterRequest(ServiceUnconnectedSend, path.Join(classID, instanceID), data), nil
}

func NewUnconnectedMessage(messageRR *MessageRouterRequest) (*CommandPacketFormat, error) {
	mr, err := messageRR.Encode()
	if err != nil {
		return nil, err
	}

	cpf := NewCommandPacketFormat(
		[]CommandPacketFormatItem{
			{
				TypeID: ItemIDUCMM,
				Data:   nil,
			},
			{
				TypeID: ItemIDUnconnectedMessage,
				Data:   mr,
			},
		},
	)

	return cpf, nil
}
