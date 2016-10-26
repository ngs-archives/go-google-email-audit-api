# go-google-email-audit-api

[![Build Status](https://travis-ci.org/ngs/go-google-email-audit-api.svg?branch=master)](https://travis-ci.org/ngs/go-google-email-audit-api)
[![GoDoc](https://godoc.org/github.com/ngs/go-google-email-audit-api/emailaudit?status.svg)](https://godoc.org/github.com/ngs/go-google-email-audit-api/emailaudit)
[![Go Report Card](https://goreportcard.com/badge/github.com/ngs/go-google-email-audit-api)](https://goreportcard.com/report/github.com/ngs/go-google-email-audit-api)
[![Coverage Status](https://coveralls.io/repos/github/ngs/go-google-email-audit-api/badge.svg?branch=master)](https://coveralls.io/github/ngs/go-google-email-audit-api?branch=master)

Go Client Library for [Google Email Audit API]

```sh
go get -u github.com/ngs/go-google-email-audit-api/emailaudit
```

## Email Monitor API

```go
import (
	// ...
	"github.com/ngs/go-google-email-audit-api/emailaudit"
)

func main() {
	// ...
	srv, err := emailaudit.New(client) // client = http.Client
	if err != nil {
		log.Fatalf("Unable to retrieve Email Audit API Client %v", err)
	}
	endDate := time.Date(2116, time.October, 31, 23, 59, 59, 0, time.UTC)

	// Create or update Email Monitor
	monitor, err := srv.MailMonitor.Update("example.com",
		"ngs", "kyohei", endDate,
		emailaudit.MailMonitorLevels{
			IncomingEmail: emailaudit.FullMessageLevel,
			OutgoingEmail: emailaudit.FullMessageLevel,
			Draft:         emailaudit.FullMessageLevel,
			Chat:          emailaudit.FullMessageLevel,
		},
	)
	if err != nil {
		log.Fatalf("Unable to update email monitor. %v", err)
	}

	// List Email Monitors
	monitors, err := srv.MailMonitor.List("example.com", "ngs")
	if err != nil {
		log.Fatalf("Unable to list email monitor. %v", err)
	}
	for _, m := range monitors {
		fmt.Printf("%v %v@%v chat:%v draft:%v incoming:%v outgoing:%v\n",
			m.Updated, m.DestUserName, m.DomainName,
			m.MonitorLevels.Chat, m.MonitorLevels.Draft,
			m.MonitorLevels.IncomingEmail, m.MonitorLevels.OutgoingEmail)
	}

	// Disable Email Monitor
	err = srv.MailMonitor.Disable("example.com", "ngs", "kyohei")
	if err != nil {
		log.Fatalf("Unable to disable email monitor. %v", err)
	}
}
```

## Mailbox Download

Not yet implemented

## Author

[Atsushi Nagase]

## License

See [LICENSE]

[Google Email Audit API]: https://developers.google.com/admin-sdk/email-audit/
[Atsushi Nagase]: https://ngs.io
[LICENSE]: LICENSE
