package emailaudit

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04"
	urlPrefix  = "https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor"
)

// MailMonitor MailMonitor
type MailMonitor struct {
	DomainName     string
	SourceUserName string
	DestUserName   string
	BeginDate      *time.Time
	EndDate        *time.Time
	MonitorLevels  MailMonitorLevels
}

// MailMonitorLevels MailMonitorLevels
type MailMonitorLevels struct {
	IncomingEmail MailMonitorLevel
	OutgoingEmail MailMonitorLevel
	Draft         MailMonitorLevel
	Chat          MailMonitorLevel
}

// NewMailMonitor returns new MailMonitor
func NewMailMonitor(domainName string, sourceUserName string, destUserName string, endDate *time.Time, monitorLevels MailMonitorLevels) MailMonitor {
	m := MailMonitor{
		DomainName:     domainName,
		SourceUserName: sourceUserName,
		DestUserName:   destUserName,
		EndDate:        endDate,
		MonitorLevels:  monitorLevels,
	}
	return m
}

func (req *MailMonitor) monitorWriteProperties() monitorWriteProperties {
	m := monitorWriteProperties{}
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

func (req *MailMonitor) toXML() []byte {
	x, _ := xml.MarshalIndent(req.monitorWriteProperties(), "", "  ")
	return x
}

// URL returns URL
func (req *MailMonitor) URL() string {
	return fmt.Sprintf("%v/%v/%v", urlPrefix, req.DomainName, req.SourceUserName)
}

func monitorFromXML(data []byte) (*MailMonitor, error) {
	var v monitorReadProperties
	if err := xml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	m := MailMonitor{
		MonitorLevels: v.toMonitorLevels(),
	}
	for _, p := range v.AppProperties {
		switch p.Name {
		case "destUserName":
			m.DestUserName = p.Value
			break
		}
	}
	if strings.HasPrefix(v.ID, urlPrefix) {
		urlparts := strings.Split(strings.Replace(v.ID, urlPrefix+"/", "", 1), "/")
		m.DomainName = urlparts[0]
		m.SourceUserName = urlparts[1]
	}
	return &m, nil
}

func (m monitorReadProperties) toMonitorLevels() MailMonitorLevels {
	ret := MailMonitorLevels{}
	for _, p := range m.AppProperties {
		l := MailMonitorLevel(p.Value)
		switch p.Name {
		case "incomingEmailMonitorLevel":
			ret.IncomingEmail = l
			break
		case "outgoingEmailMonitorLevel":
			ret.OutgoingEmail = l
			break
		case "draftMonitorLevel":
			ret.Draft = l
			break
		case "chatMonitorLevel":
			ret.Chat = l
			break
		}
	}
	return ret
}

type monitorReadProperties struct {
	XMLName       xml.Name      `xml:"http://www.w3.org/2005/Atom entry,omitempty"`
	ID            string        `xml:"id,omitempty"`
	Updated       *time.Time    `xml:"updated,omitempty"`
	AppProperties []appProperty `xml:"http://schemas.google.com/apps/2006 property"`
	Links         []link
}

type monitorWriteProperties struct {
	XMLName       xml.Name      `xml:"http://www.w3.org/2005/Atom entry,omitempty"`
	ID            string        `xml:"id,omitempty"`
	Updated       *time.Time    `xml:"updated,omitempty"`
	AppProperties []appProperty `xml:"apps:property"`
	Links         []link
}

func (m *monitorWriteProperties) addProperty(name string, value interface{}) {
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
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr,omitempty"`
	Href    string   `xml:"href,attr"`
}
