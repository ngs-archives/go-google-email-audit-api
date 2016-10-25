package emailaudit

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

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
func (svc *MailMonitorService) Update(monitor MailMonitor) (*MailMonitor, error) {
	url := monitor.URL()
	body := bytes.NewReader(monitor.toXML())
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("User-Agent", svc.s.userAgent())
	req.Header.Add("Content-Type", contentType)
	res, err := svc.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	return monitorFromXML(bytes)
}
