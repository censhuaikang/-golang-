package controllers

import (
	"community/models"
	"community/services"
	"community/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// --- 登录注册 ---

func (uc *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "输入格式错误"})
		return
	}
	res, err := uc.UserService.Register(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该用户名已被注册"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (uc *UserController) Login(c *gin.Context) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	user, err := uc.UserService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败"})
		return
	}
	token, _ := utils.GenerateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// --- 删改操作 (核心逻辑：提取 userID 并传给 Service) ---

func (uc *UserController) DeleteQuestion(c *gin.Context) {
	operatorID := c.MustGet("userID").(uint)
	idStr := c.Query("id")
	qID, _ := strconv.Atoi(idStr)

	err := uc.UserService.DeleteQuestion(uint(qID), operatorID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "问题已删除"})
}

func (uc *UserController) ModifyQuestion(c *gin.Context) {
	operatorID := c.MustGet("userID").(uint)
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}
	question.UserID = operatorID
	// 执行修改
	res, err := uc.UserService.ModifyQuestion(question, operatorID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (uc *UserController) DeleteAnswer(c *gin.Context) {
	operatorID := c.MustGet("userID").(uint)
	idStr := c.Query("id")
	aID, _ := strconv.Atoi(idStr)

	err := uc.UserService.DeleteAnswer(uint(aID), operatorID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "回答已删除"})
}

func (uc *UserController) ModifyAnswer(c *gin.Context) {
	operatorID := c.MustGet("userID").(uint)
	var answer models.Answer
	if err := c.ShouldBindJSON(&answer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}
	res, err := uc.UserService.ModifyAnswer(answer, operatorID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (uc *UserController) Unfollow(c *gin.Context) {
	followerID := c.MustGet("userID").(uint)
	targetID, _ := strconv.Atoi(c.Param("id"))

	err := uc.UserService.Unfollow(followerID, uint(targetID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已取消关注"})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	err := uc.UserService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注销失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "账号已永久注销"})
}

// --- 基础增量操作 ---

func (uc *UserController) CreateQuestion(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var q models.Question
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}
	q.UserID = userID
	res, err := uc.UserService.CreateQuestion(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建问题失败"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (uc *UserController) CreateAnswer(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var a models.Answer
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}
	a.UserID = userID
	res, err := uc.UserService.CreateAnswer(a)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建回答失败"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (uc *UserController) CreateReply(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var replyReq struct {
		Content  string `json:"content" binding:"required"`
		AnswerID uint   `json:"answer_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&replyReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据格式错误"})
		return
	}

	// 创建回复对象
	reply := models.Reply{
		Content:  replyReq.Content,
		AnswerID: replyReq.AnswerID,
		UserID:   userID,
	}

	res, err := uc.UserService.CreateReply(reply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "回复失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (uc *UserController) Follow(c *gin.Context) {
	fID := c.MustGet("userID").(uint)
	tID, _ := strconv.Atoi(c.Param("id"))
	res, err := uc.UserService.Follow(fID, uint(tID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
