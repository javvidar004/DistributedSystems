package conectors

type CreateUserRequest struct {
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	LastName     string  `json:"last_name"`
	WorkPosition string  `json:"work_position"`
	Salary       float64 `json:"salary"`
}

type UpdateUserRequest struct {
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	LastName     string  `json:"last_name"`
	WorkPosition string  `json:"work_position"`
	Salary       float64 `json:"salary"`
}

type DeleteUserRequest struct {
	Id string `json:"id"`
}

type NewLogRequest struct {
	Timestamp string `json:"timestamp"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Status    string `json:"status"`
}
