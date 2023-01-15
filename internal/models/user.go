package models

type RegisterRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email,max=25"`
	FullName string `json:"full_name,omitempty" validate:"required,max=40"`
	Password string `json:"password,omitempty" validate:"required,max=25,min=4"`
}

type RegisterResponse struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

type LoginResponse struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type LogoutResponse struct {
	Message string `json:"message,omitempty"`
}

type User struct {
	ID        int64  `json:"id,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
