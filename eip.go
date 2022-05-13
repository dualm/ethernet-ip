package eip

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/dualm/common"
	"github.com/dualm/ethernet-ip/packets"
	"github.com/dualm/ethernet-ip/packets/listidentity"
	"github.com/dualm/ethernet-ip/packets/listinterfaces"
	"github.com/dualm/ethernet-ip/packets/listservices"
	"github.com/dualm/ethernet-ip/packets/registersession"
	"github.com/dualm/ethernet-ip/packets/sendrrdata"
	"github.com/dualm/ethernet-ip/packets/sendunitdata"
	"github.com/dualm/ethernet-ip/packets/unregistersession"
	"github.com/dualm/ethernet-ip/types"
	"github.com/dualm/ethernet-ip/utils"
)

type EIPConn struct {
	config  *Config
	tcpAddr *net.TCPAddr
	tcpConn *net.TCPConn
	udpAddr *net.UDPAddr
	udpConn *net.UDPConn
	session types.UDINT

	established  bool
	connectionID types.UDINT
	seqNum       types.UINT

	requestLock *sync.Mutex
}

func (eip *EIPConn) Connect() error {
	tcpConn, err := net.DialTCP("tcp", nil, eip.tcpAddr)
	if err != nil {
		return err
	}

	err = tcpConn.SetKeepAlive(true)
	if err != nil {
		return err
	}

	eip.tcpConn = tcpConn

	if err := eip.RegisterSession(); err != nil {
		return err
	}

	return nil
}

func (eip *EIPConn) Close() error {
	if eip.tcpConn != nil {
		return eip.tcpConn.Close()
	}

	return nil
}

// todo
func (eip *EIPConn) ForwardOpen() error {
	// buffer := bufferEip.New(nil)

	// Tick_tick
	return nil
}

// todo
func (eip *EIPConn) ForwardClose() error {
	return nil
}

func (eip *EIPConn) read() (*packets.EncapsulationMessagePackets, error) {
	buf := make([]byte, 1024*64)

	length, err := eip.tcpConn.Read(buf)
	if err != nil {
		return nil, err
	}

	return eip.parse(buf[:length])
}

func (eip *EIPConn) write(data []byte) error {
	_, err := eip.tcpConn.Write(data)

	return err
}

func (eip *EIPConn) parse(buf []byte) (*packets.EncapsulationMessagePackets, error) {
	if len(buf) < 24 {
		return nil, errors.New("invalid packet, length < 24")
	}

	_packet := new(packets.EncapsulationMessagePackets)

	buffer := common.NewBuffer(buf)

	buffer.ReadLittle(&_packet.Header)
	if err := buffer.Error(); err != nil {
		return nil, err
	}

	if _packet.Header.Options != 0 {
		return nil, errors.New("wrong packet with non-zero options")
	}

	if int(_packet.Header.Length) != buffer.Len() {
		return nil, errors.New("wrong packet length")
	}

	_packet.SpecificData = make([]byte, _packet.Header.Length)
	buffer.ReadLittle(_packet.SpecificData)
	if err := buffer.Error(); err != nil {
		return nil, err
	}

	return _packet, nil
}

func (eip *EIPConn) request(packet *packets.EncapsulationMessagePackets) (*packets.EncapsulationMessagePackets, error) {
	eip.requestLock.Lock()
	defer eip.requestLock.Unlock()

	if eip.tcpConn == nil {
		return nil, errors.New("invalid tcp connection, connect first")
	}

	b, err := packet.Encode()
	if err != nil {
		return nil, err
	}

	if err := eip.write(b); err != nil {
		return nil, err
	}

	return eip.read()
}

func (eip *EIPConn) RegisterSession() error {
	ctx := utils.GetNewContext()

	request, err := registersession.New(ctx)
	if err != nil {
		return err
	}

	response, err := eip.request(request)
	if err != nil {
		return err
	}

	eip.session = response.Header.SessionHandle

	return nil
}

func (eip *EIPConn) UnRegisterSession() error {
	ctx := utils.GetNewContext()

	request, err := unregistersession.New(eip.session, ctx)
	if err != nil {
		return err
	}

	response, err := eip.request(request)
	if err != nil {
		return err
	}

	eip.session = response.Header.SessionHandle

	return nil
}

func NewEIP(address string, config *Config) (*EIPConn, error) {
	if config == nil {
		config = DefaultConfig()
	}

	tcpAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, config.TCPPort))
	if err != nil {
		return nil, err
	}

	udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, config.UDPPort))
	if err != nil {
		return nil, err
	}

	return &EIPConn{
		config:       config,
		tcpAddr:      tcpAddress,
		udpAddr:      udpAddress,
		session:      0,
		established:  false,
		connectionID: 0,
		seqNum:       0,
		requestLock:  new(sync.Mutex),
	}, nil
}

func (eip *EIPConn) ListInterface() (*listinterfaces.ListInterfaceItems, error) {
	ctx := utils.GetNewContext()

	request, err := listinterfaces.New(ctx)
	if err != nil {
		return nil, err
	}

	response, err := eip.request(request)
	if err != nil {
		return nil, err
	}

	return listinterfaces.Decode(response)
}

func (eip *EIPConn) ListServices() (*listservices.ListServicesItems, error) {
	ctx := utils.GetNewContext()

	request, err := listservices.New(ctx)
	if err != nil {
		return nil, err
	}

	response, err := eip.request(request)
	if err != nil {
		return nil, err
	}

	return listservices.Decode(response)
}

func (eip *EIPConn) ListIdentity() (*listidentity.ListIdentityItems, error) {
	ctx := utils.GetNewContext()

	request, err := listidentity.New(ctx)
	if err != nil {
		return nil, err
	}

	response, err := eip.request(request)
	if err != nil {
		return nil, err
	}

	return listidentity.Decode(response)
}

func (eip *EIPConn) SendRRData(cpf *packets.CommandPacketFormat, timeout types.UINT) (*packets.SpecificData, error) {
	ctx := utils.GetNewContext()

	request, err := sendrrdata.New(eip.session, ctx, cpf, timeout)
	if err != nil {
		return nil, err
	}

	response, err := eip.request(request)
	if err != nil {
		return nil, err
	}

	return sendrrdata.Decode(response)
}

func (eip *EIPConn) SendUnitData(cpf *packets.CommandPacketFormat) (*packets.SpecificData, error) {
	ctx := utils.GetNewContext()

	request, err := sendunitdata.New(eip.session, ctx, cpf)
	if err != nil {
		return nil, err
	}

	response, err := eip.request(request)
	if err != nil {
		return nil, err
	}

	return sendunitdata.Decode(response)
}

func (eip *EIPConn) Send(messageRouterRequest *packets.MessageRouterRequest) (*packets.SpecificData, error) {
	if !eip.established {
		mr, err := packets.UnConnectedMessageRouterRequest(
			eip.config.Slot,
			eip.config.TimeTick,
			eip.config.TimeTickOut,
			messageRouterRequest,
		)
		if err != nil {
			return nil, err
		}

		messageRouterRequest = mr
	}

	if eip.established {
		eip.seqNum += 1

		message, err := packets.NewConnectedMessage(eip.connectionID, eip.seqNum, messageRouterRequest)
		if err != nil {
			return nil, err
		}

		return eip.SendUnitData(message)
	} else {
		message, err := packets.NewUnconnectedMessage(messageRouterRequest)
		if err != nil {
			return nil, err
		}

		return eip.SendRRData(message, types.UINT(eip.config.TimeTickOut))
	}
}
