package entity


type UserID string

type AuthID string

type UserName string

type UserIcon string

type FCMToken string

type Alias string

type User struct {
	UserID   UserID   `bson:"user_id" json:"user_id" example:"user123"`
	AuthID   AuthID   `bson:"auth_id" json:"auth_id" example:"auth456"`
	UserName UserName `bson:"user_name" json:"user_name" example:"山田太郎"`
	UserIcon UserIcon `bson:"user_icon" json:"user_icon" example:"https://example.com/icon.png"`
	FCMToken FCMToken `bson:"fcm_token" json:"fcm_token" example:"fcm-token-123456"`
	Alias    Alias    `bson:"alias" json:"alias" example:"たろう"`
}
