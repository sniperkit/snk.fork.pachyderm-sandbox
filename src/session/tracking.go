package tracking

import(
	"github.com/segmentio/analytics-go"
)

func tokenizeUserSession(s sessions.Session) {
	token := util.GenerateUniqueToken()
	s.Set("user", token)
	
	analyticsClient.Identify(&analytics.Identify{ anonymousID: token})
}

func getUserToken(s sessions.Session) (string, error) {
	value := s.Get("user")

	if value == nil {
		return "", errors.New("No user ID found in session")
	}

	return value.(string), nil
}


