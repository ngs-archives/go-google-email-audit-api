package emailaudit

const (
	// FeedComplianceAuditScope FeedComplianceAuditScope OAuth2 Scope
	FeedComplianceAuditScope = "https://apps-apis.google.com/a/feeds/compliance/audit/"
)

// MailMonitorLevel MailMonitorLevel
type MailMonitorLevel string

const (
	// NoneLevel HEADER_ONLY
	NoneLevel MailMonitorLevel = ""
	// HeaderOnlyLevel HEADER_ONLY
	HeaderOnlyLevel MailMonitorLevel = "HEADER_ONLY"
	// FullMessageLevel FULL_MESSAGE
	FullMessageLevel MailMonitorLevel = "FULL_MESSAGE"
)
