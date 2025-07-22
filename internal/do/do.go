package do

import "time"

type AccountResDO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      []string  `json:"role"`
	Status    string    `json:"status"`
	RoleNames []string  `json:"role_names"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      int       `json:"type"`
}
