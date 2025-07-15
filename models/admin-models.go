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
