package internal

import "encoding/xml"

const InvalidID = -1

func GetXMLName(name string) xml.Name {
	return xml.Name{Space: "", Local: name}
}
