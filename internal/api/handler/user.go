package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/server"
	"github.com/scys12/go-attendance/internal/service"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/scys12/go-attendance/pkg/sessions"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(svc service.IUserService) UserHandler {
	return UserHandler{
		userService: svc,
	}
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "Login",
		}).Info("[UserHandler] Unable to decode request body")

		server.RenderError(w, http.StatusBadRequest, types.ErrBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(loginReq)
	errs := ValidateRequest(validate, err)

	if errs != nil {
		logrus.WithFields(logrus.Fields{
			"error":    errs[0],
			"function": "Login",
		}).Info("[UserHandler] Request validation failed")

		server.RenderError(w, http.StatusBadRequest, errs[0])
		return
	}

	user, err := u.userService.LoginUser(loginReq)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "Login",
		}).Errorln("[UserHandler] Failed to create session for user")
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	sesErr := sessions.CreateSessionAuthentication(r, user)
	if sesErr != nil {
		logrus.WithFields(logrus.Fields{
			"error":    sesErr.Error(),
			"function": "Login",
		}).Errorln("[UserHandler] Failed to create session for user")

		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	resp := &models.LoginResponse{
		ID:      int64(user.ID),
		Message: types.SuccessLogin,
	}
	server.RenderResponse(w, http.StatusCreated, resp)
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		server.RenderError(w, http.StatusBadRequest, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(registerReq)
	errs := ValidateRequest(validate, err)

	if errs != nil {
		logrus.WithFields(logrus.Fields{
			"error":    errs[0],
			"function": "Register",
		}).Info("[UserHandler] Request validation failed")

		server.RenderError(w, http.StatusBadRequest, errs[0])
		return
	}

	user, err := u.userService.RegisterUser(registerReq)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"function": "Register",
		}).Errorln("[UserHandler] Failed to register user")

		server.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	resp := models.RegisterResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
	}

	server.RenderResponse(w, http.StatusCreated, resp)
}

func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessions.DeleteSessionToken(r)
	resp := models.LogoutResponse{
		Message: types.SuccessLogout,
	}
	server.RenderResponse(w, http.StatusCreated, resp)
}
