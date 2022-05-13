package eip

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"unicode"

	"github.com/dualm/common"
	"github.com/dualm/ethernet-ip/packets"
	"github.com/dualm/ethernet-ip/path"
	"github.com/dualm/ethernet-ip/types"
)

const (
	NULL          types.UINT = 0x00
	BOOL          types.UINT = 0xC1
	SINT          types.UINT = 0xC2
	INT           types.UINT = 0xC3
	DINT          types.UINT = 0xC4
	LINT          types.UINT = 0xC5
	USINT         types.UINT = 0xC6
	UINT          types.UINT = 0xC7
	UDINT         types.UINT = 0xC8
	ULINT         types.UINT = 0xC9
	REAL          types.UINT = 0xCA
	LREAD         types.UINT = 0xCB
	STIME         types.UINT = 0xCC
	DATE          types.UINT = 0xCD
	TIME_OF_DAY   types.UINT = 0xCE
	DATE_AND_TIME types.UINT = 0xCF
	STRING        types.UINT = 0xD0
	BYTE          types.UINT = 0xD1
	WORD          types.UINT = 0xD2
	DWORD         types.UINT = 0xD3
	LWORD         types.UINT = 0xD4
	STRING2       types.UINT = 0xD5
	FTIME         types.UINT = 0xD6
	LTIME         types.UINT = 0xD7
	ITIME         types.UINT = 0xD8
	STRINGN       types.UINT = 0xD9
	SHORT_STRING  types.UINT = 0xDA
	TIME          types.UINT = 0xDB
	EPATH         types.UINT = 0xDC
	ENGUINT       types.UINT = 0xDD
	STRINGI       types.UINT = 0xDE
)

var TagTypeMap = map[types.UINT]string{
	NULL:  "NULL",
	BOOL:  "BOOL",
	SINT:  "SINT",
	INT:   "INT",
	DINT:  "DINT",
	REAL:  "REAL",
	DWORD: "DWORD",
	LINT:  "LINT",
}

type Tag struct {
	Lock *sync.Mutex
	EIP  *EIPConn

	instanceID types.UDINT
	nameLen    types.UINT
	name       []byte
	Type       types.UINT // data access page.44, 15bit: 1=struct, 0=atomic
	dim1Len    types.UDINT
	dim2Len    types.UDINT
	dim3Len    types.UDINT
	changed    bool

	value    []byte
	mValue   []byte
	OnChange func()

	readRequestMsg *packets.MessageRouterRequest
}

func (tag *Tag) SetDriver(driver interface{}) {
	tag.EIP = driver.(*EIPConn)
}

func (tag *Tag) GetValue() []byte {
	tag.Lock.Lock()
	defer tag.Lock.Unlock()

	return tag.value
}

func (tag *Tag) Read() error {
	tag.Lock.Lock()
	defer tag.Lock.Unlock()

	if tag.readRequestMsg == nil {
		readRequest, err := tag.readRequest()
		if err != nil {
			return err
		}

		tag.readRequestMsg = readRequest
	}

	res, err := tag.EIP.Send(tag.readRequestMsg)
	if err != nil {
		return err
	}

	mrres := new(packets.MessageRouterResponse)

	if err := mrres.Decode(res.Packet.Items[1].Data); err != nil {
		return fmt.Errorf("decode error, Error: %w", err)
	}

	if err := tag.readParser(mrres, nil); err != nil {
		return fmt.Errorf("readParser error, Error: %w", err)
	}

	return nil
}

func (tag *Tag) readRequest() (*packets.MessageRouterRequest, error) {
	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle(tag.count())
	if err := buffer.Error(); err != nil {
		return nil, err
	}

	// Symbolic Segment Addressing
	path, err := path.DataBuild(path.SymbolSegment, tag.name)
	if err != nil {
		return nil, err
	}

	// Symbolic Segment Addressing

	// Symbol Instance Addressing
	// symbol class id 0x6B
	// classID, err := path.LogicalBuild(path.LogicalClassID, 0x6B, 0, true)
	// if err != nil {
	// 	return nil, err
	// }

	// instanceID, err := path.LogicalBuild(path.LogicalInstaceID, types.UDINT(tag.instanceID), 0, true)
	// if err != nil {
	// 	return nil, err
	// }
	// path := path.Join(classID, instanceID)
	// Symbol Instance Addressing

	messageRouterRequest := packets.NewMessageRouterRequest(packets.ServiceReadTag, path, buffer.Bytes())

	return messageRouterRequest, nil
}

func (tag *Tag) readParser(response *packets.MessageRouterResponse, cb func(func())) error {
	buffer := common.NewBuffer(response.ResponseData)

	_t := uint16(0)
	buffer.ReadLittle(&_t)

	// 0x2a0
	// Tag Type Service Parameter for structures
	if _t == 0x02a0 {
		buffer.ReadLittle(&_t)
	}

	payload := make([]byte, buffer.Len())
	buffer.ReadLittle(payload)

	if err := buffer.Error(); err != nil {
		return err
	}

	if !bytes.Equal(tag.value, payload) {
		tag.value = payload
		if tag.OnChange != nil {
			if cb == nil {
				go tag.OnChange()
			} else {
				go cb(tag.OnChange)
			}
		}
	}

	return nil
}

func (tag *Tag) Write() error {
	tag.Lock.Lock()
	defer tag.Lock.Unlock()

	if tag.mValue == nil {
		return nil
	}

	writeRequest, err := tag.writeRequest()
	if err != nil {
		return err
	}

	multiWriteRequest, err := multiple(writeRequest)
	if err != nil {
		return err
	}

	_, err = tag.EIP.Send(multiWriteRequest)
	if err != nil {
		return err
	}

	if tag.mValue != nil {
		copy(tag.value, tag.mValue)

		tag.mValue = nil
	}

	return nil
}

func (tag *Tag) writeRequest() ([]*packets.MessageRouterRequest, error) {
	var result []*packets.MessageRouterRequest

	// atomic
	if 0x8000&tag.Type == 0 {
		buffer := common.NewEmptyBuffer()

		buffer.WriteLittle(tag.Type)
		buffer.WriteLittle(tag.count())
		buffer.WriteLittle(tag.mValue)

		// symbolic segment addressing
		path, err := path.DataBuild(path.SymbolSegment, tag.name)
		if err != nil {
			return nil, err
		}

		// symbolic instance addressing
		// classID, err := path.LogicalBuild(path.LogicalClassID, 0x6B, 0, true)
		// if err != nil {
		// 	return nil, err
		// }

		// instanceID, err := path.LogicalBuild(path.LogicalInstaceID, types.UDINT(tag.instanceID), 0, true)
		// if err != nil {
		// 	return nil, err
		// }
		// path := path.Join(classID, instanceID)

		messageRouterRequest := packets.NewMessageRouterRequest(
			packets.ServiceWriteTag,
			path,
			buffer.Bytes(),
		)

		result = append(result, messageRouterRequest)
	} else {
		buffer := common.NewEmptyBuffer()

		buffer.WriteLittle(DINT)
		buffer.WriteLittle(types.UINT(1))
		buffer.WriteLittle(types.UDINT(len(tag.mValue)))

		classID, err := path.LogicalBuild(path.LogicalClassID, 0x6B, 0, true)
		if err != nil {
			return nil, err
		}

		instanceID, err := path.LogicalBuild(path.LogicalInstaceID, types.UDINT(tag.instanceID), 0, true)
		if err != nil {
			return nil, err
		}

		data, err := path.DataBuild(path.SymbolSegment, []byte("LEN"))
		if err != nil {
			return nil, err
		}

		messageRouterRequest1 := packets.NewMessageRouterRequest(
			packets.ServiceWriteTag,
			path.Join(classID, instanceID, data),
			buffer.Bytes(),
		)

		result = append(result, messageRouterRequest1)

		buffer1 := common.NewEmptyBuffer()

		buffer1.WriteLittle(SINT)
		buffer1.WriteLittle(types.UINT(len(tag.mValue)))
		buffer1.WriteLittle(tag.mValue)

		data, err = path.DataBuild(
			path.SymbolSegment, []byte("DATA"))
		if err != nil {
			return nil, err
		}

		messageRouterRequest2 := packets.NewMessageRouterRequest(
			packets.ServiceWriteTag,
			path.Join(classID, instanceID, data),
			buffer.Bytes())

		result = append(result, messageRouterRequest2)
	}

	return result, nil
}

func (tag *Tag) SetValue(data []byte) {
	// tag.Lock.Lock()
	// defer tag.Lock.Unlock()
	tag.changed = true

	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle(data)

	tag.mValue = buffer.Bytes()
}

func (tag *Tag) SetInt32(i int32) {
	tag.changed = true

	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle(i)
	tag.mValue = buffer.Bytes()
}

func (tag *Tag) SetString(i string) {
	tag.changed = true

	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle([]byte(i))

	tag.mValue = buffer.Bytes()
}

func (tag *Tag) SetType(word types.UINT) {
	tag.Type = word
}

func (tag *Tag) dims() types.USINT {
	return types.USINT((0x6000 & tag.Type) >> 13)
}

func (tag *Tag) TypeString() string {
	var _type string
	if 0x8000&tag.Type == 0 {
		_type = "atomic"
	} else {
		_type = "struct"
	}

	return fmt.Sprintf("%#04x(%6s) | %s | %d dims", uint16(tag.Type), TagTypeMap[0xFFF&tag.Type], _type, (0x6000&tag.Type)>>13)
}

func (tag *Tag) Name() string {
	return string(tag.name)
}

func (tag *Tag) count() types.UINT {
	a := types.UINT(1)

	if tag.dim1Len > 0 {
		a = types.UINT(tag.dim1Len)
	}

	b := types.UINT(1)
	if tag.dim2Len > 0 {
		b = types.UINT(tag.dim2Len)
	}

	c := types.UINT(1)
	if tag.dim3Len > 0 {
		c = types.UINT(tag.dim3Len)
	}

	return a * b * c
}

func (tag *Tag) Int32() (int32, error) {
	buffer := common.NewBuffer(tag.value)

	var val int32

	buffer.ReadLittle(&val)
	if err := buffer.Error(); err != nil {
		return 0, err
	}

	return val, nil
}

func (tag *Tag) String() (string, error) {
	buffer := common.NewBuffer(tag.value)

	l := types.UDINT(0)

	buffer.ReadLittle(&l)
	if l > 88 {
		return "", nil
	}

	val := make([]byte, l)
	buffer.ReadLittle(val)
	for i := range val {
		if !unicode.IsPrint(rune(val[i])) {
			return "", errors.New("some rune can't print")
		}
	}

	return string(val), nil
}

func (tag *Tag) XInt32() (int32, error) {
	var val []byte
	if len(tag.mValue) > 0 {
		val = tag.mValue
	} else {
		val = tag.value
	}

	buffer := common.NewBuffer(val)

	var v int32
	buffer.ReadLittle(&v)

	if err := buffer.Error(); err != nil {
		return 0, err
	}
	return v, nil
}

func (tag *Tag) XString() (string, error) {
	var value []byte
	if len(tag.mValue) > 0 {
		value = tag.mValue
	} else {
		value = tag.value
	}

	buffer := common.NewBuffer(value)

	l := types.UDINT(0)
	buffer.ReadLittle(&l)
	if l > 88 {
		return "", nil
	}

	val := make([]byte, l)
	buffer.ReadLittle(val)
	for i := range val {
		if !unicode.IsPrint(rune(val[i])) {
			return "", errors.New("some rune can't print")
		}
	}

	if err := buffer.Error(); err != nil {
		return "", err
	}

	return string(value), nil
}

func multiple(messageRouterRequests []*packets.MessageRouterRequest) (*packets.MessageRouterRequest, error) {
	l := len(messageRouterRequests)
	if l == 1 {
		return messageRouterRequests[0], nil
	}

	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle(types.UINT(l))

	offset := 2 * (l + 1)
	buffer.WriteLittle(types.UINT(offset))
	for i := range messageRouterRequests {
		if i != l-1 {
			data, err := messageRouterRequests[i].Encode()
			if err != nil {
				return nil, err
			}

			offset += len(data)
			buffer.WriteLittle(offset)
		}
	}

	for i := range messageRouterRequests {
		data, err := messageRouterRequests[i].Encode()
		if err != nil {
			return nil, err
		}

		buffer.WriteLittle(data)
	}

	if err := buffer.Error(); err != nil {
		return nil, err
	}

	classID, err := path.LogicalBuild(path.LogicalClassID, 0x02, 0, true)
	if err != nil {
		return nil, err
	}

	instanceID, err := path.LogicalBuild(path.LogicalInstaceID, 0x01, 0, true)
	if err != nil {
		return nil, err
	}

	return packets.NewMessageRouterRequest(
		packets.ServiceMultipleServicePacket,
		path.Join(
			classID, instanceID,
		),
		buffer.Bytes()), nil
}

func (eip *EIPConn) AllTags() (map[string]*Tag, error) {
	result := make(map[string]*Tag)

	return eip.allTags(result, 0)
}

func (eip *EIPConn) allTags(tagMap map[string]*Tag, instanceID types.UDINT) (map[string]*Tag, error) {
	classPath, err := path.LogicalBuild(path.LogicalClassID, 0x6B, 0, true)
	if err != nil {
		return nil, err
	}

	instancePath, err := path.LogicalBuild(path.LogicalInstaceID, instanceID, 0, true)
	if err != nil {
		return nil, err
	}

	paths := path.Join(
		classPath,
		instancePath,
	)

	buffer := common.NewEmptyBuffer()

	buffer.WriteLittle(types.UINT(3))
	buffer.WriteLittle(types.UINT(1))
	buffer.WriteLittle(types.UINT(2))
	buffer.WriteLittle(types.UINT(8))
	if err := buffer.Error(); err != nil {
		return nil, err
	}

	messageRouterRequest := packets.NewMessageRouterRequest(
		packets.ServiceGetInstanceAttributeList, paths, buffer.Bytes())

	res, err := eip.Send(messageRouterRequest)
	if err != nil {
		return nil, err
	}

	mrres := new(packets.MessageRouterResponse)
	if err := mrres.Decode(res.Packet.Items[1].Data); err != nil {
		return nil, fmt.Errorf("decode error, Error: %w", err)
	}

	buffer1 := common.NewBuffer(mrres.ResponseData)

	for buffer1.Len() > 0 {
		tag := new(Tag)
		tag.EIP = eip
		tag.Lock = new(sync.Mutex)

		buffer1.ReadLittle(&tag.instanceID)
		buffer1.ReadLittle(&tag.nameLen)
		tag.name = make([]byte, tag.nameLen)
		buffer1.ReadLittle(&tag.name)
		buffer1.ReadLittle(&tag.Type)
		buffer1.ReadLittle(&tag.dim1Len)
		buffer1.ReadLittle(&tag.dim2Len)
		buffer1.ReadLittle(&tag.dim3Len)

		tagMap[tag.Name()] = tag
		instanceID = tag.instanceID
	}

	if mrres.GeneralStatus == 0x60 {
		return eip.allTags(tagMap, instanceID)
	}

	return tagMap, nil
}

type TagGroup struct {
	tags map[types.UDINT]*Tag
	EIP  *EIPConn
	Lock *sync.Mutex
}

func NewTagGroup(lock *sync.Mutex) *TagGroup {
	return &TagGroup{
		tags: make(map[types.UDINT]*Tag),
		Lock: lock,
	}
}

func (tg *TagGroup) Add(tag *Tag) {
	if tg.EIP == nil {
		tg.EIP = tag.EIP
	} else {
		if tg.EIP != tag.EIP {
			return
		}
	}

	tg.tags[tag.instanceID] = tag
}

func (tg *TagGroup) Remove(tag *Tag) {
	delete(tg.tags, tag.instanceID)
}

func (tg *TagGroup) Read() error {
	tg.Lock.Lock()
	defer tg.Lock.Unlock()

	if len(tg.tags) == 0 {
		return nil
	}

	if len(tg.tags) == 1 {
		for _, v := range tg.tags {
			return v.Read()
		}
	}

	var list []types.UDINT
	var mrs []*packets.MessageRouterRequest

	for i := range tg.tags {
		one := tg.tags[i]

		one.Lock.Lock()
		list = append(list, one.instanceID)
		request, err := one.readRequest()
		if err != nil {
			return err
		}

		mrs = append(mrs, request)

		one.Lock.Unlock()
	}

	_sb, err := multiple(mrs)
	if err != nil {
		return err
	}

	res, err := tg.EIP.Send(_sb)
	if err != nil {
		return err
	}

	rmr := new(packets.MessageRouterResponse)
	if err := rmr.Decode(res.Packet.Items[1].Data); err != nil {
		return fmt.Errorf("decode error, Error: %w", err)
	}

	buffer1 := common.NewBuffer(rmr.ResponseData)

	count := types.UINT(0)
	buffer1.ReadLittle(&count)

	if int(count) != len(list) {
		return nil
	}

	var offset []types.UINT
	for i := types.UINT(0); i < count; i++ {
		one := types.UINT(0)
		buffer1.ReadLittle(&one)
		offset = append(offset, one)
	}

	var cbs []func()
	for i2 := range list {
		mr := new(packets.MessageRouterResponse)

		if (i2 + 1) != len(offset) {
			if err := mr.Decode(rmr.ResponseData[offset[i2]:offset[i2+1]]); err != nil {
				return err
			}
		} else {
			if err := mr.Decode(rmr.ResponseData[offset[i2]:]); err != nil {
				return err
			}
		}

		if err := tg.tags[list[i2]].readParser(mr, func(f func()) {
			cbs = append(cbs, f)
		}); err != nil {
			return err
		}
	}

	for i := range cbs {
		go cbs[i]()
	}

	return nil
}

func (tg *TagGroup) Write() error {
	tg.Lock.Lock()
	defer tg.Lock.Unlock()

	var list []types.UDINT
	var mrs []*packets.MessageRouterRequest

	for i := range tg.tags {
		one := tg.tags[i]

		one.Lock.Lock()
		if one.changed {
			list = append(list, one.instanceID)

			writeRequest, err := one.writeRequest()
			if err != nil {
				return err
			}

			mrs = append(mrs, writeRequest...)
			one.changed = false
		}
		one.Lock.Unlock()
	}

	if len(list) == 0 {
		return nil
	}

	multiple, err := multiple(mrs)
	if err != nil {
		return err
	}

	_, err = tg.EIP.Send(multiple)
	if err != nil {
		return err
	}

	for i := range tg.tags {
		copy(tg.tags[i].value, tg.tags[i].mValue)
		tg.tags[i].mValue = nil
	}

	return nil
}

func NewTag(eip *EIPConn, name string, count int, onChange func()) *Tag {
	return &Tag{
		Lock:       &sync.Mutex{},
		EIP:        eip,
		instanceID: 0,
		nameLen:    types.UINT(len(name)),
		name:       []byte(name),
		Type:       0x00C3,
		dim1Len:    types.UDINT(count),
		dim2Len:    0,
		dim3Len:    0,
		changed:    false,
		value:      []byte{},
		mValue:     []byte{},
		OnChange:   onChange,
	}
}
