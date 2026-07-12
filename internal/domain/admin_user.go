package domain

type AdminUser struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
