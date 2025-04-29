package entity

type UserID string
type AuthID string
type UserName string
type FCMToken string
type Alias string

type User struct {
	UserID   UserID   `json:"user_id"`
	AuthID   AuthID   `json:"auth_id"`
	UserName UserName `json:"user_name"`
	FCMToken FCMToken `json:"fcm_token"`
	Alias    Alias    `json:"alias"`
}
