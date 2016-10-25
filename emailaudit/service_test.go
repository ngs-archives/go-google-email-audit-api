package emailaudit

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"

	gock "gopkg.in/h2non/gock.v1"
)

func TestNewServiceError(t *testing.T) {
	svc, err := New(nil)
	expected := "client is nil"
	if err.Error() != expected {
		t.Errorf(`Expected %v but got "%v"`, expected, err.Error())
	}
	if svc != nil {
		t.Errorf("Expected nil but got %v", svc)
	}
}

func TestNewService(t *testing.T) {
	client := &http.Client{}
	svc, err := New(client)
	if err != nil {
		t.Errorf("Expected nil but got %v", svc)
	}
	if svc == nil {
		t.Errorf("Expected not nil but got %v", svc)
	}
}

func TestServiceUserAgent(t *testing.T) {
	client := &http.Client{}
	svc, _ := New(client)
	expected := "google-api-go-client/0.5"
	if svc.userAgent() != expected {
		t.Errorf(`Expected "%v" but got "%v"`, expected, svc.userAgent())
	}
	svc, _ = New(client)
	expected = "google-api-go-client/0.5 foo"
	svc.UserAgent = "foo"
	if svc.userAgent() != expected {
		t.Errorf(`Expected "%v" but got "%v"`, expected, svc.userAgent())
	}
}

func updateEmailMonitor() (*MailMonitor, error) {
	ctx := context.Background()
	config := &oauth2.Config{}
	token := &oauth2.Token{AccessToken: "test"}
	client := config.Client(ctx, token)
	svc, _ := New(client)
	endDate := time.Date(2016, time.October, 30, 14, 59, 0, 0, time.UTC)
	monitor := NewMailMonitor("example.com", "abhishek", "namrata", &endDate, MailMonitorLevels{
		IncomingEmail: HeaderOnlyLevel,
		OutgoingEmail: HeaderOnlyLevel,
		Draft:         HeaderOnlyLevel,
		Chat:          HeaderOnlyLevel,
	})
	return svc.MailMonitor.Update(monitor)
}

func TestMailMonitorServiceUpdate(t *testing.T) {
	defer gock.Off()
	gock.New("https://apps-apis.google.com").
		Post("/a/feeds/compliance/audit/mail/monitor/example.com/abhishek").
		MatchType("application/atom\\+xml").
		MatchHeader("Authorization", "Bearer test").
		MatchHeader("User-Agent", "google-api-go-client/0.5").
		Reply(200).
		XML(`<entry xmlns='http://www.w3.org/2005/Atom' xmlns:apps='http://schemas.google.com/apps/2006'>
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
   </entry>`)

	monitor2, err := updateEmailMonitor()
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	if monitor2 == nil {
		t.Error("Expected not nil but got nil")
	}
}

func TestMailMonitorServiceUpdateHTTPError(t *testing.T) {
	defer gock.Off()
	gock.New("https://apps-apis.google.com").
		Post("/a/feeds/compliance/audit/mail/monitor/example.com/abhishek").
		MatchType("application/atom\\+xml").
		MatchHeader("Authorization", "Bearer test").
		MatchHeader("User-Agent", "google-api-go-client/0.5").
		ReplyError(errors.New("Error!"))

	monitor2, err := updateEmailMonitor()
	expected := "Post https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek: Error!"
	if err.Error() != expected {
		t.Errorf(`Expected "%v" but got "%v"`, expected, err)
	}
	if monitor2 != nil {
		t.Errorf("Expected nil but got %v", monitor2)
	}
}
