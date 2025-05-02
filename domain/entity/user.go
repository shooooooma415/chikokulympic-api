package entity

type UserID string
type AuthID string
type UserName string
type UserIcon string
type FCMToken string
type Alias string

type User struct {
	UserID   UserID   `bson:"user_id" json:"user_id"`
	AuthID   AuthID   `bson:"auth_id" json:"auth_id"`
	UserName UserName `bson:"user_name" json:"user_name"`
	UserIcon UserIcon `bson:"user_icon" json:"user_icon"`
	FCMToken FCMToken `bson:"fcm_token" json:"fcm_token"`
	Alias    Alias    `bson:"alias" json:"alias"`
}
