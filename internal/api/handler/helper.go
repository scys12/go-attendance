package handler

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/types"
)

func ValidateRequest(validate *validator.Validate, err error) (errs []error) {
	if err == nil {
		return nil
	}

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}

func TransformAttendancesToAttendancesResponse(attendances []models.Attendance) models.AttendancesResponse {
	totalAttendances := len(attendances)
	atdncResp := make([]models.AttendanceResponse, totalAttendances)

	for i := range attendances {
		atdncResp[i].AttendanceID = attendances[i].ID
		atdncResp[i].AttendanceDate = attendances[i].AttendanceDate.Format(types.FormatOnlyDate)

		if attendances[i].CheckIn.IsZero() {
			atdncResp[i].CheckInTime = "-"
		} else {
			atdncResp[i].CheckInTime = attendances[i].CheckIn.Format(types.FormatTime)
		}

		if attendances[i].CheckOut.IsZero() {
			atdncResp[i].CheckOutTime = "-"
		} else {
			atdncResp[i].CheckOutTime = attendances[i].CheckOut.Format(types.FormatTime)
		}
	}

	resp := models.AttendancesResponse{
		Attendances:      atdncResp,
		TotalAttendances: totalAttendances,
	}
	return resp
}
