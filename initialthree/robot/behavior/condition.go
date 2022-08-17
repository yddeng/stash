package behavior

import (
	"encoding/xml"
	"fmt"
	"initialthree/node/common/attr"
	"initialthree/robot/internal"
	"initialthree/robot/robot/module"
	"initialthree/robot/types"
	"initialthree/zaplogger"
	"strconv"
	"strings"

	. "github.com/GodYY/bevtree"
	"github.com/pkg/errors"
)

type StatusID = types.StatusID
type StatusValue = types.StatusValue

// 条件不满足的原因
type CondReason interface {
	String() string
}

// 条件检查结果
type CondResult struct {
	reasonNotMet CondReason
}

func (r CondResult) IsMet() bool { return r.reasonNotMet == nil }

func (r CondResult) ReasonNotMet() CondReason { return r.reasonNotMet }

// 条件类型
type CondType int8

const (
	// Cond_None   = CondType(iota)
	Cond_Status = CondType(iota)
	Cond_Attr
	cond_count
)

func (t CondType) String() string {
	return condDefines[t].name
}

func (t CondType) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if def, err := getCondDefine(t); err != nil {
		return xml.Attr{}, errors.WithMessage(err, "MarshalXMLAttr condition type")
	} else {
		return xml.Attr{Name: name, Value: def.name}, nil
	}
}

func (t *CondType) UnmarshalXMLAttr(attr xml.Attr) error {
	if tt, ok := condName2Types[attr.Value]; !ok {
		return fmt.Errorf("UnmarshalXMLAttr condition type: invalid condition type name %s", attr.Value)
	} else {
		*t = tt
		return nil
	}
}

type condDefine struct {
	name   string
	create func() Condition
}

var condDefines = make([]*condDefine, cond_count)

var condName2Types = map[string]CondType{}

func regCondDefine(cond Condition, def *condDefine) {
	t := cond.CondType()

	if def == nil {
		panic(fmt.Errorf("condition type %d define nil", t))
	}

	if condDefines[t] != nil {
		panic(fmt.Errorf("condition type %s define duplicate", t))
	}

	condDefines[t] = def
	condName2Types[def.name] = t
}

func getCondDefine(t CondType) (*condDefine, error) {
	if t < 0 || t >= cond_count {
		return nil, fmt.Errorf("invalid condition type %d", t)
	}

	if def := condDefines[t]; def == nil {
		return nil, fmt.Errorf("define of condition type %d not exist", t)
	} else {
		return def, nil
	}
}

func init() {
	regCondDefine(&CondStatus{}, &condDefine{name: "status", create: func() Condition { return new(CondStatus) }})
	regCondDefine(&CondAttr{}, &condDefine{name: "attr", create: func() Condition { return new(CondAttr) }})
}

func CreateCondition(ct CondType) (Condition, error) {
	if def, err := getCondDefine(ct); err != nil {
		return nil, err
	} else {
		return def.create(), nil
	}
}

type Condition interface {
	CondType() CondType
	CheckIsMet(player) CondResult
}

// type CondNone struct {
// }

// func (c CondNone) CondType() CondType {
// 	return Cond_None
// }

// func (c CondNone) CheckIsMet(r Robot) Result {
// 	return Result{}
// }

// 状态条件
type Status struct {
	status StatusID
	value  interface{}
	eq     bool
}

func NewStatusBool(status StatusID, b ...bool) Status {
	if !status.IsBoolean() {
		panic("not boolean status")
	}

	val := true
	if len(b) > 0 {
		val = b[0]
	}

	return Status{
		status: status,
		value:  val,
	}
}

func NewStatus(status StatusID, value StatusValue, eq ...bool) Status {
	if status.IsBoolean() {
		panic("boolean status")
	}

	equal := true
	if len(eq) > 0 {
		equal = eq[0]
	}

	return Status{
		status: status,
		value:  value,
		eq:     equal,
	}
}

func (s Status) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var err error
	var attr xml.Attr

	if attr, err = s.status.MarshalXMLAttr(internal.GetXMLName("name")); err != nil {
		return err
	} else {
		start.Attr = append(start.Attr, attr)
	}

	if s.status.IsBoolean() {
		if !s.value.(bool) {
			attr.Name = internal.GetXMLName("value")
			attr.Value = "false"
			start.Attr = append(start.Attr, attr)
		}
	} else if attr, err = s.value.(StatusValue).MarshalXMLAttr(internal.GetXMLName("value")); err != nil {
		return err
	} else {
		start.Attr = append(start.Attr, attr)
		if !s.eq {
			attr.Name = internal.GetXMLName("eq")
			attr.Value = "false"
			start.Attr = append(start.Attr, attr)
		}
	}

	if err = e.EncodeToken(start); err != nil {
		return err
	}

	if err = e.EncodeToken(start.End()); err != nil {
		return err
	}

	return nil
}

func (s *Status) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var err error
	var attrValue *xml.Attr

	s.eq = true

	for i, v := range start.Attr {
		if v.Name == internal.GetXMLName("name") {
			if err = s.status.UnmarshalXMLAttr(v); err != nil {
				return err
			}
		} else if v.Name == internal.GetXMLName("value") {
			attrValue = &start.Attr[i]
		} else if v.Name == internal.GetXMLName("eq") {
			if strings.ToLower(v.Value) == "true" {
				s.eq = true
			} else {
				s.eq = false
			}
		}
	}

	if s.status.IsBoolean() {
		if attrValue == nil || strings.ToLower(attrValue.Value) == "true" {
			s.value = true
		} else {
			s.value = false
		}
	} else {
		if attrValue == nil {
			return xml.UnmarshalError("miss value attr")
		}

		if err = s.status.UnmarshalXMLAttr(*attrValue); err != nil {
			return err
		}
	}

	return d.Skip()
}

func (s *Status) IsMet(r player) bool {
	if s.status.IsBoolean() {
		return s.value.(bool) == r.IsStatus(s.status)
	} else {
		b := s.value.(StatusValue) == r.GetStatusValue(s.status)
		return b && s.eq
	}
}

// 状态条件
type CondStatus struct {
	Or     bool     `xml:"or,attr,omitempty"`
	Status []Status `xml:"status"`
}

func NewCondStatus(or bool, status ...Status) *CondStatus {
	if len(status) == 0 {
		panic("no status")
	}

	return &CondStatus{
		Or:     or,
		Status: status,
	}
}

func (c *CondStatus) CondType() CondType {
	return Cond_Status
}

type reasonStatusNotMet struct {
	status Status
}

func (r reasonStatusNotMet) String() string {
	if r.status.status.IsBoolean() {
		if r.status.value.(bool) {
			return fmt.Sprintf("status \"%s\" is false", r.status.status)
		} else {
			return fmt.Sprintf("status \"%s\" is true", r.status.status)
		}
	} else {
		if r.status.eq {
			return fmt.Sprintf("status %s is not %s", r.status.status, r.status.value)
		} else {
			return fmt.Sprintf("status %s is %s", r.status.status, r.status.value)
		}
	}
}

func (c *CondStatus) CheckIsMet(r player) CondResult {
	for _, v := range c.Status {
		if !v.IsMet(r) {
			return CondResult{reasonNotMet: reasonStatusNotMet{status: v}}
		}
	}

	return CondResult{}
}

type Op int8

const (
	Lower = Op(iota)
	LowerEqual
	Equal
	Greater
	GreaterEqual
)

var opStrings = [...]string{
	Lower:        "lower",
	LowerEqual:   "lowerequal",
	Equal:        "equal",
	Greater:      "greater",
	GreaterEqual: "greaterequal",
}

var opStringToOp = map[string]Op{
	"lower":        Lower,
	"lowerequal":   LowerEqual,
	"equal":        Equal,
	"greater":      Greater,
	"greaterequal": GreaterEqual,
}

func (op Op) String() string { return opStrings[op] }

func (op Op) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if op < Lower || op > GreaterEqual {
		return xml.Attr{}, fmt.Errorf("invalid op: %d", op)
	}

	return xml.Attr{Name: name, Value: op.String()}, nil
}

func (op *Op) UnmarshalXMLAttr(attr xml.Attr) error {
	if o, ok := opStringToOp[attr.Value]; ok {
		*op = o
		return nil
	} else {
		return fmt.Errorf("invalid op: %s", attr.Value)
	}
}

type Attr struct {
	ID    int32 `xml:"id,attr"`
	Count int32 `xml:"count,attr"`
	Op    Op    `xml:"op,attr"`
}

func NewAttr(id, count int32, op Op) Attr {
	if id < 1 || id > attr.AttrMax {
		panic(fmt.Errorf("attr %d invalid", id))
	}

	if op < Lower || op > GreaterEqual {
		panic(fmt.Errorf("op %d invalid", op))
	}

	return Attr{
		ID:    id,
		Count: count,
		Op:    op,
	}
}

func (a Attr) IsMet(r player) bool {
	module := r.GetModule(module.Module_Attr).(*module.ModuleAttr)
	if attr := module.GetAttr(a.ID); attr == nil {
		zaplogger.GetSugar().Errorf("attr %d invalid", a.ID)
		return false
	} else {
		switch a.Op {
		case Lower:
			return int32(attr.GetVal()) < a.Count

		case LowerEqual:
			return int32(attr.GetVal()) <= a.Count

		case Equal:
			return int32(attr.GetVal()) == a.Count

		case Greater:
			return int32(attr.GetVal()) > a.Count

		case GreaterEqual:
			return int32(attr.GetVal()) >= a.Count

		default:
			panic(fmt.Errorf("invalid op: %d", a.Op))
		}
	}
}

// func (a Attr) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
// 	start.Attr = append(start.Attr, xml.Attr{Name: internal.GetXMLName("name"), Value: attr.GetNameById(a.ID)})
// 	start.Attr = append(start.Attr, xml.Attr{Name: internal.GetXMLName("count"), Value: strconv.Itoa(int(a.Count))})

// 	if err := e.EncodeToken(start); err != nil {
// 		return err
// 	}

// 	return e.EncodeToken(start.End())
// }

// func (a *Attr) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	a.ID = 0
// 	a.Count = 0

// 	for _, v := range start.Attr {
// 		if v.Name == internal.GetXMLName("name") {
// 			a.ID = attr.GetIdByName(v.Value)
// 			if a.ID <= 0 {
// 				return xml.UnmarshalError(fmt.Sprintf("invalid attrname \"%s\"", v.Value))
// 			}
// 		} else if v.Name == internal.GetXMLName("count") {
// 			if n, err := strconv.Atoi(v.Value); err != nil {
// 				return xml.UnmarshalError(fmt.Sprintf("invalid attrcount \"%s\"", v.Value))
// 			} else {
// 				a.Count = int32(n)
// 			}
// 		}
// 	}

// 	if a.ID == 0 {
// 		return xml.UnmarshalError("miss attrname")
// 	}

// 	if a.Count == 0 {
// 		return xml.UnmarshalError("miss attrcount")
// 	}

// 	return d.Skip()
// }

type CondAttr struct {
	Attr []Attr `xml:"attr"`
}

func NewCondAttr(attr ...Attr) *CondAttr {
	if len(attr) == 0 {
		panic("no attr")
	}

	return &CondAttr{Attr: attr}
}

func (c *CondAttr) CondType() CondType {
	return Cond_Attr
}

type reasonAttrNotMet struct {
	attr Attr
}

func (r reasonAttrNotMet) String() string {
	return fmt.Sprintf("attr %s %s %d", attr.GetNameById(r.attr.ID), r.attr.Op.String(), r.attr.Count)
}

func (c *CondAttr) CheckIsMet(r player) CondResult {
	for _, v := range c.Attr {
		if !v.IsMet(r) {
			return CondResult{reasonNotMet: reasonAttrNotMet{v}}
		}
	}

	return CondResult{}
}

const condition = NodeType("condition")

func init() {
	regNodeType(condition, func() Node { return new(ConditionNode) }, func() Task { return new(conditionTask) })
}

type ConditionNode struct {
	parent         Node
	noChildSuccess bool
	child          Node
	conds          []Condition
	comment        string
}

func NewConditionNode() *ConditionNode {
	return &ConditionNode{}
}

func (c *ConditionNode) NodeType() NodeType    { return condition }
func (c *ConditionNode) Parent() Node          { return c.parent }
func (c *ConditionNode) SetParent(parent Node) { c.parent = parent }

func (c *ConditionNode) SetNoChildSuccess(b bool) { c.noChildSuccess = b }

func (c *ConditionNode) Child() Node { return c.child }

func (c *ConditionNode) SetChild(child Node) {
	if child == nil || child.Parent() != nil {
		return
	}

	child.SetParent(c)
	c.child = child
}

func (c *ConditionNode) AddCond(cond Condition) {
	if cond == nil {
		return
	}
	c.conds = append(c.conds, cond)
}

func (c *ConditionNode) Cond(idx int) Condition {
	if idx < 0 || idx >= len(c.conds) {
		return nil
	}

	return c.conds[idx]
}

func (c *ConditionNode) RemoveCond(idx int) Condition {
	if idx < 0 || idx >= len(c.conds) {
		return nil
	}

	cond := c.conds[idx]
	c.conds = append(c.conds[:idx], c.conds[idx+1:]...)
	return cond
}

func (c *ConditionNode) SetComment(comment string) { c.comment = comment }
func (c *ConditionNode) Comment() string           { return c.comment }

func (c *ConditionNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Space: "", Local: "nochildsuccess"}})

	if err := e.EncodeSE(start, func(x *XMLEncoder) error {
		if len(c.conds) > 0 {
			condsStart := xml.StartElement{Name: xml.Name{Space: "", Local: "conds"}}
			condsStart.Attr = append(condsStart.Attr, xml.Attr{Name: xml.Name{Space: "", Local: "count"}, Value: strconv.Itoa(len(c.conds))})

			err := e.EncodeToken(condsStart)
			if err == nil {
				condStart := xml.StartElement{Name: xml.Name{Space: "", Local: "cond"}}
				for i, cond := range c.conds {
					attr, err := cond.CondType().MarshalXMLAttr(xml.Name{Space: "", Local: "condtype"})
					if err == nil {
						condStart.Attr = []xml.Attr{attr}
						err = e.EncodeElement(cond, condStart)
					}

					if err != nil {
						return errors.WithMessagef(err, "Marshal No.%d cond", i)
					}
				}

				err = e.EncodeToken(condsStart.End())
			}

			if err != nil {
				return err
			}
		}

		if c.child != nil {
			return e.EncodeNode(c.child, xml.StartElement{Name: xml.Name{Space: "", Local: "child"}})
		}

		return nil
	}); err != nil {
		return errors.WithMessage(err, "ConditionNode Marshal")
	}

	return nil
}

func (c *ConditionNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "nochildsuccess" {
			c.noChildSuccess = true
		}
	}

	if err := d.DecodeAtUntil(xml.Name{Space: "", Local: "conds"}, start.End(), func(d *XMLDecoder, s xml.StartElement) error {
		var xmlCondCountName = xml.Name{Space: "", Local: "count"}
		var condCount int
		var conds []Condition

		for _, attr := range s.Attr {
			if attr.Name == xmlCondCountName {
				var err error
				if condCount, err = strconv.Atoi(attr.Value); err != nil {
					return errors.WithMessage(err, "Unmarshal cond count")
				} else if condCount <= 0 {
					return fmt.Errorf("invalid cond count %d", condCount)
				} else {
					conds = make([]Condition, 0, condCount)
					break
				}
			}
		}

		if conds == nil {
			return XMLAttrNotFoundError(xmlCondCountName)
		}

		xmlCondname := xml.Name{Space: "", Local: "cond"}

		if err := d.DecodeAtUntil(xmlCondname, s.End(), func(d *XMLDecoder, s xml.StartElement) error {
			if len(conds) >= condCount {
				return errors.New("too many cond")
			}

			xmlCondTypeName := xml.Name{Space: "", Local: "condtype"}
			var err error
			var cond Condition

			for _, attr := range s.Attr {
				if attr.Name == xmlCondTypeName {
					var condType CondType
					if err = condType.UnmarshalXMLAttr(attr); err == nil {
						cond, err = CreateCondition(condType)
					}
					break
				}
			}

			if cond != nil {
				err = d.DecodeElement(cond, s)
			}

			if err != nil {
				return err
			}

			conds = append(conds, cond)
			return nil

		}); err != nil {
			return errors.WithMessagef(err, "Unmarshal No.%d cond", len(conds))
		}

		c.conds = conds

		if err := d.Skip(); err != nil {
			return err
		}

		return ErrXMLDecodeStop

	}); err != nil {
		return errors.WithMessagef(err, "ConditionNode Unmarshal %s conds", start)
	}

	var child Node

	if err := d.DecodeNodeUntil(&child, xml.Name{Space: "", Local: "child"}, start.End()); err != nil {
		return errors.WithMessagef(err, "ConditionNode Unmarshal %s child", start)
	}

	c.SetChild(child)

	return d.Skip()
}

type conditionTask struct {
	node *ConditionNode
}

func (c *conditionTask) TaskType() TaskType {
	return Serial
}

func (c *conditionTask) OnCreate(node Node) {
	c.node = node.(*ConditionNode)
}

func (c *conditionTask) OnDestroy() {
	c.node = nil
}

func (c *conditionTask) OnInit(nextNodes NodeList, ctx Context) bool {
	entity := ctx.UserData().(player)

	for i, cond := range c.node.conds {
		if result := cond.CheckIsMet(entity); !result.IsMet() {
			entity.Debugf("No.%d cond of ConditionNode \"%s\" not met: %s", i, c.node.Comment(), result.ReasonNotMet())
			return false
		}
	}

	if c.node.Child() == nil {
		if c.node.noChildSuccess {
			return true
		} else {
			return false
		}
	}

	nextNodes.PushNode(c.node.Child())
	return true
}

func (c *conditionTask) OnUpdate(ctx Context) Result {
	if c.node.Child() == nil && c.node.noChildSuccess {
		return Success
	} else {
		return Running
	}
}

func (c *conditionTask) OnTerminate(ctx Context) {}

func (c *conditionTask) OnChildTerminated(result Result, nextNodes NodeList, ctx Context) Result {
	return result
}
