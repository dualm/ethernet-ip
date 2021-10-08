package packets

import (
	"github.com/dualm/ethernet-ip/bufferEip"
	"github.com/dualm/ethernet-ip/types"
	"github.com/dualm/ethernet-ip/packets/command"
	"errors"
)

type EncapsulationHeader struct {
	Command       command.Command
	Length        types.UINT
	SessionHandle types.UDINT
	Status        types.UDINT
	SenderContext types.ULINT
	Options       types.UDINT
}

type EncapsulationMessagePackets struct {
	Header       EncapsulationHeader
	SpecificData []byte
}

func (p *EncapsulationMessagePackets) Encode() ([]byte, error) {
	if p.Header.Length > 65511 {
		return nil, errors.New("specific data over length 65511")
	}

	if !command.CheckCommandValid(p.Header.Command) {
		return nil, errors.New("command not supported")
	}

	buf := bufferEip.New(nil)
	buf.WriteLittle(p.Header)
	buf.WriteLittle(p.SpecificData)

	if err := buf.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
