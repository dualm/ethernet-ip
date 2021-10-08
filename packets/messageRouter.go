package packets

import (
	"github.com/dualm/ethernet-ip/bufferEip"
	"github.com/dualm/ethernet-ip/types"
	"github.com/dualm/ethernet-ip/utils"
)

type MessageRouterRequest struct {
	Service         types.USINT
	RequestPathSize types.USINT
	RequestPath     []byte
	RequestData     []byte
}

func (m *MessageRouterRequest) Encode() ([]byte, error) {
	if m.RequestPathSize == 0 {
		m.RequestPathSize = utils.Len(m.RequestPath)
	}

	buffer := bufferEip.New(nil)
	buffer.WriteLittle(m.Service)
	buffer.WriteLittle(m.RequestPathSize)
	buffer.WriteLittle(m.RequestPath)
	buffer.WriteLittle(m.RequestData)
	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func NewMessageRouterRequest(service types.USINT, path []byte, data []byte)(*MessageRouterRequest) {
	return &MessageRouterRequest{
		Service:         service,
		RequestPathSize: utils.Len(path),
		RequestPath:     path,
		RequestData:     data,
	}
}

type MessageRouterResponse struct {
	ReplyService           types.USINT
	Reserved               types.USINT
	GeneralStatus          types.USINT
	SizeOfAdditionalStatus types.USINT
	AdditionalStatus       []byte
	ResponseData           []byte
}

func (m *MessageRouterResponse) Decode(raw []byte) error {
	buffer := bufferEip.New(raw)
	buffer.ReadLittle(&m.ReplyService)
	buffer.ReadLittle(&m.Reserved)
	buffer.ReadLittle(&m.GeneralStatus)
	buffer.ReadLittle(&m.SizeOfAdditionalStatus)
	m.AdditionalStatus = make([]byte, m.SizeOfAdditionalStatus*2)
	buffer.ReadLittle(&m.AdditionalStatus)
	m.ResponseData = make([]byte, buffer.Len())
	buffer.ReadLittle(&m.ResponseData)

	return buffer.Error()
}
