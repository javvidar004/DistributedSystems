package conectors

type UserInfoResponse struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UpdatePasswordRequest struct {
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}
