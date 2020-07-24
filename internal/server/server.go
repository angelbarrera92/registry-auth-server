package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/angelbarrera92/registry-auth-server/internal/auth/authn"
	"github.com/angelbarrera92/registry-auth-server/internal/auth/authz"
	"github.com/angelbarrera92/registry-auth-server/internal/configs"
	"github.com/angelbarrera92/registry-auth-server/internal/token"
	"github.com/angelbarrera92/registry-auth-server/pkg/api"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type registryAuthServer struct {
	Authenticator api.Authenticator
	Authorization api.Authorization
	Token         api.IToken

	//
	address      string
	pem, key     string
	secureModule bool
}

// NewRegistryAuthServer return a instance of registry auth server
func NewRegistryAuthServer(cfg *configs.Configs) *registryAuthServer {
	tcfg := &api.TokenConfig{
		CertFile: cfg.Tls.Cert,
		KeyFile:  cfg.Tls.Key,
		Issuer:   cfg.Token.Issuer,
	}
	authnController := authn.NewStaticBasicAuthenticator()
	if authnController == nil {
		logrus.WithField("State", "Build Authenticator Failed").Errorf("Authenticator is nil")
		return nil
	}
	authoController := authz.NewStaticAclAuthorization()
	if authoController == nil {
		logrus.WithField("State", "Build NewStaticAclAuthorization Failed").Errorf("NewStaticAclAuthorization is nil")
		return nil
	}
	tokenController := token.NewTokenController(tcfg)
	if tokenController == nil {
		logrus.WithField("State", "Build NewTokenController Failed").Errorf("NewTokenController is nil")
		return nil
	}
	return &registryAuthServer{
		Authenticator: authnController,
		Authorization: authoController,
		Token:         tokenController,
		pem:           cfg.Tls.Cert,
		key:           cfg.Tls.Key,
		address:       cfg.Server.Address + ":" + cfg.Server.Port,
		secureModule:  cfg.SecureModule,
	}
}

//
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logrus.Infof("Uri:%s\tMethod:%s", request.URL.String(), request.Method)
		if request.Method == "OPTION" {
			// todo to deal cors
			return
		}
		next.ServeHTTP(writer, request)
	})
}

// run registry server
func (rs *registryAuthServer) Run(ctx context.Context) error {
	logrus.Info("Docker registry token server begin running.....")
	route := mux.NewRouter()
	route.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello world")
		return
	})
	route.Handle("/auth", rs)
	route.Use(loggingMiddleware)
	var errchan chan error
	defer close(errchan)
	var err error

	if rs.secureModule {
		logrus.Infof("Docker Registry Auth server Run as TLS Module")
		err = http.ListenAndServeTLS(rs.address, rs.pem, rs.key, route)
	} else {
		logrus.Infof("Docker registry auth server run as insecure module")
		err = http.ListenAndServe(rs.address, route)
	}
	select {
	case errchan <- err:
		return <-errchan
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (rs *registryAuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// parse request
	authReq, err := parseRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse Request authRequest Failed:%s", err.Error()), http.StatusBadRequest)
		return
	}
	logrus.Debugf("[1].parse request Successfully")

	// auth username and password
	if ok, err := rs.Authenticator.Authenticate(authReq.Username, authReq.Password); !ok {
		logrus.Errorf("[!]Authenticated Failed:%s", err.Error())
		http.Error(w, fmt.Sprintf("[!]Authenticated Failed:%s", err.Error()), http.StatusUnauthorized)
		return
	}
	logrus.Infof("[2]Auth username and password Successfully")

	// Acl
	_, err = rs.Authorization.Authorize(authRequestHandler(authReq))
	if err != nil {
		http.Error(w, fmt.Sprintf("Authorization User Acl Faile: %s", err.Error()), http.StatusForbidden)
		return
	}
	logrus.Infof("[3]AUthorization action successfully")

	// token
	tokenstring, err := rs.Token.GenerateToken(generateTokenClaimHandler(authReq))
	if err != nil {
		http.Error(w, fmt.Sprintf("Generate token for %s Failed: %s", authReq.Username, err), http.StatusInternalServerError)
		return
	}
	logrus.Infof("[4]Token Generation Successfully: ")
	data, _ := json.Marshal(&map[string]string{"access_token": tokenstring, "token": tokenstring})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
