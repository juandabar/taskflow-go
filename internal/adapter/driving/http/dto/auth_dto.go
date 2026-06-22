package dto

import "time"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	if r.Name == "" {
		return errRequired("name")
	}
	if r.Email == "" {
		return errRequired("email")
	}
	if len(r.Password) < 8 {
		return errInvalid("password", "must be at least 8 characters")
	}
	return nil
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	if r.Email == "" {
		return errRequired("email")
	}
	if r.Password == "" {
		return errRequired("password")
	}
	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}
