package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/LucienVen/golang-auth-service/internal/errors"
	"github.com/LucienVen/golang-auth-service/internal/middleware"
	"github.com/LucienVen/golang-auth-service/internal/request"
	"github.com/LucienVen/golang-auth-service/internal/response"
	"github.com/LucienVen/golang-auth-service/internal/service"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// GetAuthService 获取认证服务
func (ctrl *AuthController) GetAuthService() service.AuthService {
	return ctrl.authService
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户使用账号密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param loginRequest body request.LoginRequest true "登录请求"
// @Success 200 {object} response.LoginResponse "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "认证失败"
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数格式错误: "+err.Error())
		return
	}

	// 设置客户端信息
	req.DeviceID = c.GetHeader("X-Device-ID")
	req.UserAgent = c.GetHeader("User-Agent")
	req.IPAddress = c.ClientIP()

	// 调用认证服务
	loginResp, err := ctrl.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, errors.ErrCodeUnauthorized, err.Error())
		return
	}

	response.Success(c, loginResp)
}

// Register 用户注册
// @Summary 用户注册
// @Description 新用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param registerRequest body request.RegisterRequest true "注册请求"
// @Success 200 {object} response.RegisterResponse "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 409 {object} response.Response "用户已存在"
// @Router /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数格式错误: "+err.Error())
		return
	}

	// 设置客户端信息
	req.DeviceID = c.GetHeader("X-Device-ID")
	req.UserAgent = c.GetHeader("User-Agent")
	req.IPAddress = c.ClientIP()

	// 调用认证服务
	registerResp, err := ctrl.authService.Register(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, errors.ErrCodeParamInvalid, err.Error())
		return
	}

	response.Success(c, registerResp)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出，使当前令牌失效
// @Tags 认证
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer令牌"
// @Success 200 {object} response.Response "登出成功"
// @Failure 401 {object} response.Response "未认证"
// @Router /auth/logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	// 从Authorization头获取令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Unauthorized(c)
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		response.Unauthorized(c)
		return
	}

	tokenString := authHeader[len(bearerPrefix):]

	// 调用认证服务登出
	err := ctrl.authService.Logout(c.Request.Context(), tokenString)
	if err != nil {
		response.Error(c, errors.ErrCodeUnauthorized, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "登出成功"})
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param refreshRequest body request.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} response.RefreshTokenResponse "刷新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "刷新令牌无效"
// @Router /auth/refresh [post]
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数格式错误: "+err.Error())
		return
	}

	// 调用认证服务刷新令牌
	refreshResp, err := ctrl.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, errors.ErrCodeRefreshTokenInvalid, err.Error())
		return
	}

	response.Success(c, refreshResp)
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前用户的详细资料
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer令牌"
// @Success 200 {object} response.UserProfileResponse "用户资料"
// @Failure 401 {object} response.Response "未认证"
// @Router /user/profile [get]
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 调用认证服务获取用户资料
	profileResp, err := ctrl.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, errors.ErrCodeUserNotFound, err.Error())
		return
	}

	response.Success(c, profileResp)
}

// GetSessions 获取用户会话列表
// @Summary 获取用户会话列表
// @Description 获取当前用户的所有活跃会话
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer令牌"
// @Success 200 {object} response.SessionsResponse "会话列表"
// @Failure 401 {object} response.Response "未认证"
// @Router /user/sessions [get]
func (ctrl *AuthController) GetSessions(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 调用认证服务获取会话列表
	sessions, err := ctrl.authService.GetSessions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, errors.ErrCodeSystemError, err.Error())
		return
	}

	sessionsResp := &response.SessionsResponse{
		Sessions: sessions,
		Count:    len(sessions),
	}

	response.Success(c, sessionsResp)
}

// RevokeSession 撤销指定会话
// @Summary 撤销指定会话
// @Description 撤销用户指定的会话
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer令牌"
// @Param sessionId path string true "会话ID"
// @Success 200 {object} response.Response "撤销成功"
// @Failure 401 {object} response.Response "未认证"
// @Failure 400 {object} response.Response "参数错误"
// @Router /user/sessions/{sessionId} [delete]
func (ctrl *AuthController) RevokeSession(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 获取会话ID
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		response.BadRequest(c, "会话ID不能为空")
		return
	}

	// 调用认证服务撤销会话
	err := ctrl.authService.RevokeSession(c.Request.Context(), userID, sessionID)
	if err != nil {
		response.Error(c, errors.ErrCodeSystemError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "会话已撤销"})
}

// RevokeAllSessions 撤销所有会话
// @Summary 撤销所有会话
// @Description 撤销用户的所有会话（除当前会话外）
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer令牌"
// @Success 200 {object} response.Response "撤销成功"
// @Failure 401 {object} response.Response "未认证"
// @Router /user/sessions [delete]
func (ctrl *AuthController) RevokeAllSessions(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		response.Unauthorized(c)
		return
	}

	// 调用认证服务撤销所有会话
	err := ctrl.authService.RevokeAllSessions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, errors.ErrCodeSystemError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "所有会话已撤销"})
}

// ValidateToken 验证令牌
// @Summary 验证令牌
// @Description 验证令牌的有效性并返回用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer令牌"
// @Success 200 {object} response.UserResponse "令牌有效"
// @Failure 401 {object} response.Response "令牌无效"
// @Router /auth/validate [get]
func (ctrl *AuthController) ValidateToken(c *gin.Context) {
	// 从Authorization头获取令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Unauthorized(c)
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		response.Unauthorized(c)
		return
	}

	tokenString := authHeader[len(bearerPrefix):]

	// 调用认证服务验证令牌
	userResp, err := ctrl.authService.ValidateToken(c.Request.Context(), tokenString)
	if err != nil {
		response.Error(c, errors.ErrCodeTokenInvalid, err.Error())
		return
	}

	response.Success(c, userResp)
}