package approval

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
	Status      string `json:"status"`
	ResolvedBy  string `json:"resolvedBy"`
}