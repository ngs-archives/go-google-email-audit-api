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
		{NewMailMonitor("littleapps.co.jp", "src", "dest", &endDate, MailMonitorLevels{}),
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

func TestMailMonitorFromXML(t *testing.T) {
	x := `<entry xmlns='http://www.w3.org/2005/Atom' xmlns:apps='http://schemas.google.com/apps/2006'>
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
	m, err := monitorFromXML([]byte(x))
	for _, test := range []struct {
		actual   MailMonitorLevel
		expected MailMonitorLevel
	}{
		{m.MonitorLevels.Chat, FullMessageLevel},
		{m.MonitorLevels.Draft, FullMessageLevel},
		{m.MonitorLevels.IncomingEmail, FullMessageLevel},
		{m.MonitorLevels.OutgoingEmail, FullMessageLevel},
	} {
		if test.actual != test.expected {
			t.Errorf(`Expected "%v" but got "%v"`, test.expected, test.actual)
		}
	}
	for _, test := range []struct {
		actual   string
		expected string
	}{
		{m.DestUserName, "namrata"},
		{m.SourceUserName, "abhishek"},
		{m.DomainName, "example.com"},
	} {
		if test.actual != test.expected {
			t.Errorf(`Expected "%v" but got "%v"`, test.expected, test.actual)
		}
	}
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
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
