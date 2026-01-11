package services

import (
	"community/models"
	"community/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo *repositories.UserRepository
	DB       *gorm.DB
}

func NewUserService(userRepo *repositories.UserRepository, db *gorm.DB) *UserService {
	return &UserService{UserRepo: userRepo, DB: db}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// --- 用户相关 ---

func (us *UserService) Register(user models.User) (*models.User, error) {
	hashedPwd, err := hashPassword(user.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}
	user.Password = hashedPwd
	return us.UserRepo.CreateUser(&user)
}

func (us *UserService) Login(username, password string) (*models.User, error) {
	var user models.User
	err := us.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

func (us *UserService) DeleteUser(id uint) error {
	return us.DB.Delete(&models.User{}, id).Error
}

// --- 问题相关 ---

func (us *UserService) CreateQuestion(q models.Question) (*models.Question, error) {
	return us.UserRepo.CreateQuestion(&q)
}

func (us *UserService) ModifyQuestion(q models.Question, operatorID uint) (*models.Question, error) {
	var original models.Question
	if err := us.DB.First(&original, q.ID).Error; err != nil {
		return nil, errors.New("问题不存在")
	}

	if original.UserID != operatorID {
		return nil, errors.New("【越权警告】你不是该问题的作者，无法修改")
	}

	return us.UserRepo.ModifyQuestion(&q)
}

func (us *UserService) DeleteQuestion(qID uint, operatorID uint) error {
	var original models.Question
	if err := us.DB.First(&original, qID).Error; err != nil {
		return errors.New("问题不存在")
	}
	if original.UserID != operatorID {
		return errors.New("权限不足：你无权删除此内容")
	}
	return us.UserRepo.DeleteQuestion(&models.Question{ID: qID})
}

// --- 回答相关 ---

func (us *UserService) CreateAnswer(a models.Answer) (*models.Answer, error) {
	return us.UserRepo.CreateAnswer(&a)
}

func (us *UserService) ModifyAnswer(a models.Answer, operatorID uint) (*models.Answer, error) {
	var original models.Answer
	if err := us.DB.First(&original, a.ID).Error; err != nil {
		return nil, errors.New("回答不存在")
	}
	if original.UserID != operatorID {
		return nil, errors.New("权限不足：你不是该回答的作者")
	}
	return us.UserRepo.ModifyAnswer(&a)
}

func (us *UserService) DeleteAnswer(aID uint, operatorID uint) error {
	var original models.Answer
	if err := us.DB.First(&original, aID).Error; err != nil {
		return errors.New("回答不存在")
	}
	if original.UserID != operatorID {
		return errors.New("权限不足：无法删除他人内容")
	}
	return us.UserRepo.DeleteAnswer(&models.Answer{ID: aID})
}

// --- 社交与评论 ---

func (us *UserService) CreateReply(reply models.Reply) (*models.Reply, error) {
	// 验证被回复的答案是否存在
	var parentAnswer models.Answer
	if err := us.DB.First(&parentAnswer, reply.AnswerID).Error; err != nil {
		return nil, errors.New("被回复的答案不存在")
	}

	return us.UserRepo.CreateReply(&reply)
}

func (us *UserService) Follow(f, t uint) (*models.Follow, error) {
	if f == t {
		return nil, errors.New("不能关注你自己")
	}
	follow := models.Follow{FollowerID: f, FollowingID: t}
	err := us.DB.Create(&follow).Error
	return &follow, err
}

func (us *UserService) Unfollow(f, t uint) error {
	return us.DB.Where("follower_id = ? AND following_id = ?", f, t).Delete(&models.Follow{}).Error
}
