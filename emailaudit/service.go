package emailaudit

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/api/googleapi"
)

const contentType = "application/atom+xml"

// Service Service
type Service struct {
	client      *http.Client
	MailMonitor *MailMonitorService
	UserAgent   string
}

// MailMonitorService MailMonitorService
type MailMonitorService struct {
	s *Service
}

// New returns new Service
func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client}
	s.MailMonitor = NewMailMonitorService(s)
	return s, nil
}

func (s *Service) userAgent() string {
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}

// NewMailMonitorService returns new MailMonitorService
func NewMailMonitorService(s *Service) *MailMonitorService {
	rs := &MailMonitorService{s: s}
	return rs
}

// Update creates or updates EmailMonitor
// - https://developers.google.com/admin-sdk/email-audit/#creating_a_new_email_monitor
// - https://developers.google.com/admin-sdk/email-audit/#updating_an_email_monitor
func (svc *MailMonitorService) Update(domainName string, sourceUserName string, destUserName string, endDate time.Time, monitorLevels MailMonitorLevels) (*MailMonitor, error) {
	monitor := NewMailMonitor(domainName, sourceUserName, destUserName, endDate, monitorLevels)
	url := monitor.URL()
	body := bytes.NewReader(monitor.toXML())
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("User-Agent", svc.s.userAgent())
	req.Header.Add("Content-Type", contentType)
	res, err := svc.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, _ := ioutil.ReadAll(res.Body)
	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return nil, errors.New(string(bytes))
	}
	return monitorFromXML(bytes)
}

// List Retrieving all email monitors of a source user
// - https://developers.google.com/admin-sdk/email-audit/#retrieving_all_email_monitors_of_a_source_user
func (svc *MailMonitorService) List(domain string, sourceUserName string) ([]MailMonitor, error) {
	url := fmt.Sprintf("%v/%v/%v", baseURL, domain, sourceUserName)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", svc.s.userAgent())
	res, err := svc.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, _ := ioutil.ReadAll(res.Body)
	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return nil, errors.New(string(bytes))
	}
	return monitorsFromXML(bytes)
}

// Disable Deleting an email monitor
// - https://developers.google.com/admin-sdk/email-audit/#deleting_an_email_monitor
func (svc *MailMonitorService) Disable(domain string, sourceUserName string, destUserName string) error {
	url := fmt.Sprintf("%v/%v/%v/%v", baseURL, domain, sourceUserName, destUserName)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("User-Agent", svc.s.userAgent())
	res, err := svc.s.client.Do(req)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return errors.New(string(bytes))
	}
	return err
}
