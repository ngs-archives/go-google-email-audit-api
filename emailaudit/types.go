package emailaudit

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04"
	baseURL    = "https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor"
)

// MailMonitor MailMonitor
type MailMonitor struct {
	DomainName     string
	SourceUserName string
	DestUserName   string
	BeginDate      *time.Time
	EndDate        *time.Time
	MonitorLevels  MailMonitorLevels
	Updated        *time.Time
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
	xstr := string(x)
	xstr = strings.Replace(xstr, `<entry xmlns="http://www.w3.org/2005/Atom">`, `<atom:entry xmlns:atom="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">`, 1)
	xstr = strings.Replace(xstr, `</entry>`, `</atom:entry>`, 1)
	return []byte(xstr)
}

// URL returns URL
func (req *MailMonitor) URL() string {
	return fmt.Sprintf("%v/%v/%v", baseURL, req.DomainName, req.SourceUserName)
}

func monitorFromXML(data []byte) (*MailMonitor, error) {
	var v monitorReadProperties
	if err := xml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	m := v.toMonitor()
	return &m, nil
}

func monitorsFromXML(data []byte) ([]MailMonitor, error) {
	var list monitorReadListProperties
	var entries []MailMonitor
	if err := xml.Unmarshal(data, &list); err != nil {
		return entries, err
	}
	for _, v := range list.Entries {
		entries = append(entries, v.toMonitor())
	}
	return entries, nil
}

func (m monitorReadProperties) toMonitor() MailMonitor {
	mm := MailMonitor{
		MonitorLevels: m.toMonitorLevels(),
	}
	for _, p := range m.AppProperties {
		switch p.Name {
		case "destUserName":
			mm.DestUserName = p.Value
			break
		case "beginDate":
			d, _ := time.Parse(timeFormat, p.Value)
			mm.BeginDate = &d
			break
		case "endDate":
			d, _ := time.Parse(timeFormat, p.Value)
			mm.EndDate = &d
			break
		}
	}
	if strings.HasPrefix(m.ID, baseURL) {
		urlparts := strings.Split(strings.Replace(m.ID, baseURL+"/", "", 1), "/")
		mm.DomainName = urlparts[0]
		mm.SourceUserName = urlparts[1]
	}
	mm.Updated = m.Updated
	return mm
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

type monitorReadListProperties struct {
	XMLName xml.Name                `xml:"http://www.w3.org/2005/Atom feed,omitempty"`
	Entries []monitorReadProperties `xml:"entry"`
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
	AppProperties []appProperty `xml:"apps:property"`
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
