package emailaudit

import (
	"testing"
	"time"
)

func TestMonitorRequestToXML(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	beginDate := time.Date(2016, time.September, 1, 0, 0, 0, 0, loc)
	endDate := time.Date(2016, time.October, 30, 23, 59, 59, 0, loc)
	m := NewMonitorRequest("littleapps.co.jp", "src", "dest", &endDate, MonitorRequestMonitorLevels{
		IncomingEmail: FullMessageLevel,
		OutgoingEmail: FullMessageLevel,
		Draft:         FullMessageLevel,
		Chat:          FullMessageLevel,
	})
	m.BeginDate = &beginDate

	x := string(m.ToXML())
	expected := `<entry xmlns="http://www.w3.org/2005/Atom">
  <apps:property name="destUserName" value="dest"></apps:property>
  <apps:property name="endDate" value="2016-10-30 14:59"></apps:property>
  <apps:property name="incomingEmailMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="outgoingEmailMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="draftMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="chatMonitorLevel" value="FULL_MESSAGE"></apps:property>
  <apps:property name="beginDate" value="2016-08-31 15:00"></apps:property>
</entry>`
	if x != expected {
		t.Errorf(`Expected "%v" but got "%v"`, expected, x)
	}
}
