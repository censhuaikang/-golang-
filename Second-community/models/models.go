package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Password  string
	Questions []Question `gorm:"foreignKey:UserID"`
	Answers   []Answer   `gorm:"foreignKey:UserID"`
	Replies   []Reply    `gorm:"foreignKey:UserID"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

func (u User) Error() string {
	panic("implement me")
}

type Question struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string
	UserID    uint
	Answers   []Answer  `gorm:"foreignKey:QuestionID"`
	Likes     uint      `gorm:"default:0"` // 新增点赞计数
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Answer struct {
	ID         uint `gorm:"primaryKey"`
	Content    string
	QuestionID uint `json:"question_id"`
	UserID     uint
	Likes      uint      `gorm:"default:0"` // 新增点赞计数
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// 关联关系
	Replies  []Reply  `gorm:"foreignKey:AnswerID;constraint:OnDelete:CASCADE"`
	Question Question `gorm:"foreignKey:QuestionID;references:ID"`
	User     User     `gorm:"foreignKey:UserID;references:ID"`
}

type Follow struct {
	ID          uint `gorm:"primaryKey"`
	FollowerID  uint //关注者
	FollowingID uint //被关注者
}

type Reply struct {
	ID        uint `gorm:"primaryKey"`
	Content   string
	ParentID  uint // 被回复的 AnswerID
	UserID    uint
	AnswerID  uint      `json:"answer-id"` // 关联到原始回答`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// 关联关系
	Answer Answer `gorm:"foreignKey:AnswerID;references:ID"`
	User   User   `gorm:"foreignKey:UserID;references:ID"`
}
