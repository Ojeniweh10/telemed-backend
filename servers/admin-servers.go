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
	err := Db.QueryRow(Ctx, "SELECT password, admintag FROM admins WHERE email = $1", data.Email).Scan(&hash, &admin.Usertag)
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
	_, err = Db.Exec(Ctx, "UPDATE admins SET otp = $1, otp_expiry = NOW()+ INTERVAL '5 minutes' WHERE email = $2", otp, data.Email)
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
	var otpExpiryTime time.Time
	err := Db.QueryRow(Ctx, "SELECT otp, otp_expiry, role FROM admins WHERE admintag = $1", data.Usertag).Scan(&dbOtp, &otpExpiryTime)
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
	_, err = Db.Exec(Ctx, `UPDATE admins SET otp = NULL, otp_expiry = NULL WHERE admintag = $1`, data.Usertag)
	if err != nil {
		log.Println("Failed to clear OTP:", err)
	}

	token, err := utils.GenerateJWT(data.Usertag)
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
	err := Db.QueryRow(Ctx, "SELECT email FROM admins WHERE email = $1", data.Email).Scan(&exists)
	if err != nil {
		log.Println(err)
		return nil, errors.New(responses.ACCOUNT_NON_EXISTENT)
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		log.Println("Failed to generate OTP:", err)
		return nil, errors.New("failed to generate OTP")
	}

	_, err = Db.Exec(Ctx, "UPDATE admins SET otp = $1, otp_expiry = NOW() + INTERVAL '10 minutes' WHERE email = $2", otp, data.Email)
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

	err := Db.QueryRow(Ctx, "SELECT otp, otp_expiry FROM admins WHERE email = $1", data.Email).
		Scan(&dbOtp, &otpExpiryTime)
	if err != nil {
		log.Println(err)
		return nil, errors.New("invalid email or OTP")
	}

	if data.OTP != dbOtp {
		log.Println("Invalid OTP for admin")
		return nil, errors.New("invalid OTP")
	}

	if time.Now().After(otpExpiryTime) {
		log.Println("OTP has expired")
		return nil, errors.New("OTP has expired")
	}
	_, err = Db.Exec(Ctx, `UPDATE admins SET otp = NULL, otp_expiry = NULL WHERE email = $1`, data.Email)
	if err != nil {
		log.Println("Failed to clear OTP:", err)
	}

	return map[string]interface{}{
		"message": "OTP verified successfully",
	}, nil
}

func (AdminServer) ResetPassword(data models.ResetPassword) (any, error) {
	var exists string
	err := Db.QueryRow(Ctx, "SELECT email FROM admins WHERE email = $1", data.Email).Scan(&exists)
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

	_, err = Db.Exec(Ctx, "UPDATE admins SET password = $1 WHERE email = $2", hashedPwd, data.Email)
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

func (AdminServer) GetAppointments() (any, error) {
	//rememmebr to modify to fetch using filters
	var appointments []models.Appointment

	rows, err := Db.Query(Ctx, "SELECT appointment_id, patient_tag, doctor_tag, scheduled_at, reason, file_url FROM appointments")
	if err != nil {
		log.Println("Failed to fetch appointments:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer rows.Close()

	for rows.Next() {
		var appointment models.Appointment
		if err := rows.Scan(&appointment.ID, &appointment.UserTag, &appointment.DoctorTag, &appointment.Scheduled_at, &appointment.Reason, &appointment.Fileurl, &appointment.Status, &appointment.Created_at); err != nil {
			log.Println("Failed to scan appointment:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		appointments = append(appointments, appointment)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over appointments:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return appointments, nil
}

func (AdminServer) GetAppointmentByID(payload models.AppointmentID) (any, error) {
	var data models.AppointmentIDResp

	query := `
		SELECT a.patient_tag, a.doctor_tag, a.scheduled_at, a.reason, a.file_url, a.status, a.created_at,
		       u.firstname, u.lastname, u.phone_no, u.gender, u.date_of_birth,
		       d.fullname, d.price_per_session
		FROM appointments a
		JOIN users u ON a.patient_tag = u.usertag
		JOIN doctors d ON a.doctor_tag = d.doctortag
		WHERE a.appointment_id = $1
	`

	err := Db.QueryRow(Ctx, query, payload.ID).Scan(
		&data.UserTag,
		&data.DoctorTag,
		&data.Scheduled_At,
		&data.Reason,
		&data.File_URL,
		&data.Status,
		&data.Created_At,
		&data.First_name,
		&data.Last_name,
		&data.Phone_No,
		&data.Gender,
		&data.Dob,
		&data.Doctor_Fullname,
		&data.Price,
	)

	if err != nil {
		log.Println("Failed to fetch full appointment details:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return data, nil
}

func (AdminServer) GetDoctorByID(data models.Doctorreq) (any, error) {
	var doctor models.Doctor

	query := `
	SELECT 
		d.doctortag,
		d.fullname,
		d.date_of_birth,
		d.phone_number,
		d.gender,
		d.specialization,
		d.country,
		d.city,
		d.yrs_of_experience,
		d.price_per_session,
		d.about,
		d.availability,
		d.profile_pic_url,
		h.name AS hospital_affiliation
	FROM doctors d
	LEFT JOIN hospitals h ON d.hospital_id = h.hospital_id
	WHERE d.fullname = $1
	`

	err := Db.QueryRow(Ctx, query, data).Scan(
		&doctor.DoctorTag,
		&doctor.FullName,
		&doctor.Dob,
		&doctor.Phone_no,
		&doctor.Gender,
		&doctor.Specialization,
		&doctor.Country,
		&doctor.City,
		&doctor.YearsOfExperience,
		&doctor.Price,
		&doctor.About,
		&doctor.Availability,
		&doctor.ProfilePicURL,
		&doctor.HospitalAffiliation, // ← from hospitals.name
	)

	if err != nil {
		log.Println("Error fetching doctor by fullname:", err)
		if err.Error() == "no rows in result set" {
			return nil, errors.New("doctor not found")
		}
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return doctor, nil
}

func (AdminServer) UpdateAppointmentStatus(payload models.UpdateAppointmentStatus) (any, error) {
	switch payload.Status {
	case "cancel":
		_, err := Db.Exec(Ctx, "UPDATE appointments SET status = 'cancelled' WHERE appointment_id = $2", payload.Status, payload.Appointment_id)
		if err != nil {
			log.Println("Failed to update appointment status:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		return map[string]string{"message": "Appointment cancelled successfully"}, nil
	case "completed":
		_, err := Db.Exec(Ctx, "UPDATE appointments SET status = 'completed' WHERE appointment_id = $2", payload.Status, payload.Appointment_id)
		if err != nil {
			log.Println("Failed to update appointment status:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		return map[string]string{"message": "Appointment completed successfully"}, nil
	case "pending":
		_, err := Db.Exec(Ctx, "UPDATE appointments SET status = 'pending' WHERE appointment_id = $2", payload.Status, payload.Appointment_id)
		if err != nil {
			log.Println("Failed to update appointment status:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		return map[string]string{"message": "Appointment status updated to pending"}, nil
	default:
		return nil, errors.New("invalid appointment status")
	}
}

func (a *AdminServer) RescheduleAppointment(data models.RescheduleAppointmentReq) (any, error) {
	_, err := time.Parse(time.RFC3339, data.NewScheduledAt)
	if err != nil {
		return errors.New("invalid datetime format"), nil
	}
	query := `UPDATE appointments SET scheduled_at = $1 WHERE appointment_id = $2`
	_, err = Db.Exec(Ctx, query, data.NewScheduledAt, data.Appointment_id)
	if err != nil {
		log.Println("Error updating appointment schedule:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}

	return map[string]string{"message": "Appointment rescheduled successfully"}, nil
}

func (AdminServer) GetDoctors() (any, error) {
	//rememmebr to modify to fetch using filters
	var doctors []models.Doctor

	rows, err := Db.Query(Ctx, "SELECT doctortag, fullname, date_of_birth, phone_number, gender, specialization, country, yrs_of_experience, price_per_session FROM doctors")
	if err != nil {
		log.Println("Failed to fetch doctors:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer rows.Close()
	for rows.Next() {
		var doctor models.Doctor
		if err := rows.Scan(&doctor.DoctorTag, &doctor.FullName, &doctor.Dob, &doctor.Phone_no, &doctor.Gender, &doctor.Specialization, &doctor.Country, &doctor.YearsOfExperience, &doctor.Price); err != nil {
			log.Println("Failed to scan doctor:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		doctors = append(doctors, doctor)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over doctors:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return doctors, nil
}

func (AdminServer) DeleteDoctor(data models.Doctorreq) error {
	_, err := Db.Exec(Ctx, "DELETE FROM doctors WHERE doctortag = $1", data.DoctorTag)
	if err != nil {
		log.Println("Failed to delete doctor:", err)
		return errors.New(responses.SOMETHING_WRONG)
	}
	return nil
}

func (AdminServer) GetPatients() (any, error) {
	var patients []models.Patient

	rows, err := Db.Query(Ctx, "SELECT usertag, firstname, lastname, email, phone_no, gender, date_of_birth FROM users WHERE role = 'user'")
	if err != nil {
		log.Println("Failed to fetch patients:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	defer rows.Close()
	for rows.Next() {
		var patient models.Patient
		if err := rows.Scan(&patient.UserTag, &patient.Firstname, &patient.Lastname, &patient.Email, &patient.Phone_no, &patient.Gender, &patient.Dob); err != nil {
			log.Println("Failed to scan patient:", err)
			return nil, errors.New(responses.SOMETHING_WRONG)
		}
		patients = append(patients, patient)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over patients:", err)
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	return patients, nil
}

func (AdminServer) GetPatientByUsertag(data models.PatientIdReq) (any, error) {
	var patient models.PatientIdResp
	var firstname, lastname string

	query1 := `select u.usertag, u.firstname, u.lastname, u.phone_no, u.gender, u.date_of_birth , ap.doctor_tag, ap.reason, ap.file_url, ap.status
				from users AS u
				inner join appointments AS ap
				on u.usertag = ap.patient_tag
				where usertag = $1 `
	err := Db.QueryRow(Ctx, query1, data.Usertag).Scan(&patient.UserTag, &firstname, &lastname, &patient.Phone_No, &patient.Gender, &patient.Dob, &patient.Attending_Doctor, &patient.Reason, &patient.File_URL, &patient.Status)
	if err != nil {
		log.Println("Failed to fetch patient:", err)
		if err.Error() == "no rows in result set" {
			return nil, errors.New("patient not found")
		}
		return nil, errors.New(responses.SOMETHING_WRONG)
	}
	patient.Name = firstname + " " + lastname

	return patient, nil

}
