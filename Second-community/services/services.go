package services

import (
	"community/models"
	"community/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func hashPassword(password string) (string, error) {
	// GenerateFromPassword 自动处理盐值生成
	// DefaultCost 默认值为 10，平衡了安全性和性能
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 3. 新增验证函数：比较明文和 Hash
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
	// 第一步：先根据用户名查找用户 (不能直接在 SQL 里对比 bcrypt hash)
	err := us.UserRepo.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 第二步：验证密码
	if !checkPasswordHash(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

func (us *UserService) DeleteUser(id uint) error {
	// 只能注销自己（由 Controller 传入当前 Token 的 ID）
	return us.UserRepo.DB.Delete(&models.User{}, id).Error
}

// --- 问题相关 ---

func (us *UserService) CreateQuestion(q models.Question) (*models.Question, error) {
	return us.UserRepo.CreateQuestion(&q)
}

func (us *UserService) ModifyQuestion(q models.Question, operatorID uint) (models.Question, error) {
	var original models.Question
	// 第一步：先去数据库看一眼，这个问题的原主人是谁
	if err := us.UserRepo.DB.First(&original, q.ID).Error; err != nil {
		return models.Question{}, errors.New("问题不存在")
	}

	// 第二步：强行核对！新账号的 ID (operatorID) 必须等于原作者 ID
	if original.UserID != operatorID {
		// 哪怕你的 Token 是真的，只要你不是作者，就驳回！
		return models.Question{}, errors.New("【越权警告】你不是该问题的作者，无法修改")
	}

	// 第三步：只有核对通过，才允许调用 Repo 执行更新
	res, err := us.UserRepo.ModifyQuestion(&q)
	return *res, err
}

func (us *UserService) DeleteQuestion(qID uint, operatorID uint) error {
	var original models.Question
	if err := us.UserRepo.DB.First(&original, qID).Error; err != nil {
		return errors.New("问题不存在")
	}
	if original.UserID != operatorID {
		return errors.New("权限不足：你无权删除此内容")
	}
	_, err := us.UserRepo.DeleteQuestion(&models.Question{ID: qID})
	return err
}

// --- 回答相关 ---

func (us *UserService) CreateAnswer(a models.Answer) (*models.Answer, error) {
	return us.UserRepo.CreateAnswer(&a)
}

func (us *UserService) ModifyAnswer(a models.Answer, operatorID uint) (models.Answer, error) {
	var original models.Answer
	if err := us.UserRepo.DB.First(&original, a.ID).Error; err != nil {
		return models.Answer{}, errors.New("回答不存在")
	}
	if original.UserID != operatorID {
		return models.Answer{}, errors.New("权限不足：你不是该回答的作者")
	}
	res, err := us.UserRepo.ModifyAnswer(&a)
	return *res, err
}

func (us *UserService) DeleteAnswer(aID uint, operatorID uint) error {
	var original models.Answer
	if err := us.UserRepo.DB.First(&original, aID).Error; err != nil {
		return errors.New("回答不存在")
	}
	if original.UserID != operatorID {
		return errors.New("权限不足：无法删除他人内容")
	}
	_, err := us.UserRepo.DeleteAnswer(&models.Answer{ID: aID})
	return err
}

// --- 社交与评论 ---

func (us *UserService) Answer(a1 models.Answers1) (models.Answers1, error) {
	res, err := us.UserRepo.Answer(&a1)
	return *res, err
}

func (us *UserService) Follow(f, t uint) (models.Follow, error) {
	if f == t {
		return models.Follow{}, errors.New("不能关注你自己")
	}
	follow := models.Follow{FollowerID: f, FollowingID: t}
	err := us.UserRepo.DB.Create(&follow).Error
	return follow, err
}

func (us *UserService) Unfollow(f, t uint) error {
	return us.UserRepo.DB.Where("follower_id = ? AND following_id = ?", f, t).Delete(&models.Follow{}).Error
}
