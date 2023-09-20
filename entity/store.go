package entity

// Store represents data about an store.
type Store struct {
	User      User   `json:"user"`
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Block     string `json:"block"`
	City      string `json:"city"`
	State     string `json:"state"`
	PhotoPath string `json:"photo_path"`
	UserID    int    `json:"user_id"`
}
