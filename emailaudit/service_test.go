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
	return svc.MailMonitor.Update("example.com", "abhishek", "namrata", endDate, MailMonitorLevels{
		IncomingEmail: HeaderOnlyLevel,
		OutgoingEmail: HeaderOnlyLevel,
		Draft:         HeaderOnlyLevel,
		Chat:          HeaderOnlyLevel,
	})
}

func listEmailMonitors() ([]MailMonitor, error) {
	ctx := context.Background()
	config := &oauth2.Config{}
	token := &oauth2.Token{AccessToken: "test"}
	client := config.Client(ctx, token)
	svc, _ := New(client)
	return svc.MailMonitor.List("example.com", "abhishek")
}

func TestMailMonitorServiceUpdate(t *testing.T) {
	defer gock.Off()
	gock.New("https://apps-apis.google.com").
		Post("/a/feeds/compliance/audit/mail/monitor/example.com/abhishek").
		MatchType("application/atom\\+xml").
		MatchHeader("Authorization", "Bearer test").
		MatchHeader("User-Agent", "google-api-go-client/0.5").
		Reply(200).
		XML(monitorXML)

	monitor2, err := updateEmailMonitor()
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	_TestMonitor(monitor2, t)
}

func TestMailMonitorServiceList(t *testing.T) {
	defer gock.Off()
	gock.New("https://apps-apis.google.com").
		Get("/a/feeds/compliance/audit/mail/monitor/example.com/abhishek").
		MatchHeader("Authorization", "Bearer test").
		MatchHeader("User-Agent", "google-api-go-client/0.5").
		Reply(200).
		XML(monitorsXML)

	m, err := listEmailMonitors()
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}
	_TestMonitors(m, t)
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

func TestMailMonitorServiceListHTTPError(t *testing.T) {
	defer gock.Off()
	gock.New("https://apps-apis.google.com").
		Get("/a/feeds/compliance/audit/mail/monitor/example.com/abhishek").
		MatchHeader("Authorization", "Bearer test").
		MatchHeader("User-Agent", "google-api-go-client/0.5").
		ReplyError(errors.New("Error!"))

	monitor2, err := listEmailMonitors()
	expected := "Get https://apps-apis.google.com/a/feeds/compliance/audit/mail/monitor/example.com/abhishek: Error!"
	if err.Error() != expected {
		t.Errorf(`Expected "%v" but got "%v"`, expected, err)
	}
	if monitor2 != nil {
		t.Errorf("Expected nil but got %v", monitor2)
	}
}
