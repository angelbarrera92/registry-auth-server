package authz

import (
	"github.com/angelbarrera92/registry-auth-server/pkg/api"
	"github.com/sirupsen/logrus"
)

type staticAclAuthorization struct {
}

func NewStaticAclAuthorization() *staticAclAuthorization {
	logrus.Info("Load Authorization Controller Successfully\n")
	return &staticAclAuthorization{}
}

func (sa *staticAclAuthorization) Authorize(req *api.AuthRequestInfo) ([]string, error) {
	return nil, nil
}
