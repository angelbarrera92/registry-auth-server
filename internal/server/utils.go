package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/angelbarrera92/registry-auth-server/pkg/api"
	"github.com/sirupsen/logrus"
)

func parseRequest(request *http.Request) (*authRequestInfo, error) {
	username, password, ok := request.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("Username or password must be supplied")
	}
	account := request.FormValue("account")
	if account != username {
		logrus.Warningf("Username:%s and Account:%s is not same", username, account)
	}

	service := request.FormValue("service")
	authReq := &authRequestInfo{
		Username:     username,
		Password:     password,
		Service:      service,
		Account:      account,
		Actions:      nil,
		Type:         "",
		ResourceName: "",
	}
	parts := strings.Split(request.URL.Query().Get("scope"), ":")

	if len(parts) > 0 {
		authReq.Type = parts[0]
	}
	if len(parts) > 1 {
		authReq.ResourceName = parts[1]
	}
	if len(parts) > 2 {
		authReq.Actions = strings.Split(parts[2], ",")
	}
	return authReq, nil
}

func authRequestHandler(info *authRequestInfo) *api.AuthRequestInfo {
	authReq := &api.AuthRequestInfo{
		Account:      info.Account,
		Actions:      info.Actions,
		ResourceName: info.ResourceName,
		Type:         info.Type,
	}
	return authReq
}

func generateTokenClaimHandler(info *authRequestInfo) *api.TokenClaim {
	sr := &api.TokenClaim{
		Type:    info.Type,
		Account: info.Account,
		Name:    info.ResourceName,
		Actions: info.Actions,
		Service: info.Service,
	}
	return sr
}
