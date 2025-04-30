package entity

type UserID string
type AuthID string
type UserName string
type FCMToken string
type Alias string

type User struct {
	UserID   UserID   `bson:"user_id" json:"user_id"`
	AuthID   AuthID   `bson:"auth_id" json:"auth_id"`
	UserName UserName `bson:"user_name" json:"user_name"`
	FCMToken FCMToken `bson:"fcm_token" json:"fcm_token"`
	Alias    Alias    `bson:"alias" json:"alias"`
}
