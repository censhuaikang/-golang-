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

func (ur *UserRepository) CreateReply(reply *models.Reply) (*models.Reply, error) {
	result := ur.DB.Create(&reply)
	if result.Error != nil {
		return nil, result.Error
	}

	// 预加载关联数据以便返回完整信息
	ur.DB.Preload("User").Preload("Answer").First(&reply, reply.ID)

	return reply, nil
}

func (ur *UserRepository) ModifyQuestion(question *models.Question) (*models.Question, error) {
	result := ur.DB.Model(&question).Updates(question)
	if result.Error != nil {
		return nil, result.Error
	}
	return question, nil
}

func (ur *UserRepository) ModifyAnswer(answer *models.Answer) (*models.Answer, error) {
	result := ur.DB.Model(&answer).Updates(answer)
	if result.Error != nil {
		return nil, result.Error
	}
	return answer, nil
}

func (ur *UserRepository) DeleteQuestion(question *models.Question) error {
	result := ur.DB.Delete(&question)
	return result.Error
}

func (ur *UserRepository) DeleteAnswer(answer *models.Answer) error {
	result := ur.DB.Delete(&answer)
	return result.Error
}

func (ur *UserRepository) DeleteUser(user *models.User) error {
	result := ur.DB.Delete(&user)
	return result.Error
}
