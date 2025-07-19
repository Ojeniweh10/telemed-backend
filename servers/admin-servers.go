package servers

import (
	"context"
	"errors"
	"log"
	"telemed/models"
	"telemed/responses"
	"telemed/utils"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AdminServer struct{}

var Ctx context.Context
var Db *pgxpool.Pool

func (AdminServer) Login(data models.Adminlogin) (any, error) {
	var hash string
	var admin models.AdminLoginResponse
	err := Db.QueryRow(Ctx, "SELECT password, usertag FROM users WHERE email = $1 AND role = 'admin'", data.Email).Scan(&hash, &admin.Usertag)
	if err != nil {
		log.Println(err)
		return nil, errors.New(responses.ACCOUNT_NON_EXISTENT)
	}

	pwdCheck := utils.VerifyPassword(data.Password, hash)
	if !pwdCheck {
		log.Println("Invalid password for admin login")
		return nil, errors.New(responses.INVALID_PASSWORD)
	}
	otp, err := utils.GenerateOTP()
	if err != nil {
		log.Println("Failed to generate OTP:", err)
		return nil, errors.New("failed to generate OTP")
	}
	_, err = Db.Exec(Ctx, "UPDATE users SET otp = $1, otp_expiry = NOW()+ INTERVAL '5 minutes' WHERE email = $2", otp, data.Email)
	if err != nil {
		log.Println("failed to save OTP", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	err = utils.SendEmailOTP(data.Email, otp)
	if err != nil {
		log.Println("Failed to send OTP email:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return admin, nil
}

func (AdminServer) VerifyOTP(data models.OTPVerify) (any, error) {
	var dbOtp string
	var role string
	var otpExpiryTime time.Time
	err := Db.QueryRow(Ctx, "SELECT otp, otp_expiry, role FROM users WHERE usertag = $1", data.Usertag).Scan(&dbOtp, &otpExpiryTime, &role)
	if err != nil {
		log.Println(err)
		return nil, errors.New("invalid email or OTP")
	}

	if data.OTP != dbOtp {
		log.Println("Invalid OTP for admin login")
		return nil, errors.New("invalid OTP")
	}

	if time.Now().After(otpExpiryTime) {
		log.Println("OTP has expired")
		return nil, errors.New("OTP has expired")
	}

	token, err := utils.GenerateJWT(data.Usertag, role)
	if err != nil {
		log.Println("Failed to generate JWT token:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	}, nil
}

func (AdminServer) ForgotPassword(data models.ForgotPassword) (any, error) {
	var exists string
	err := Db.QueryRow(Ctx, "SELECT email FROM users WHERE email = $1 AND role = 'admin'", data.Email).Scan(&exists)
	if err != nil {
		log.Println(err)
		return nil, errors.New(responses.ACCOUNT_NON_EXISTENT)
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		log.Println("Failed to generate OTP:", err)
		return nil, errors.New("failed to generate OTP")
	}

	_, err = Db.Exec(Ctx, "UPDATE users SET otp = $1, otp_expiry = NOW() + INTERVAL '10 minutes' WHERE email = $2", otp, data.Email)
	if err != nil {
		log.Println("failed to save OTP", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	err = utils.SendEmailOTP(data.Email, otp)
	if err != nil {
		log.Println("Failed to send OTP email:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return nil, err
}

func (AdminServer) VerifyPwdOTP(data models.VerifyPwdOTP) (any, error) {
	var dbOtp string
	var otpExpiryTime time.Time
	var role string

	err := Db.QueryRow(Ctx, "SELECT otp, otp_expiry, role FROM users WHERE email = $1", data.Email).
		Scan(&dbOtp, &otpExpiryTime, &role)
	if err != nil {
		log.Println(err)
		return nil, errors.New("invalid email or OTP")
	}

	if role != "admin" {
		log.Println("Not an admin account")
		return nil, errors.New("unauthorized")
	}

	if data.OTP != dbOtp {
		log.Println("Invalid OTP for admin")
		return nil, errors.New("invalid OTP")
	}

	if time.Now().After(otpExpiryTime) {
		log.Println("OTP has expired")
		return nil, errors.New("OTP has expired")
	}

	return map[string]interface{}{
		"message": "OTP verified successfully",
	}, nil
}

func (AdminServer) ResetPassword(data models.ResetPassword) (any, error) {
	var exists string
	err := Db.QueryRow(Ctx, "SELECT email FROM users WHERE email = $1 AND role = 'admin'", data.Email).Scan(&exists)
	if err != nil {
		log.Println(err)
		return nil, errors.New(responses.ACCOUNT_NON_EXISTENT)
	}

	if data.NewPassword == "" {
		return nil, errors.New(responses.INCOMPLETE_DATA)
	}

	hashedPwd, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	_, err = Db.Exec(Ctx, "UPDATE users SET password = $1 WHERE email = $2", hashedPwd, data.Email)
	if err != nil {
		log.Println("Failed to reset password:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return map[string]interface{}{
		"message": responses.PASSWORD_RESET_SUCCESS,
	}, nil
}
