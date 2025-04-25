package user

type UserID string
type UserName string
type FCMToken string
type Alias string

type User struct {
	UserID   UserID   `json:"user_id"`
	UserName UserName `json:"user_name"`
	FCMToken FCMToken `json:"fcm_token"`
	Alias    Alias    `json:"alias"`
}
