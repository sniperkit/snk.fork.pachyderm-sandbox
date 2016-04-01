package session

import(
	"errors"
	"fmt"

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
	fmt.Printf("IDd a new user %v\n", token)

	err := analyticsClient.Identify(&analytics.Identify{UserId: token})

	if err != nil {
		fmt.Printf("Segment.io error %v\n", err)
	}
}

func GetUserToken(s sessions.Session) (string, error) {
	value := s.Get("user")

	if value == nil {
		return "", errors.New("No user ID found in session")
	}

	return value.(string), nil
}


