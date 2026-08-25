package conectors

type BookingRequest struct {
	UserID      int    `json:"user_id"`
	WorkspaceID int    `json:"workspace_id"`
	BookingDate string `json:"booking_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
