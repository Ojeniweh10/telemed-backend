package models

type Adminlogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	Usertag string `json:"usertag"`
}

type OTPVerify struct {
	OTP     string `json:"otp"`
	Usertag string `json:"usertag"`
}

type ForgotPassword struct {
	Email string `json:"email"`
}

type VerifyPwdOTP struct {
	OTP   string `json:"otp"`
	Email string `json:"email"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}
