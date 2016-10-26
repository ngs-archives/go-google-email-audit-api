package emailaudit

import (
	"testing"
	"time"
)

func TestMailMonitorToXML(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	beginDate := time.Date(2016, time.September, 1, 0, 0, 0, 0, loc)
	endDate := time.Date(2016, time.October, 30, 23, 59, 59, 0, loc)
	for _, test := range []struct {
		monitor     MailMonitor
		expectedXML string
	}{
		{NewMailMonitor("littleapps.co.jp", "src", "dest", &endDate, MailMonitorLevels{
			IncomingEmail: NoneLevel,
			OutgoingEmail: NoneLevel,
			Draft:         NoneLevel,
			Chat:          NoneLevel,
		}),
			`<atom:entry xmlns:atom="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <apps:property name="destUserName" value="dest"></apps:property>
  <apps:property name="endDate" value="2016-10-30 14:59"></apps:property>
</atom:entry>`},
		{func() MailMonitor {
			m := NewMailMonitor("littleapps.co.jp", "src", "dest", &endDate, MailMonitorLevels{
				IncomingEmail: HeaderOnlyLevel,
				OutgoingEmail: HeaderOnlyLevel,
				Draft:         HeaderOnlyLevel,
				Chat:          HeaderOnlyLevel,
			})
			return m
		}(),
			`<atom:entry xmlns:atom="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <apps:property name="destUserName" value="dest"></apps:property>
  <apps:property name="endDate" value="2016-10-30 14:59"></apps:property>
  <apps:property name="incomingEmailMonitorLevel" value="HEADER_ONLY"></apps:property>
  <apps:property name="outgoingEmailMonitorLevel" value="HEADER_ONLY"></apps:property>
  <apps:property name="draftMonitorLevel" value="HEADER_ONLY"></apps:property>
  <apps:property name="chatMonitorLevel" value="HEADER_ONLY"></apps:property>
</atom:entry>`},
		{func() MailMonitor {
			m := NewMailMonitor("littleapps.co.jp", "src", "dest", &endDate, MailMonitorLevels{
				IncomingEmail: FullMessageLevel,
				OutgoingEmail: FullMessageLevel,
				Draft:         FullMessageLevel,
				Chat:          FullMessageLevel,
			})
			m.BeginDate = &beginDate
			return m
		}(),
			`<atom:entry xmlns:atom="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <apps:property name="destUserName" value="dest"></apps:property>
  <apps:property name="endDate" value="2016-10-30 14:59"></apps:property>
  <apps:property name="incomingEmailMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="outgoingEmailMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="draftMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="chatMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="beginDate" value="2016-08-31 15:00"></apps:property>
</atom:entry>`},
	} {
		x := string(test.monitor.toXML())
		if x != test.expectedXML {
			t.Errorf(`Expected "%v" but got "%v"`, test.expectedXML, x)
		}
	}
}

func TestMailMonitorURL(t *testing.T) {
	m := MailMonitor{
		SourceUserName: "src",
		DomainName:     "example.com",
	}
	expected := "https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/src"
	actual := m.URL()
	if expected != actual {
		t.Errorf(`Expected "%v" but got "%v"`, expected, actual)
	}
}

const monitorXML = `<entry xmlns='http://www.w3.org/2005/Atom' xmlns:apps='http://schemas.google.com/apps/2006'>
  <id>https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek/namrata</id>
  <updated>2009-08-20T00:28:57.319Z</updated>
  <link rel='self' type='application/atom+xml' href="https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek/namrata" />
  <link rel='edit' type='application/atom+xml' href="https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek/namrata" />
  <apps:property name="destUserName" value="namrata"></apps:property>
  <apps:property name="endDate" value="2016-10-30 14:59"></apps:property>
  <apps:property name="incomingEmailMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="outgoingEmailMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="draftMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="chatMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="beginDate" value="2016-08-31 15:00"></apps:property>
</entry>`

func TestMailMonitorFromXML(t *testing.T) {
	m, err := monitorFromXML([]byte(monitorXML))
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	_TestMonitor(m, t)
}

func _TestMonitor(m *MailMonitor, t *testing.T) {
	for _, test := range []struct {
		actual   interface{}
		expected interface{}
	}{
		{m.MonitorLevels.Chat, FullMessageLevel},
		{m.MonitorLevels.Draft, FullMessageLevel},
		{m.MonitorLevels.IncomingEmail, FullMessageLevel},
		{m.MonitorLevels.OutgoingEmail, FullMessageLevel},
		{m.DestUserName, "namrata"},
		{m.SourceUserName, "abhishek"},
		{m.DomainName, "example.com"},
		{m.EndDate.UnixNano(), time.Date(2016, time.October, 30, 14, 59, 0, 0, time.UTC).UnixNano()},
		{m.BeginDate.UnixNano(), time.Date(2016, time.August, 31, 15, 0, 0, 0, time.UTC).UnixNano()},
		{m.Updated.UnixNano(), time.Date(2009, time.August, 20, 0, 28, 57, 319000000, time.UTC).UnixNano()},
	} {
		if test.actual != test.expected {
			t.Errorf(`Expected "%v" but got "%v"`, test.expected, test.actual)
		}
	}
}

const monitorsXML = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:openSearch="http://a9.com/-/spec/opensearchrss/1.0/" xmlns:apps="http://schemas.google.com/apps/2006">
<id>https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek</id>
<updated>2016-10-26T03:03:03.192Z</updated>
<link rel="http://schemas.google.com/g/2005#feed" type="application/atom+xml" href="https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek"/>
<link rel="http://schemas.google.com/g/2005#post" type="application/atom+xml" href="https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek"/>
<link rel="self" type="application/atom+xml" href="https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek"/>
<openSearch:startIndex>1</openSearch:startIndex>
<entry>
    <id>https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek/namrata</id>
    <updated>2009-04-17T15:29:21.064Z</updated>
    <link rel="self" type="application/atom+xml" href="https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek/namrata"/>
    <link rel="edit" type="application/atom+xml" href="https://apps-apis.google.com/feeds/compliance/audit/mail/monitor/example.com/abhishek/namrata"/>
    <apps:property name="requestId" value="53156"/>
    <apps:property name="destUserName" value="namrata"/>
    <apps:property name="beginDate" value="2009-06-15 00:00"/>
    <apps:property name="endDate" value="2009-06-30 23:20"/>
    <apps:property name="incomingEmailMonitorLevel" value="FULL_MESSAGE"/>
    <apps:property name="outgoingEmailMonitorLevel" value="FULL_MESSAGE"/>
    <apps:property name="draftMonitorLevel" value="FULL_MESSAGE"/>
    <apps:property name="chatMonitorLevel" value="FULL_MESSAGE"/>
 </entry>
 <entry>
    <id>https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek/joe</id>
    <updated>2009-05-17T15:29:21.064Z</updated>
    <link rel="self" type="application/atom+xml" href="https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek/joe"/>
    <link rel="edit" type="application/atom+xml" href="https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek/joe"/>
    <apps:property name="requestId" value="22405"/>
    <apps:property name="destUserName" value="joe"/>
    <apps:property name="beginDate" value="2009-06-20 00:00"/>
    <apps:property name="endDate" value="2009-07-30 23:20"/>
    <apps:property name="incomingEmailMonitorLevel" value="FULL_MESSAGE"/>
    <apps:property name="outgoingEmailMonitorLevel" value="FULL_MESSAGE"/>
    <apps:property name="draftMonitorLevel" value="FULL_MESSAGE"/>
    <apps:property name="chatMonitorLevel" value="FULL_MESSAGE"/>
  </entry>
</feed>
`

func _TestMonitors(m []MailMonitor, t *testing.T) {
	for _, test := range []struct {
		actual   interface{}
		expected interface{}
	}{
		{len(m), 2},
		{m[0].MonitorLevels.Chat, FullMessageLevel},
		{m[0].MonitorLevels.Draft, FullMessageLevel},
		{m[0].MonitorLevels.IncomingEmail, FullMessageLevel},
		{m[0].MonitorLevels.OutgoingEmail, FullMessageLevel},
		{m[0].DestUserName, "namrata"},
		{m[0].SourceUserName, "abhishek"},
		{m[0].DomainName, "example.com"},
		{m[0].BeginDate.String(), time.Date(2009, time.June, 15, 0, 0, 0, 0, time.UTC).String()},
		{m[0].EndDate.String(), time.Date(2009, time.June, 30, 23, 20, 0, 0, time.UTC).String()},
		{m[0].Updated.String(), time.Date(2009, time.April, 17, 15, 29, 21, 64000000, time.UTC).String()},
		{m[1].MonitorLevels.Chat, FullMessageLevel},
		{m[1].MonitorLevels.Draft, FullMessageLevel},
		{m[1].MonitorLevels.IncomingEmail, FullMessageLevel},
		{m[1].MonitorLevels.OutgoingEmail, FullMessageLevel},
		{m[1].DestUserName, "joe"},
		{m[1].SourceUserName, "abhishek"},
		{m[1].DomainName, "example.com"},
		{m[1].BeginDate.String(), time.Date(2009, time.June, 20, 0, 0, 0, 0, time.UTC).String()},
		{m[1].EndDate.String(), time.Date(2009, time.July, 30, 23, 20, 0, 0, time.UTC).String()},
		{m[1].Updated.String(), time.Date(2009, time.May, 17, 15, 29, 21, 64000000, time.UTC).String()},
	} {
		if test.actual != test.expected {
			t.Errorf(`Expected "%v" but got "%v"`, test.expected, test.actual)
		}
	}
}

func TestMailMonitorsFromXML(t *testing.T) {
	m, err := monitorsFromXML([]byte(monitorsXML))
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	_TestMonitors(m, t)
}

func TestMonitorFromXMLError(t *testing.T) {
	x := []byte("<foo />")
	m, err := monitorFromXML([]byte(x))
	expected := "expected element type <entry> but have <foo>"
	if err.Error() != expected {
		t.Errorf(`Expected "%v" but got "%v"`, expected, err.Error())
	}
	if m != nil {
		t.Errorf("Expected nil but got %v", m)
	}
}

func TestMonitorsFromXMLError(t *testing.T) {
	x := []byte("<foo />")
	m, err := monitorsFromXML([]byte(x))
	expected := "expected element type <feed> but have <foo>"
	if err.Error() != expected {
		t.Errorf(`Expected "%v" but got "%v"`, expected, err.Error())
	}
	if len(m) != 0 {
		t.Errorf(`Expected 0 but got %v`, len(m))
	}
	if m != nil {
		t.Errorf("Expected nil but got %v", m)
	}
}
