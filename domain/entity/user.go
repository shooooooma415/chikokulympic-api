package entity

// UserID はユーザーIDの型定義
// @Description ユーザーを一意に識別するID
type UserID string

// AuthID は認証IDの型定義
// @Description 外部認証システムで使用されるID
type AuthID string

// UserName はユーザー名の型定義
// @Description ユーザーの表示名
type UserName string

// UserIcon はユーザーアイコンの型定義
// @Description ユーザーのアイコンURL
type UserIcon string

// FCMToken はFCMトークンの型定義
// @Description Firebase Cloud Messagingで使用されるトークン
type FCMToken string

// Alias はユーザーのエイリアスの型定義
// @Description ユーザーの別名
type Alias string

// User はユーザーエンティティの構造体
// @Description ユーザーエンティティ
type User struct {
	UserID   UserID   `bson:"user_id" json:"user_id" example:"user123"`
	AuthID   AuthID   `bson:"auth_id" json:"auth_id" example:"auth456"`
	UserName UserName `bson:"user_name" json:"user_name" example:"山田太郎"`
	UserIcon UserIcon `bson:"user_icon" json:"user_icon" example:"https://example.com/icon.png"`
	FCMToken FCMToken `bson:"fcm_token" json:"fcm_token" example:"fcm-token-123456"`
	Alias    Alias    `bson:"alias" json:"alias" example:"たろう"`
}
