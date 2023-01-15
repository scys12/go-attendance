package repositories

import (
	"database/sql"

	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/sirupsen/logrus"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) IUserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) IsEmailExist(email string) (bool, error) {
	var count int
	isExist := false

	if err := u.db.QueryRow(`
		SELECT COUNT(1) 
		FROM users 
		WHERE email = ?
	`, email,
	).Scan(&count); err != nil {
		return isExist, err
	}

	if count > 0 {
		isExist = true
	}
	return isExist, nil
}

func (u *UserRepo) InsertUser(user models.User) (int64, error) {
	var id int64
	res, err := u.db.Exec(`
	INSERT INTO users(
		email,
		password,
		full_name
	) VALUES (
		?,
		?,
		?
	)`,
		user.Email,
		user.Password,
		user.FullName)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "InsertUser",
			"error":    err.Error(),
		}).Errorln("[UserRepo] Failed to insert user")
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "InsertUser",
			"error":    err.Error(),
		}).Errorln("[UserRepo] Failed to get ID")

		return 0, types.ErrFailedGetLastIDDB
	}

	return id, nil
}

func (u *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)

	if err := u.db.QueryRow(`
	SELECT 
		id, 
		email, 
		full_name, 
		password,
		created_at, 
		updated_at
	FROM users
	WHERE
		email = ?
	`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FullName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetUserByEmail",
			"error":    err.Error(),
		}).Errorln("[UserRepo] Failed to get user detail")

		return nil, err
	}

	return user, nil
}
