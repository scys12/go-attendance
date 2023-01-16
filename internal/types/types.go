package types

import "errors"

var (
	ErrFailedInsertActivity = errors.New("failed to insert activity")
	ErrFailedUpdateActivity = errors.New("failed to update activity")
	ErrFailedDeleteActivity = errors.New("failed to delete activity")
	ErrGetActivityList      = errors.New("failed to get list of activity")
	ErrGetAttendanceList    = errors.New("failed to get list of attendance")

	ErrFailedGetLastIDDB  = errors.New("failed to get last inserted id")
	ErrNoRowsAffected     = errors.New("there is not any rows affected")
	ErrFailedRowsAffected = errors.New("failed to get rows affected")

	ErrQueryParamActivityDate = errors.New("error query param activity_date not found")
	ErrBadRequest             = errors.New("error bad request")
	ErrNotCheckIn             = errors.New("user has not checked in")
	ErrAlreadyCheckIn         = errors.New("user has checked in")
	ErrAlreadyCheckOut        = errors.New("user has checked out")

	ErrTokenAuthentication = errors.New("token authentication failed")
	ErrEmailPasswordWrong  = errors.New("email/password is wrong")
	ErrEmailNotFound       = errors.New("email is not found")
	ErrEmailExists         = errors.New("email exists in the system")

	ErrTokenUnexpectedMethod = errors.New("unexpected signed methods")
	ErrTokenFailedClaim      = errors.New("failed when claiming token")
)

const (
	SuccessLogout   = "successfully logout from account"
	SuccessLogin    = "successfully login to account"
	SuccessCheckIn  = "successfully check in today"
	SuccessCheckOut = "successfully check out today"
)

const (
	FormatTime     = "15:04:05"
	FormatOnlyDate = "2006-01-02"
)
