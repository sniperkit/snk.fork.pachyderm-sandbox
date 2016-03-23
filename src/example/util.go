package example

import(
	"strings"
	"fmt"

	"github.com/satori/go.uuid"
)

func generateUniqueName(prefix string) string {
	unique_suffix := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	unique_name := uniqueNameFromToken(prefix, unique_suffix[0:12])

	return unique_name
}

func uniqueNameFromToken(prefix string, token string) string {
	return fmt.Sprintf("%v-%v", prefix, token)
}
