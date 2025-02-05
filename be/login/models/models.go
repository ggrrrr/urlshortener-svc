package models

type (
	UserPasswordRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	UserPasswordResponse struct {
		Token string `json:"token"`
	}
)
