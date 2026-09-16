package conectors

type NewLogRequest struct {
	Timestamp string `json:"timestamp"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Status    string `json:"status"`
}
