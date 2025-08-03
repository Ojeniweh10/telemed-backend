package models

import (
	"time"

	"gorm.io/datatypes"
)

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

type AnalyticsReq struct {
	Metric string `json:"metric"`
	Month  string `json:"month"`
	Year   string `json:"year"`
}

type AnalyticsResp struct {
	Metric          string  `json:"metric"`
	Month           string  `json:"month"`
	Year            string  `json:"year"`
	Total_amount    float64 `json:"total_amount"`
	Payment_count   int     `json:"payment_count"`
	Average_payment float64 `json:"average_payment"`
}

type Appointment struct {
	ID           string `json:"id"`
	UserTag      string `json:"usertag"`
	DoctorTag    string `json:"doctortag"`
	Scheduled_at string `json:"appointment_date"`
	Reason       string `json:"reason"`
	Status       string `json:"status"`
	Fileurl      string `json:"fileurl"`
	Created_at   string `json:"created_at"`
}

type AppointmentID struct {
	ID string `json:"id"`
}

type Userdata struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
	Dob       string `json:"dob"`
	Phone_no  string `json:"phone_no"`
}

type Doctordata struct {
	Fullname string  `json:"fullname"`
	Price    float64 `json:"price"`
}

type AppointmentIDResp struct {
	UserTag         string    `json:"usertag"`
	DoctorTag       string    `json:"doctortag"`
	Scheduled_At    string    `json:"appointment_date"`
	Reason          string    `json:"reason"`
	File_URL        string    `json:"fileurl"`
	Status          string    `json:"status"`
	Created_At      time.Time `json:"created_at"`
	First_name      string    `json:"firstname"`
	Last_name       string    `json:"lastname"`
	Phone_No        string    `json:"phone_no"`
	Gender          string    `json:"gender"`
	Dob             string    `json:"dob"`
	Doctor_Fullname string    `json:"doctor_fullname"`
	Price           float64   `json:"price"`
}

type Doctorreq struct {
	DoctorTag string
}

type Doctor struct {
	DoctorTag           string         `json:"doctortag"`
	FullName            string         `json:"fullname"`
	Dob                 string         `json:"date_of_birth"`
	Phone_no            string         `json:"phone_number"`
	Gender              string         `json:"gender"`
	Specialization      string         `json:"specialization"`
	Country             string         `json:"country"`
	City                string         `json:"city"`
	YearsOfExperience   int            `json:"yrs_of_experience"`
	Price               float64        `json:"price_per_session"`
	About               string         `json:"about"`
	Availability        datatypes.JSON `json:"availability"` // or []string if unmarshalled
	ProfilePicURL       string         `json:"profile_pic_url"`
	HospitalAffiliation string         `json:"hospital_affiliation"` // from hospital.name
}

type UpdateAppointmentStatus struct {
	Status         string `json:"status"`
	Appointment_id string `json:"appointment_id"`
}

type RescheduleAppointmentReq struct {
	Appointment_id string `json:"appointment_id"`
	NewScheduledAt string `json:"new_scheduled_at"`
}

type Patient struct {
	UserTag   string `json:"usertag"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone_no  string `json:"phone_no"`
	Gender    string `json:"gender"`
	Dob       string `json:"dob"`
}

type PatientIdReq struct {
	Usertag string
}

type PatientIdResp struct {
	UserTag          string `json:"usertag"`
	Name             string `json:"name"`
	Phone_No         string `json:"phone_no"`
	Gender           string `json:"gender"`
	Dob              string `json:"dob"`
	Reason           string `json:"reason"`
	Attending_Doctor string `json:"attending_doctor"`
	File_URL         string `json:"file_url"`
	Status           string `json:"status"`
}
