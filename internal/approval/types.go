package approval

type Status string

const (
	StatusPending  Status = "PENDING"
	StatusApproved Status = "APPROVED"
	StatusRejected Status = "REJECTED"
	StatusTimeout  Status = "TIMEOUT"
)

func (s Status) IsTerminal() bool {
	return s == StatusApproved || s == StatusRejected || s == StatusTimeout
}

type CreateRequest struct {
	RequestedBy string `json:"requestedBy"`
	HostName    string `json:"hostName"`
	Command     string `json:"command"`
}

type ActionRequest struct {
	ResolvedBy string `json:"resolvedBy"`
}

type Response struct {
	ID          int64  `json:"id"`
	RequestedBy string `json:"requestedBy"`
	HostName    string `json:"hostName"`
	Command     string `json:"command"`
	RiskLevel   string `json:"riskLevel"`
	Status      Status `json:"status"`
	ResolvedBy  string `json:"resolvedBy"`
}
