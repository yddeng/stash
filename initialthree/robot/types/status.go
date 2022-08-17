package types

import (
	"encoding/xml"
	"fmt"
)

const statusBits = 32

type Status int32

type StatusID int8

const (
	Status_IsLogin = StatusID(iota)
	Status_IsFirstlogin
)

var statusDefines = [...]statusDefine{
	Status_IsLogin:      newBooleanStatus("isLogin", 0),
	Status_IsFirstlogin: newBooleanStatus("isFirstLogin", 1),
}

func getStatusDefine(s StatusID) statusDefine {
	if s < 0 || int(s) >= len(statusDefines) {
		panic(fmt.Errorf("status %d no define", s))
	}
	return statusDefines[s]
}

var statusName2ID = map[string]StatusID{}

type statusDefine struct {
	name     string
	offset   int
	bit      int
	values   []StatusValue
	valueSet map[StatusValue]struct{}
}

func newBooleanStatus(name string, offset int) statusDefine {
	if offset < 0 || offset > statusBits-1 {
		panic("offset out of range")
	}

	return statusDefine{
		name:   name,
		offset: offset,
		bit:    1,
	}
}

func newStatus(name string, offset, bit int, value ...StatusValue) statusDefine {
	if offset < 0 || offset > statusBits-1 {
		panic("offset out of range")
	}

	if bit <= 0 {
		panic("invalid bit")
	}

	if offset+bit > statusBits {
		panic("overflow")
	}

	d := statusDefine{
		name:   name,
		offset: offset,
		bit:    bit,
	}

	if len(value) == 0 {
		panic("no value")
	}

	d.valueSet = make(map[StatusValue]struct{})
	for i, v := range value {
		if v > StatusValue(d.max()) {
			panic(fmt.Errorf("the No.%d value %d is overflow", i, v))
		}
		d.valueSet[v] = struct{}{}
	}
	d.values = value

	return d
}

func (s statusDefine) isBoolean() bool {
	return s.bit == 1
}

func (s statusDefine) mask() int {
	return (1<<s.bit - 1) << s.offset
}

func (s statusDefine) max() int {
	return (1<<s.bit - 1)
}

func (s statusDefine) IsOverlappedWith(o statusDefine) bool {
	l1, h1 := s.offset, s.offset+s.bit-1
	l2, h2 := o.offset, o.offset+o.bit-1
	return h1 >= l2 && l1 <= h2
}

func (s StatusID) String() string {
	return statusDefines[s].name
}

func (s StatusID) IsBoolean() bool {
	return statusDefines[s].isBoolean()
}

func (s StatusID) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(statusDefines[s].name, start)
}

func (s *StatusID) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var name string

	if err := d.DecodeElement(&name, &start); err != nil {
		return err
	}

	*s = statusName2ID[name]
	return nil

}

func (s StatusID) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: statusDefines[s].name}, nil
}

func (s *StatusID) UnmarshalXMLAttr(attr xml.Attr) error {
	*s = statusName2ID[attr.Value]
	return nil
}

type StatusValue int8

func (v StatusValue) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: ""}, nil
}

func (v *StatusValue) UnmarshalXMLAttr(attr xml.Attr) error {
	// todo
	return nil
}

func (ss Status) Is(s StatusID) bool {
	def := getStatusDefine(s)
	if !def.isBoolean() {
		panic(fmt.Errorf("status %s not boolean", s))
	}
	return ss&Status(def.mask()) != 0
}

func (ss *Status) set(s StatusID, b bool) {
	def := getStatusDefine(s)
	if !def.isBoolean() {
		panic(fmt.Errorf("status %s not boolean", s))
	}
	if b {
		*ss |= Status(def.mask())
	} else {
		*ss &= ^Status(def.mask())
	}
}

func (ss *Status) Set(s StatusID) {
	ss.set(s, true)
}

func (ss *Status) Unset(s StatusID) {
	ss.set(s, false)
}

func (ss Status) GetStatusValue(s StatusID) StatusValue {
	def := getStatusDefine(s)
	if def.isBoolean() {
		panic(fmt.Errorf("status %s is boolean", s))
	}
	return StatusValue(ss & Status(def.mask()) >> Status(def.offset))
}

func (ss *Status) SetStatusValue(s StatusID, value StatusValue) {
	def := getStatusDefine(s)
	if def.isBoolean() {
		panic(fmt.Errorf("status %s is boolean", s))
	}
	*ss &= Status(value << StatusValue(def.offset))
}

func (ss *Status) Reset() {
	*ss = 0
}

func init() {
	totalBit := 0
	for i := 0; i < len(statusDefines); i++ {
		totalBit += statusDefines[i].bit
		if totalBit > 32 {
			panic(fmt.Errorf("robot-status %s bits overflow", statusDefines[i].name))
		}

		for j := i + 1; j < len(statusDefines); j++ {
			if statusDefines[i].IsOverlappedWith(statusDefines[j]) {
				panic(fmt.Errorf("the section of robot-status %s and %s are intersected", StatusID(i).String(), StatusID(j).String()))
			}
		}

		statusName2ID[statusDefines[i].name] = StatusID(i)
	}
}
