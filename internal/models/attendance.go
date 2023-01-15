package models

import "time"

type AttendanceResponse struct {
	AttendanceID   int64  `json:"id,omitempty"`
	AttendanceDate string `json:"attendance_date,omitempty"`
	Message        string `json:"message,omitempty"`
	CheckOutTime   string `json:"check_out_time,omitempty"`
	CheckInTime    string `json:"check_in_time,omitempty"`
}

type AttendancesResponse struct {
	Attendances      []AttendanceResponse `json:"attendances,omitempty"`
	TotalAttendances int                  `json:"total_attendances,omitempty"`
}

type Attendance struct {
	ID             int64     `json:"id,omitempty"`
	CheckIn        time.Time `json:"check_in,omitempty"`
	CheckOut       time.Time `json:"check_out,,omitempty"`
	AttendanceDate time.Time `json:"attendance_date,omitempty"`
	UserID         int64     `json:"user_id,omitempty"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}
