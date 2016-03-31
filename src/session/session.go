package session

import(
	"github.com/gin-gonic/contrib/sessions"
	"github.com/segmentio/analytics-go"
)

func Reset(s sessions.Session) {
	user = s.Get("user")
	s.Clear()

	if user != nil {
		s.Set("user", user.(string))
	}
	s.Save()
}

func TagUserSession(s sessions.Session) {
	token := util.GenerateUniqueToken()
	s.Set("user", token)
	
	analyticsClient.Identify(&analytics.Identify{ anonymousID: token})
}

func GetUserToken(s sessions.Session) (string, error) {
	value := s.Get("user")

	if value == nil {
		return "", errors.New("No user ID found in session")
	}

	return value.(string), nil
}


