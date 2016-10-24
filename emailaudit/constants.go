package emailaudit

const (
	// FeedComplianceAuditScope FeedComplianceAuditScope OAuth2 Scope
	FeedComplianceAuditScope = "https://apps-apis.google.com/a/feeds/compliance/audit/"
)

// MonitorLevel MonitorLevel
type MonitorLevel string

const (
	// HeaderOnlyLevel HEADER_ONLY
	HeaderOnlyLevel MonitorLevel = "HEADER_ONLY"
	// FullMessageLevel FULL_MESSAGE
	FullMessageLevel MonitorLevel = "FULL_MESSAGE"
)
