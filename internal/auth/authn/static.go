package authn

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type staticBasicAuthenticator struct {
}

func NewStaticBasicAuthenticator() *staticBasicAuthenticator {
	logrus.Info("Load Authenticator Controller Successfully")
	return &staticBasicAuthenticator{}
}

func (authn *staticBasicAuthenticator) Authenticate(username, password string) (bool, error) {
	if pwd, ok := whilteList[username]; ok {
		if pwd == password {
			return true, nil
		}
	}
	return false, fmt.Errorf("Username or password is invalid")
}
