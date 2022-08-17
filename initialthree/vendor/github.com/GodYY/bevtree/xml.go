package bevtree

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/GodYY/gutils/assert"
	"github.com/pkg/errors"
)

const indent = "    "

// help function to create xml.Name in namespace bevtree.
func XMLName(name string) xml.Name {
	return xml.Name{Space: "", Local: name}
}

// XML strings.
const (

	// xml name for Tree.
	XMLStringTree = "bevtree"

	// xml name for name.
	XMLStringName = "name"

	// xml name for comment.
	XMLStringComment = "comment"

	// xml name for NodeType.
	XMLStringNodeType = "nodetype"

	// xml name for BevType
	XMLStringBevType = "bevtype"

	// xml name for root node.
	XMLStringRoot = "root"

	// xml name for childs.
	XMLStringChilds = "childs"

	// xml name for child.
	XMLStringChild = "child"

	// xml name for limited.
	XMLStringLimited = "limited"

	// xml name for success on fail.
	XMLStringSuccessOnFail = "successonfail"

	// xml name for subtree.
	XMLStringSubtree = "subtree"

	XMLStringConfig = "config"
)

func XMLNameToString(name xml.Name) string {
	if name.Space == "" {
		return name.Local
	} else {
		return name.Space + name.Local
	}
}

func XMLTokenToString(token xml.Token) string {
	var sb strings.Builder

	switch o := token.(type) {
	case xml.StartElement:
		nameStr := XMLNameToString(o.Name)
		sb.WriteString(fmt.Sprintf("<%s", nameStr))

		for _, attr := range o.Attr {
			sb.WriteString(fmt.Sprintf(" %s=\"%s\"", XMLNameToString(attr.Name), attr.Value))
		}

		sb.WriteString(">")

		return sb.String()

	case xml.EndElement:
		nameStr := XMLNameToString(o.Name)
		sb.WriteString(fmt.Sprintf("</%s>", nameStr))
		return sb.String()

	default:
		panic("not support yet")
	}
}

func XMLTokenError(token xml.Token, err error) error {
	return errors.WithMessage(err, XMLTokenToString(token))
}

func XMLTokenErrorf(token xml.Token, f string, args ...interface{}) error {
	return errors.Errorf(XMLTokenToString(token)+": "+f, args...)
}

func XMLAttrNotFoundError(attrName xml.Name) error {
	return errors.Errorf("attribute \"%s\" not found", XMLNameToString(attrName))
}

// XMLUnmarshal is the interface implemented by objects that can marshal
// themselves into valid bevtree XML elements.
type XMLMarshaler interface {
	MarshalBTXML(*XMLEncoder, xml.StartElement) error
}

// XMLUnmarshaler is the interface implemented by objects that can unmarshal
// an bevtree XML element description of themselves.
type XMLUnmarshaler interface {
	UnmarshalBTXML(*XMLDecoder, xml.StartElement) error
}

// An Encoder writes bevtree XML data to an output stream.
type XMLEncoder struct {
	*xml.Encoder
	framework *Framework
}

// newXMLEncoder returns a new encoder that writes to w.
func newXMLEncoder(framework *Framework, w io.Writer) *XMLEncoder {
	assert.Assert(framework != nil, "framework nil")
	assert.Assert(w != nil, "writer nil")

	return &XMLEncoder{
		Encoder:   xml.NewEncoder(w),
		framework: framework,
	}
}

// EncodeElement writes the bevtree XML encoding of v to the stream,
// using start as the outermost tag in the encoding.
//
// EncodeElement calls Flush before returning.
func (e *XMLEncoder) EncodeElement(v interface{}, start xml.StartElement) error {
	if marshaler, ok := v.(XMLMarshaler); ok {
		if err := marshaler.MarshalBTXML(e, start); err != nil {
			return err
		}

		return e.Flush()
	} else {
		return e.Encoder.EncodeElement(v, start)
	}
}

// EncodeSE writes the bevtree XML encoding of v to the stream
// between start and start.End(), using start as the outermost tag
// in the encoding.
//
// EncodeSE calls Flush before returning.
func (e *XMLEncoder) EncodeSE(start xml.StartElement, f func(*XMLEncoder) error) error {
	// if f == nil {
	// 	return errors.New("f nil")
	// }

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if f != nil {
		if err := f(e); err != nil {
			return err
		}
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}

	return e.Flush()
}

// EncodeElementSE writes the bevtree XML encoding of v to the stream
// between start and start.End(), using start as the outermost tag
// in the encoding.
//
// EncodeElementSE calls Flush before returning.
func (e *XMLEncoder) EncodeElementSE(v interface{}, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeElement(v, start); err != nil {
		return err
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}

	return e.Flush()
}

// EncodeNode encode the behavior tree node to the stream with start as
// start element. EncodeNode automatically encode the type of the node
// as xml.Attr and append it to start.
// <start ... nodetype="nodetype">
// ...
// </end>
func (e *XMLEncoder) EncodeNode(n Node, start xml.StartElement) error {
	if ntAttr, err := e.marshalNodeTypeAttr(n.NodeType(), (XMLName(XMLStringNodeType))); err == nil {
		start.Attr = append(start.Attr, ntAttr)
	} else {
		return err
	}

	if n.Comment() != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: XMLName(XMLStringComment), Value: n.Comment()})
	}

	return e.EncodeElement(n, start)
}

// Marshal nodeType as xml.Attr with name.
func (e *XMLEncoder) marshalNodeTypeAttr(nodeType NodeType, name xml.Name) (xml.Attr, error) {
	if meta := e.framework.getNodeMeta(nodeType); meta == nil {
		return xml.Attr{}, errors.Errorf("meta of node type \"%s\" not found", nodeType.String())
	} else {
		return xml.Attr{Name: name, Value: nodeType.String()}, nil
	}
}

// Marshal bevType as xml.Attr with name.
func (e *XMLEncoder) marshalBevTypeAttr(bevType BevType, name xml.Name) (xml.Attr, error) {
	if meta := e.framework.getBevMeta(bevType); meta == nil {
		return xml.Attr{}, errors.Errorf("meta of bev type \"%s\" not found", bevType.String())
	} else {
		return xml.Attr{Name: name, Value: bevType.String()}, nil
	}
}

// A XMLDecoder represents an bevtree XML parser reading a particular
// input stream.
type XMLDecoder struct {
	*xml.Decoder

	// The cached token for next decoding.
	tokenCached xml.Token

	// behavior tree system.
	framework *Framework
}

// newXMLDecoder creates a new bevtree XML parser reading from r.
// If r does not implement io.ByteReader, newXMLDecoder will
// do its own buffering.
func newXMLDecoder(framework *Framework, r io.Reader) *XMLDecoder {
	assert.Assert(framework != nil, "framework nil")
	assert.Assert(r != nil, "reader nil")

	return &XMLDecoder{
		Decoder:   xml.NewDecoder(r),
		framework: framework,
	}
}

func (d *XMLDecoder) Token() (xml.Token, error) {
	if d.tokenCached != nil {
		token := d.tokenCached
		d.tokenCached = nil
		return token, nil
	} else {
		return d.Decoder.Token()
	}
}

func (d *XMLDecoder) Skip() error {
	if d.tokenCached != nil {
		if _, ok := d.tokenCached.(xml.EndElement); ok {
			d.tokenCached = nil
			return nil
		}

		d.tokenCached = nil
	}

	return d.Decoder.Skip()
}

func (d *XMLDecoder) Framework() *Framework { return d.framework }

// DecodeElement read element from start to parse into v.
func (d *XMLDecoder) DecodeElement(v interface{}, start xml.StartElement) error {
	if unmarshal, ok := v.(XMLUnmarshaler); ok {
		return unmarshal.UnmarshalBTXML(d, start)
	} else {
		return d.Decoder.DecodeElement(v, &start)
	}
}

// DecodeAt search the start element with name at. If the
// start elment was found, DecodeAt invoke f with the start
// element and return the result of f.
func (d *XMLDecoder) DecodeAt(at xml.Name, f func(*XMLDecoder, xml.StartElement) error) error {
	var err error
	var token xml.Token

	for token, err = d.Token(); err == nil; token, err = d.Token() {
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name == at {
				if err = f(d, t); err != nil {
					return err
				} else {
					return nil
				}
			}

		case xml.EndElement:
			if t.Name == at {
				return fmt.Errorf("primarily found %s", XMLTokenToString(t))
			}
		}
	}

	return err
}

// The error used to stop decoding.
var ErrXMLDecodeStop = errors.New("XML decode stop")

// DecodeUntil invokes f with every start element, until
// until end element has been readed, or the result of f
// is non-nil or ErrXMLDecodeStop.
func (d *XMLDecoder) DecodeUntil(until xml.EndElement, f func(*XMLDecoder, xml.StartElement) error) error {
	var err error
	var token xml.Token

	for token, err = d.Token(); err == nil; token, err = d.Token() {
		switch t := token.(type) {
		case xml.StartElement:
			if err = f(d, t); err == nil {
				continue
			} else if err != ErrXMLDecodeStop {
				return err
			} else {
				return nil
			}

		case xml.EndElement:
			if t == until {
				d.tokenCached = t
				return nil
			}
		}
	}

	return err
}

// DecodeAtUntil invokes f with every start element named at,
// until until end element has been read, or the result of
// f is non-nil or ErrXMLDecodeStop.
func (d *XMLDecoder) DecodeAtUntil(at xml.Name, until xml.EndElement, f func(*XMLDecoder, xml.StartElement) error) error {
	var err error
	var token xml.Token
	found := false

	for token, err = d.Token(); err == nil; token, err = d.Token() {
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name == at {
				if found {
					return errors.Errorf("<%s...> self nested", XMLNameToString(at))
				} else {
					found = true

					if err = f(d, t); err == nil {
						found = false
						continue
					} else if err == ErrXMLDecodeStop {
						return nil
					} else {
						return err
					}
				}
			}

		case xml.EndElement:
			if t.Name == at {
				return fmt.Errorf("primarily found %s", XMLTokenToString(t))
			} else if t == until {
				d.tokenCached = t
				return nil
			}
		}
	}

	return err
}

// DecodeElementAt looks for the start element named at, and invoke
// DecodeElement with it.
func (d *XMLDecoder) DecodeElementAt(v interface{}, at xml.Name) error {
	return d.DecodeAt(at, func(d *XMLDecoder, s xml.StartElement) error {
		return d.DecodeElement(v, s)
	})
}

// DecodeElementUntil looks for the start element named name to decode
// to v, until the until end element has been read.
func (d *XMLDecoder) DecodeElementUntil(v interface{}, name xml.Name, until xml.EndElement) error {
	return d.DecodeUntil(until, func(d *XMLDecoder, t xml.StartElement) error {
		if t.Name == name {
			if err := d.DecodeElement(v, t); err != nil {
				return err
			} else {
				return ErrXMLDecodeStop
			}
		}

		return nil
	})
}

// DecodeNode decode the node type from start, then create the node
// with the type to decode into it.
func (d *XMLDecoder) DecodeNode(start xml.StartElement) (Node, error) {
	nodeTypeXMLName := XMLName(XMLStringNodeType)
	commentXMLName := XMLName(XMLStringComment)
	var node Node
	var comment string
	for _, attr := range start.Attr {
		if attr.Name == nodeTypeXMLName {
			if nodeType, err := d.unmarshalNodeTypeAttr(attr); err == nil {
				node = d.framework.getNodeMeta(nodeType).createNode()
			} else {
				return nil, XMLTokenError(start, err)
			}
		} else if attr.Name == commentXMLName {
			comment = attr.Value
		}
	}

	if node == nil {
		return nil, XMLTokenError(start, XMLAttrNotFoundError(nodeTypeXMLName))
	}

	node.SetComment(comment)

	if err := d.DecodeElement(node, start); err != nil {
		return nil, err
	}

	return node, nil
}

// DecodeNodeAt first looks for the start element named startName of node;
// then, parse attr NodeType from it and create new node; finally, invoke
// DecodeElement with the start element and new node.
func (d *XMLDecoder) DecodeNodeAt(pnode *Node, at xml.Name) error {
	return d.DecodeAt(at, func(d *XMLDecoder, s xml.StartElement) error {

		node, err := d.DecodeNode(s)
		if err != nil {
			return err
		} else {
			*pnode = node
			return nil
		}

	})
}

// DecodeNodeUntil works like DecodeUntil, except it works on behavior tree
// node.
func (d *XMLDecoder) DecodeNodeUntil(pnode *Node, name xml.Name, until xml.EndElement) error {
	return d.DecodeUntil(until, func(d *XMLDecoder, s xml.StartElement) error {

		node, err := d.DecodeNode(s)
		if err != nil {
			return err
		} else {
			*pnode = node
			return ErrXMLDecodeStop
		}

	})
}

// Unmarshal xml.Attr attr as NodeType.
func (d *XMLDecoder) unmarshalNodeTypeAttr(attr xml.Attr) (NodeType, error) {
	nt := NodeType(attr.Value)
	if meta := d.framework.getNodeMeta(nt); meta == nil {
		return NodeType(""), errors.Errorf("meta of node type %s not found", attr.Value)
	} else {
		return nt, nil
	}
}

// Unmarshal xml.Attr attr as BevType.
func (d *XMLDecoder) unmarshalBevTypeAttr(attr xml.Attr) (BevType, error) {
	bt := BevType(attr.Value)
	if meta := d.framework.getBevMeta(bt); meta == nil {
		return BevType(""), errors.Errorf("meta of bev type %s not found", attr.Value)
	} else {
		return bt, nil
	}
}

func MarshalXMLTree(framework *Framework, t *tree) ([]byte, error) {
	if framework == nil {
		return nil, errors.New("marshal xml tree: framework nil")
	}

	if t == nil {
		return nil, nil
	}

	var buf = bytes.NewBuffer(nil)
	e := newXMLEncoder(framework, buf)
	e.Indent("", indent)

	start := xml.StartElement{Name: XMLName(XMLStringTree)}
	if err := e.EncodeElement(t, start); err != nil {
		return nil, errors.WithMessagef(err, "marshal xml tree")
	}

	return buf.Bytes(), nil
}

func UnmarshalXMLTree(framework *Framework, data []byte, t *tree) error {
	if framework == nil {
		return errors.New("unmarshal xml tree: framework nil")
	}

	if data == nil || t == nil {
		return nil
	}

	var buf = bytes.NewReader(data)
	d := newXMLDecoder(framework, buf)

	if err := d.DecodeElementAt(t, XMLName(XMLStringTree)); err != nil {
		return errors.WithMessagef(err, "unmarshal xml tree")
	}

	return nil
}

func EncodeXMLTreeFile(framework *Framework, path string, t *tree) (err error) {
	if framework == nil {
		return errors.New("encode xml tree file: framework nil")
	}

	if t == nil {
		return nil
	}

	var file *os.File

	file, err = os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		if e := file.Close(); err == nil {
			err = e
		}
	}()

	enc := newXMLEncoder(framework, file)
	enc.Indent("", indent)

	start := xml.StartElement{Name: XMLName(XMLStringTree)}
	if err := enc.EncodeElement(t, start); err != nil {
		return errors.WithMessagef(err, "encode xml tree file: \"%s\"", path)
	}

	return nil
}

func DecodeXMLTreeFile(framework *Framework, path string, t *tree) error {
	if framework == nil {
		return errors.New("decode xml tree file: framework nil")
	}

	if t == nil {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := newXMLDecoder(framework, file)

	if err := dec.DecodeElementAt(t, XMLName(XMLStringTree)); err != nil {
		return errors.WithMessagef(err, "decode xml tree file: \"%s\"", path)
	}

	return nil
}

// MarshalXMLTree return an bevtree XML encoding of t.
func (f *Framework) MarshalXMLTree(t *tree) ([]byte, error) {
	if data, err := MarshalXMLTree(f, t); err != nil {
		return nil, errors.WithMessage(err, "framework")
	} else {
		return data, nil
	}
}

// UnmarshalXMLTree parses the bevtree XML-encoded Tree
// data and stores the result in the Tree pointed to by t.
func (f *Framework) UnmarshalXMLTree(data []byte, t *tree) error {
	if err := UnmarshalXMLTree(f, data, t); err != nil {
		return errors.WithMessage(err, "framework")
	} else {
		return nil
	}
}

// EncodeXMLTreeFile works like MarshalXMLTree but write
// encoded data to file.
func (f *Framework) EncodeXMLTreeFile(path string, t *tree) (err error) {
	if err := EncodeXMLTreeFile(f, path, t); err != nil {
		return errors.WithMessage(err, "framework")
	} else {
		return nil
	}
}

// DecodeXMLTreeFile works like UnmarshalXMLTree but read
// encoded data from file.
func (f *Framework) DecodeXMLTreeFile(path string, t *tree) error {
	if err := DecodeXMLTreeFile(f, path, t); err != nil {
		return errors.WithMessage(err, "framework")
	} else {
		return nil
	}
}

func (t *tree) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("Tree.MarshalBTXML start:%v", start)
	}

	if t.name == "" {
		return errors.New("Tree has no name")
	}

	if t.name != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: XMLName(XMLStringName), Value: t.name})
	}

	if t.comment != "" {
		start.Attr = append(start.Attr, xml.Attr{Name: XMLName(XMLStringComment), Value: t.comment})
	}

	if err := e.EncodeSE(start, func(x *XMLEncoder) error {
		rootStart := xml.StartElement{Name: XMLName(XMLStringRoot)}
		if err := e.EncodeElementSE(t._root, rootStart); err != nil {
			return errors.WithMessagef(err, "Marshal root")
		}

		return nil
	}); err != nil {
		return errors.WithMessagef(err, "Tree %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (t *tree) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("Tree.UnmarshalBTXML start:%v", start)
	}

	for _, attr := range start.Attr {
		if attr.Name == XMLName(XMLStringName) {
			t.name = attr.Value
		} else if attr.Name == XMLName(XMLStringComment) {
			t.comment = attr.Value
		}
	}

	if t.name == "" {
		return errors.New("tree has no name")
	}

	if t._root == nil {
		t._root = newRootNode()
	}

	if err := d.DecodeElementAt(t._root, XMLName(XMLStringRoot)); err != nil {
		return errors.WithMessagef(err, "Tree %s Unmarshal root", XMLTokenToString(start))
	}

	return d.Skip()
}

func (r *rootNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("rootNode.MarshalBTXML start:%v", start)
	}

	if r.child == nil {
		return nil
	}

	if err := e.EncodeNode(r.child, xml.StartElement{Name: XMLName(XMLStringChild)}); err != nil {
		return errors.WithMessage(err, "rootNode Marshal child")
	}

	return nil
}

func (r *rootNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("rootNode.UnmarshalBTXML start:%v", start)
	}

	var child Node
	if err := d.DecodeNodeUntil(&child, XMLName(XMLStringChild), start.End()); err != nil {
		return errors.WithMessage(err, "rootNode Unmarshal child")
	}

	if child != nil {
		r.SetChild(child)
	}

	return d.Skip()
}

func (d *decoratorNode) marshalXML(e *XMLEncoder) error {
	if d.child == nil {
		return nil
	}

	if err := e.EncodeNode(d.child, xml.StartElement{Name: XMLName(XMLStringChild)}); err != nil {
		return errors.WithMessage(err, "Marshal child")
	}

	return nil
}

func (d *decoratorNode) unmarshalXML(dec *XMLDecoder, start xml.StartElement) error {
	var child Node
	if err := dec.DecodeNodeUntil(&child, XMLName(XMLStringChild), start.End()); err != nil {
		return errors.WithMessage(err, "Unmarshal child")
	}

	if child != nil {
		d.setChild(child)
	}

	return nil
}

func (i *InverterNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("InverterNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {
		return i.decoratorNode.marshalXML(e)
	}); err != nil {
		return errors.WithMessagef(err, "InverterNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (i *InverterNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("InverterNode.UnmarshalBTXML start:%v", start)
	}

	if err := i.decoratorNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "InverterNode %s Unmarshal", XMLTokenToString(start))
	}

	if i.child != nil {
		i.child.SetParent(i)
	}

	return d.Skip()
}

func (s *SucceederNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("SucceederNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {
		return s.decoratorNode.marshalXML(e)
	}); err != nil {
		return errors.WithMessagef(err, "SucceederNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (s *SucceederNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("SucceederNode.UnmarshalBTXML start:%v", start)
	}

	if err := s.decoratorNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "SucceederNode %s Unmarshal", XMLTokenToString(start))
	}

	if s.child != nil {
		s.child.SetParent(s)
	}

	return d.Skip()
}

func (r *RepeaterNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("RepeaterNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {

		if err := e.EncodeElement(r.limited, xml.StartElement{Name: XMLName(XMLStringLimited)}); err != nil {
			return errors.WithMessage(err, "Marshal limited")
		}

		return r.decoratorNode.marshalXML(e)

	}); err != nil {
		return errors.WithMessagef(err, "RepeaterNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (r *RepeaterNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("RepeaterNode.UnmarshalBTXML start:%v", start)
	}

	if err := d.DecodeElementAt(&r.limited, XMLName(XMLStringLimited)); err != nil {
		return errors.WithMessagef(err, "RepeaterNode %s Unmarshal limited", XMLTokenToString(start))
	}

	if err := r.decoratorNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "RepeaterNode %s Unmarshal", XMLTokenToString(start))
	}

	if r.child != nil {
		r.child.SetParent(r)
	}

	return d.Skip()
}

func (r *RepeatUntilFailNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("RepeatUntilFailNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {

		if err := e.EncodeElement(r.successOnFail, xml.StartElement{Name: XMLName(XMLStringSuccessOnFail)}); err != nil {
			return errors.WithMessage(err, "Marshal successOnFail")
		}

		return r.decoratorNode.marshalXML(e)

	}); err != nil {
		return errors.WithMessagef(err, "RepeatUntilFailNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (r *RepeatUntilFailNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("RepeatUntilFailNode.UnmarshalBTXML start:%v", start)
	}

	if err := d.DecodeElementAt(&r.successOnFail, XMLName(XMLStringSuccessOnFail)); err != nil {
		return errors.WithMessagef(err, "RepeatUntilFailNode %s Unmarshal successOnFail", XMLTokenToString(start))
	}

	if err := r.decoratorNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "RepeatUntilFailNode %s Unmarshal", XMLTokenToString(start))
	}

	if r.child != nil {
		r.child.SetParent(r)
	}

	return d.Skip()
}

func (c *compositeNode) marshalXML(e *XMLEncoder) error {
	childCount := c.ChildCount()
	if childCount > 0 {
		var err error

		childsStart := xml.StartElement{Name: XMLName(XMLStringChilds)}
		childsStart.Attr = append(childsStart.Attr, xml.Attr{Name: XMLName("count"), Value: strconv.Itoa(childCount)})

		if err = e.EncodeToken(childsStart); err != nil {
			return err
		}

		for i := 0; i < childCount; i++ {
			child := c.Child(i)
			if err = e.EncodeNode(child, xml.StartElement{Name: XMLName(XMLStringChild)}); err != nil {
				return errors.WithMessagef(err, "Marshal No.%d child", i)
			}
		}

		if err = e.EncodeToken(childsStart.End()); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (c *compositeNode) unmarshalXML(d *XMLDecoder, start xml.StartElement) error {
	if err := d.DecodeAtUntil(XMLName(XMLStringChilds), start.End(), func(d *XMLDecoder, s xml.StartElement) error {
		xmlCountName := XMLName("count")
		var childCount int
		var childs []Node
		for _, attr := range s.Attr {
			if attr.Name == xmlCountName {
				var err error
				if childCount, err = strconv.Atoi(attr.Value); err != nil {
					return errors.WithMessage(err, "Unmarshal child count")
				} else if childCount <= 0 {
					return fmt.Errorf("invalid child count: %d", childCount)
				} else {
					childs = make([]Node, 0, childCount)
					break
				}
			}
		}

		if childCount > 0 {
			xmlChildName := XMLName(XMLStringChild)
			if err := d.DecodeAtUntil(xmlChildName, s.End(), func(d *XMLDecoder, s xml.StartElement) error {
				if len(childs) >= childCount {
					return errors.New("too many children")
				}

				if node, err := d.DecodeNode(s); err != nil {
					return err
				} else {
					childs = append(childs, node)
					return nil
				}

			}); err != nil {
				return errors.WithMessagef(err, "Unmarshal No.%d child", len(childs))
			}

			if len(childs) < childCount {
				return errors.New("too few children")
			}
		}

		c.children = childs

		if err := d.Skip(); err != nil {
			return err
		}

		return ErrXMLDecodeStop

	}); err != nil {
		return err
	}

	return nil
}

func (s *SequenceNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("SequenceNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {
		return s.compositeNode.marshalXML(e)
	}); err != nil {
		return errors.WithMessagef(err, "SequenceNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (s *SequenceNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("SequenceNode.UnmarshalBTXML start:%v", start)
	}

	if err := s.compositeNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "SequenceNode %s Unmarshal", XMLTokenToString(start))
	}

	for _, v := range s.children {
		v.SetParent(s)
	}

	return d.Skip()
}

func (s *SelectorNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("SelectorNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {
		return s.compositeNode.marshalXML(e)
	}); err != nil {
		return errors.WithMessagef(err, "SelectorNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (s *SelectorNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("SelectorNode.UnmarshalBTXML start:%v", start)
	}

	if err := s.compositeNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "SelectorNode %s Unmarshal", XMLTokenToString(start))
	}

	for _, v := range s.children {
		v.SetParent(s)
	}

	return d.Skip()
}

func (r *RandSequenceNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("RandSequenceNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {
		return r.compositeNode.marshalXML(e)
	}); err != nil {
		return errors.WithMessagef(err, "RandSequenceNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (r *RandSequenceNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("RandSequenceNode.UnmarshalBTXML start:%v", start)
	}

	if err := r.compositeNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "RandSequenceNode %s Unmarshal", XMLTokenToString(start))
	}

	for _, v := range r.children {
		v.SetParent(r)
	}

	return d.Skip()
}

func (r *RandSelectorNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("RandSelectorNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {
		return r.compositeNode.marshalXML(e)
	}); err != nil {
		return errors.WithMessagef(err, "RandSelectorNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (r *RandSelectorNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("RandSelectorNode.MarshalBTXML start:%v", start)
	}

	if err := r.compositeNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "RandSelectorNode %s Unmarshal", XMLTokenToString(start))
	}

	for _, v := range r.children {
		v.SetParent(r)
	}

	return d.Skip()
}

func (p *ParallelNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("ParallelNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(e *XMLEncoder) error {
		return p.compositeNode.marshalXML(e)
	}); err != nil {
		return errors.WithMessagef(err, "ParallelNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (p *ParallelNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("ParallelNode.MarshalBTXML start:%v", start)
	}

	if err := p.compositeNode.unmarshalXML(d, start); err != nil {
		return errors.WithMessagef(err, "ParallelNode %s Unmarshal", XMLTokenToString(start))
	}

	for _, v := range p.children {
		v.SetParent(p)
	}

	return d.Skip()
}

func (b *BevNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("BevNode.MarshalBTXML start:%v", start)
	}

	var err error
	var bevTypeAttr xml.Attr
	if bevTypeAttr, err = e.marshalBevTypeAttr(b.bev.BevType(), XMLName(XMLStringBevType)); err == nil {
		start.Attr = append(start.Attr, bevTypeAttr)

		switch o := b.bev.(type) {
		case XMLMarshaler:
			err = o.MarshalBTXML(e, start)
		default:
			err = e.Encoder.EncodeElement(b.bev, start)
		}

	}

	if err != nil {
		return errors.WithMessagef(err, "BevNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (b *BevNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("BevNode.UnmarshalBTXML start:%v", start)
	}

	var err error
	var bev Bev

	xmlNameBevType := XMLName(XMLStringBevType)
	for _, attr := range start.Attr {
		if attr.Name == xmlNameBevType {
			var bevType BevType
			if bevType, err = d.unmarshalBevTypeAttr(attr); err == nil {
				bev = d.Framework().getBevMeta(bevType).createBev()
				err = d.DecodeElement(bev, start)
			}

			break
		}
	}

	if err != nil {
		return errors.WithMessagef(err, "BevNode %s Unmarshal", XMLTokenToString(start))
	} else if bev == nil {
		return errors.Errorf("BevNode %s Unmarshal: no bev", XMLTokenToString(start))
	} else {
		b.bev = bev
		return nil
	}
}

func (s *SubtreeNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("SubtreeNode.MarshalBTXML start:%v", start)
	}

	start.Attr = append(start.Attr, xml.Attr{Name: XMLName(XMLStringSubtree), Value: s.subtree.Name()})

	if err := e.EncodeSE(start, nil); err != nil {
		return errors.WithMessagef(err, "SubtreeNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (s *SubtreeNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("SubtreeNode.UnmarshalBTXML start:%v", start)
	}

	for _, attr := range start.Attr {
		if attr.Name == XMLName(XMLStringSubtree) {
			subtree, err := d.Framework().getOrLoadTree(attr.Value)
			if err != nil {
				return errors.WithMessagef(err, "SubtreeNode %s Unmarshal", XMLTokenToString(start))
			} else if subtree == nil {
				return errors.Errorf("SubtreeNode %s Unmarshal: subtree \"%s\" not exsit", XMLTokenToString(start), attr.Value)
			} else {
				s.subtree = subtree
				break
			}
		}
	}

	if s.subtree == nil {
		return errors.Errorf("SubtreeNode %s Unmarshal: attr subtree not exist", XMLTokenToString(start))
	}

	return d.Skip()
}

func (n *WeightSelectorNode) MarshalBTXML(e *XMLEncoder, start xml.StartElement) error {
	if debug {
		log.Printf("WeightSelectorNode.MarshalBTXML start:%v", start)
	}

	if err := e.EncodeSE(start, func(x *XMLEncoder) error {
		childCount := n.ChildCount()
		if childCount == 0 {
			return nil
		}

		childsStart := xml.StartElement{Name: XMLName(XMLStringChilds)}
		childsStart.Attr = append(childsStart.Attr, xml.Attr{Name: XMLName("count"), Value: strconv.Itoa(childCount)})

		if err := e.EncodeToken(childsStart); err != nil {
			return err
		}

		for i := 0; i < childCount; i++ {
			child, weight := n.Child(i)
			childStart := xml.StartElement{Name: XMLName(XMLStringChild)}
			childStart.Attr = append(childStart.Attr, xml.Attr{Name: XMLName("weight"), Value: strconv.FormatFloat(float64(weight), 'f', -1, 32)})
			if err := e.EncodeNode(child, childStart); err != nil {
				return errors.WithMessagef(err, "Marshal No.%d child", i)
			}
		}

		if err := e.EncodeToken(childsStart.End()); err != nil {
			return err
		}

		return nil

	}); err != nil {
		return errors.WithMessagef(err, "WeightSelectorNode %s Marshal", XMLTokenToString(start))
	}

	return nil
}

func (n *WeightSelectorNode) UnmarshalBTXML(d *XMLDecoder, start xml.StartElement) error {
	if debug {
		log.Printf("WeightSelectorNode.UnmarshalBTXML start:%v", start)
	}

	if err := d.DecodeAtUntil(XMLName(XMLStringChilds), start.End(), func(d *XMLDecoder, s xml.StartElement) error {
		xmlCountName := XMLName("count")
		var childCount int
		var childs []*weightNode
		for _, attr := range s.Attr {
			if attr.Name == xmlCountName {
				var err error
				if childCount, err = strconv.Atoi(attr.Value); err != nil {
					return errors.WithMessage(err, "Unmarshal child count")
				} else if childCount <= 0 {
					return fmt.Errorf("invalid child count: %d", childCount)
				} else {
					childs = make([]*weightNode, 0, childCount)
					break
				}
			}
		}

		if childCount > 0 {
			xmlChildName := XMLName(XMLStringChild)
			if err := d.DecodeAtUntil(xmlChildName, s.End(), func(d *XMLDecoder, s xml.StartElement) error {
				if len(childs) >= childCount {
					return errors.New("too many children")
				}

				var weight = math.NaN()
				xmlNameWeight := XMLName("weight")
				for _, attr := range s.Attr {
					if attr.Name == xmlNameWeight {
						if w, err := strconv.ParseFloat(attr.Value, 32); err != nil {
							return errors.WithMessagef(err, "Unmarshal attr weight")
						} else {
							weight = w
						}
					}
				}

				if math.IsNaN(weight) {
					return errors.New("attr weight not found")
				}

				if node, err := d.DecodeNode(s); err != nil {
					return err
				} else {
					childs = append(childs, &weightNode{node: node, weight: float32(weight)})
					return nil
				}

			}); err != nil {
				return errors.WithMessagef(err, "Unmarshal No.%d child", len(childs))
			}

			if len(childs) < childCount {
				return errors.New("too few children")
			}
		}

		n.children = childs

		if err := d.Skip(); err != nil {
			return err
		}

		return ErrXMLDecodeStop

	}); err != nil {
		return errors.WithMessagef(err, "WeightSelectorNode %s Unmarshal", XMLTokenToString(start))
	}

	return d.Skip()
}
