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
	_, err = Db.Exec(Ctx, `UPDATE users SET otp = NULL, otp_expiry = NULL WHERE usertag = $1`, data.Usertag)
	if err != nil {
		log.Println("Failed to clear OTP:", err)
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
	_, err = Db.Exec(Ctx, `UPDATE users SET otp = NULL, otp_expiry = NULL WHERE email = $1`, data.Email)
	if err != nil {
		log.Println("Failed to clear OTP:", err)
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

func (AdminServer) GetDashboardSummary() (any, error) {
	var patientsCount, doctorsCount, appointmentsCount, ordersCount, doctorRequests int

	queries := []struct {
		query string
		dest  *int
	}{
		{"SELECT COUNT(*) FROM users WHERE role = 'user'", &patientsCount},
		{"SELECT COUNT(*) FROM users WHERE role = 'doctor'", &doctorsCount},
		{"SELECT COUNT(*) FROM appointments", &appointmentsCount},
		{"SELECT COUNT(*) FROM orders", &ordersCount},
		{"SELECT COUNT(*) FROM users WHERE role = 'doctor' AND status = 'pending'", &doctorRequests},
	}

	for _, q := range queries {
		if err := Db.QueryRow(Ctx, q.query).Scan(q.dest); err != nil {
			log.Printf("Dashboard query failed: %s — %v", q.query, err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
	}

	return map[string]interface{}{
		"patients_count":     patientsCount,
		"doctors_count":      doctorsCount,
		"appointments_count": appointmentsCount,
		"orders_count":       ordersCount,
		"doctor_requests":    doctorRequests,
	}, nil
}
func (AdminServer) GetAnalytics(data models.AnalyticsReq) (any, error) {

	if data.Metric != "payments" {
		return nil, errors.New("unsupported metric: " + data.Metric)
	} else {
		res, err := getPaymentAnalytics(data.Month, data.Year)
		if err != nil {
			log.Printf("Failed to get analytics for %s: %v", data.Metric, err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		return res, nil

	}

}

func getPaymentAnalytics(month, year string) (any, error) {

	var analytics models.AnalyticsResp

	query := `SELECT COALESCE(SUM(amount), 0), COUNT(*) FROM payments WHERE EXTRACT(MONTH FROM payment_date) = $1 AND 
			EXTRACT(YEAR FROM payment_date) = $2 AND status = 'completed'`

	err := Db.QueryRow(Ctx, query, month, year).Scan(&analytics.Total_amount, &analytics.Payment_count)
	if err != nil {
		log.Printf("Payment analytics query failed: %v", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	average := 0.0
	if analytics.Payment_count > 0 {
		average = analytics.Total_amount / float64(analytics.Payment_count)
	}

	analytics.Metric = "payments"
	analytics.Month = month
	analytics.Year = year
	analytics.Average_payment = average

	return analytics, nil
}
