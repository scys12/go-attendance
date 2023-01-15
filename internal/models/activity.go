package models

type Activity struct {
	ID           int64  `json:"id,omitempty"`
	AttendanceID int64  `json:"attendance_id,omitempty"`
	UserID       int64  `json:"user_id,omitempty"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ActivityRequest struct {
	Description string `json:"description,omitempty" validate:"required"`
}

type ActivityResponse struct {
	ID           int64  `json:"id,omitempty"`
	AttendanceID int64  `json:"attendance_id,omitempty"`
	Description  string `json:"description,omitempty"`
}

type ActivitiesResponse struct {
	Activities      []Activity `json:"activities,omitempty"`
	TotalActivities int        `json:"total_activities,omitempty"`
}
