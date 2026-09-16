package conectors

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name         string  `json:"name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	WorkPosition string  `json:"work_position"`
	Salary       float64 `json:"salary"`
	Password     string  `json:"password"`
}

type UpdatePasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}
