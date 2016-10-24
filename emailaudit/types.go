package emailaudit

import (
	"encoding/xml"
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04"

// MonitorRequest MonitorRequest
type MonitorRequest struct {
	DomainName     string
	SourceUserName string
	DestUserName   string
	BeginDate      *time.Time
	EndDate        *time.Time
	MonitorLevels  MonitorRequestMonitorLevels
}

// MonitorRequestMonitorLevels MonitorRequestMonitorLevels
type MonitorRequestMonitorLevels struct {
	IncomingEmail MonitorLevel
	OutgoingEmail MonitorLevel
	Draft         MonitorLevel
	Chat          MonitorLevel
}

// NewMonitorRequest returns new MonitorRequest
func NewMonitorRequest(domainName string, sourceUserName string, destUserName string, endDate *time.Time, monitorLevels MonitorRequestMonitorLevels) MonitorRequest {
	m := MonitorRequest{
		DomainName:     domainName,
		SourceUserName: sourceUserName,
		DestUserName:   destUserName,
		EndDate:        endDate,
		MonitorLevels:  monitorLevels,
	}
	return m
}

func (req *MonitorRequest) monitorProperties() monitorProperties {
	m := monitorProperties{}
	m.addProperty("destUserName", req.DestUserName)
	m.addProperty("endDate", req.EndDate)

	if req.MonitorLevels.IncomingEmail != "" {
		m.addProperty("incomingEmailMonitorLevel", req.MonitorLevels.IncomingEmail)
	}
	if req.MonitorLevels.OutgoingEmail != "" {
		m.addProperty("outgoingEmailMonitorLevel", req.MonitorLevels.OutgoingEmail)
	}
	if req.MonitorLevels.Draft != "" {
		m.addProperty("draftMonitorLevel", req.MonitorLevels.Draft)
	}
	if req.MonitorLevels.Chat != "" {
		m.addProperty("chatMonitorLevel", req.MonitorLevels.Chat)
	}
	if req.BeginDate != nil {
		m.addProperty("beginDate", req.BeginDate)
	}

	return m
}

// ToXML convert monitor request to xml
func (req *MonitorRequest) ToXML() []byte {
	x, _ := xml.MarshalIndent(req.monitorProperties(), "", "  ")
	return x
}

type monitorProperties struct {
	XMLName       xml.Name      `xml:"http://www.w3.org/2005/Atom entry"`
	ID            string        `xml:"id,omitempty"`
	Updated       *time.Time    `xml:"updated,omitempty"`
	AppProperties []appProperty `xml:"apps:property"`
	Links         []link        `xml:"link"`
}

func (m *monitorProperties) addProperty(name string, value interface{}) {
	if date, ok := value.(*time.Time); ok {
		value = date.UTC().Format(timeFormat)
	}
	m.AppProperties = append(m.AppProperties, appProperty{Name: name, Value: fmt.Sprintf("%v", value)})
}

type appProperty struct {
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}

type link struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr"`
}
