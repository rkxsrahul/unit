package methods

import (
	"testing"
)

// test Slugify Email
func TestSlugifyEmail(t *testing.T) {

	email := "test@testing.com"
	slugify_email := SlugifyEmail(email)
	if len(slugify_email) > 5 {
		t.Error("test case fail")
	}

}

func TestRandomString(t *testing.T) {
	RandomString(5)

}
