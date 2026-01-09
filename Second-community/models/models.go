package models

type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Password  string
	Questions []Question
	Answers   []Answer
}

func (u User) Error() string {
	panic("implement me")
}

type Question struct {
	ID      uint `gorm:"primaryKey"`
	Title   string
	Content string
	UserID  uint
	Answers []Answer
	// 新增点赞计数

}

type Answer struct {
	ID         uint `gorm:"primaryKey"`
	Content    string
	QuestionID uint
	UserID     uint
	// 新增点赞计数

}

type Follow struct {
	ID          uint `gorm:"primaryKey"`
	FollowerID  uint //关注者
	FollowingID uint //被关注者
}

type Answers1 struct {
	ID       uint
	Content  string
	AnswerID uint
	UserID   uint
}
