package entity

// GroupID はグループIDの型定義
// @Description グループを一意に識別するID
type GroupID string

// GroupName はグループ名の型定義
// @Description グループの名前
type GroupName string

// GroupPassword はグループパスワードの型定義
// @Description グループへの参加に必要なパスワード
type GroupPassword string

// GroupDescription はグループ説明の型定義
// @Description グループの説明文
type GroupDescription string

// GroupMembers はグループメンバーのリスト型定義
// @Description グループに所属するメンバーのIDリスト
type GroupMembers []UserID

// GroupEvents はグループイベントのリスト型定義
// @Description グループに関連するイベントのIDリスト
type GroupEvents []EventID

// Group はグループエンティティの構造体
// @Description グループエンティティ
type Group struct {
	GroupID          GroupID          `bson:"group_id" json:"group_id" example:"group123"`
	GroupName        GroupName        `bson:"name" json:"group_name" example:"テストグループ"`
	GroupPassword    GroupPassword    `bson:"password" json:"group_password" example:"password123"`
	GroupManagerID   UserID           `bson:"manager_id" json:"group_manager_id" example:"user456"`
	GroupDescription GroupDescription `bson:"description" json:"group_description" example:"これはテストグループです"`
	GroupMembers     GroupMembers     `bson:"members" json:"group_members" example:"[\"user123\",\"user456\"]"`
	GroupEvents      GroupEvents      `bson:"events" json:"group_events" example:"[\"event123\",\"event456\"]"`
}
