package request

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	PasswordAgain string `json:"password_again" binding:"required,eqfield=Password"`
	Email         string `json:"email" binding:"required,email"`
	CaptchaID     string `json:"captcha_id" binding:"required"`
	CaptchaCode   string `json:"captcha_code" binding:"required"`
}
