package repositories

import (
	"community/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	result := ur.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
func (ur *UserRepository) CreateQuestion(question *models.Question) (*models.Question, error) {
	result := ur.DB.Create(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return question, nil
}
func (ur *UserRepository) CreateAnswer(answer *models.Answer) (*models.Answer, error) {
	result := ur.DB.Create(&answer)
	if result.Error != nil {
		return nil, result.Error
	}
	return answer, nil
}
func (ur *UserRepository) Answer(answer *models.Answers1) (*models.Answers1, error) {
	result := ur.DB.Create(&answer)
	if result.Error != nil {
		return nil, result.Error

	}
	return answer, nil
}

func (ur *UserRepository) ModifyQuestion(question *models.Question) (*models.Question, error) {
	result := ur.DB.Model(&question).Updates(question)
	if result.Error != nil {
		return nil, result.Error
	}
	return question, nil
}

//	func (ur *UserRepository) ModifyQuestion(q *models.Question, realOwnerID uint) error {
//		// 强制：只有 ID 匹配 且 UserID 匹配当前 Token 用户时，才执行更新
//		result := ur.DB.Model(&models.Question{}).
//			Where("id = ? AND user_id = ?", q.ID, realOwnerID).
//			Select("Title", "Content"). // 只允许修改标题和内容，拒绝修改 UserID
//			Updates(q)
//
//		if result.RowsAffected == 0 {
//			return errors.New("权限不足或资源不存在")
//		}
//		return nil
//	}
func (ur *UserRepository) ModifyAnswer(answer *models.Answer) (*models.Answer, error) {
	result := ur.DB.Model(&answer).Updates(answer)
	if result.Error != nil {
		return nil, result.Error
	}
	return answer, nil
}
func (ur *UserRepository) DeleteQuestion(question *models.Question) (*models.Question, error) {
	result := ur.DB.Delete(&question)
	if result.Error != nil {
		return nil, result.Error
	}
	return question, nil
}
func (ur *UserRepository) DeleteAnswer(answer *models.Answer) (*models.Answer, error) {
	result := ur.DB.Delete(&answer)
	if result.Error != nil {
		return nil, result.Error
	}
	return answer, nil
}
func (ur *UserRepository) DeleteUser(user *models.User) (*models.User, error) {
	result := ur.DB.Delete(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil

}
