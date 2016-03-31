package session

import(
	"errors"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/segmentio/analytics-go"

	"github.com/pachyderm/sandbox/src/util"
)

func Reset(s sessions.Session) {
	user := s.Get("user")
	s.Clear()

	if user != nil {
		s.Set("user", user.(string))
	}
	s.Save()
}

func TagUserSession(analyticsClient *analytics.Client, s sessions.Session) {
	token := util.GenerateUniqueToken()
	s.Set("user", token)
	
	analyticsClient.Identify(&analytics.Identify{AnonymousId: token})
}

func GetUserToken(s sessions.Session) (string, error) {
	value := s.Get("user")

	if value == nil {
		return "", errors.New("No user ID found in session")
	}

	return value.(string), nil
}


