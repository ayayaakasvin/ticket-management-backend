package request

type UserRequest struct {
	Username 	string `json:"username,omitempty"`
	Password 	string `json:"password"`
	Email		string `json:"email,omitempty"`
}