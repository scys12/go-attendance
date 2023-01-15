package sessions

import (
	"net/http"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/pkg/tokenizer"
)

var (
	lock           = &sync.Mutex{}
	sessionManager *scs.SessionManager
)

const sessionStr = "session"

func GetSessionInstance() *scs.SessionManager {
	if sessionManager == nil {
		lock.Lock()
		defer lock.Unlock()
		if sessionManager == nil {
			sessionManager = scs.New()
			sessionManager.Lifetime = 24 * time.Hour
		}
	}

	return sessionManager
}

func CreateSessionAuthentication(r *http.Request, user *models.User) error {
	sess := GetSessionInstance()
	token, err := tokenizer.CreateToken(user.ID, user.Email)
	if err != nil {
		return err
	}
	sess.Put(r.Context(), sessionStr, token)
	return nil
}

func GetSessionToken(r *http.Request) string {
	sess := GetSessionInstance()
	token := sess.GetString(r.Context(), sessionStr)
	return token
}

func DeleteSessionToken(r *http.Request) {
	sess := GetSessionInstance()
	sess.Destroy(r.Context())
}
