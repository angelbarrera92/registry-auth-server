package token

import (
	"crypto"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/angelbarrera92/registry-auth-server/pkg/api"
	"github.com/docker/distribution/registry/auth/token"
	"github.com/docker/libtrust"
	"github.com/sirupsen/logrus"
)

var t *tokenController

const (
	ALG = crypto.SHA256
)

type tokenController struct {
	config     *api.TokenConfig
	privateKey libtrust.PrivateKey
	publicKey  libtrust.PublicKey
}

func NewTokenController(config *api.TokenConfig) *tokenController {
	if t == nil {
		t = &tokenController{}
		if err := t.loadTokenConfig(config); err != nil {
			logrus.WithField("State", "NewTokenController").Errorf("LoadTokenConfig Failed:%s", err)
			return nil
		}
		logrus.Infof("Load token Controller Successfully")
	}
	return t
}

func (tc *tokenController) loadTokenConfig(config *api.TokenConfig) error {
	tc.config = config
	prikey, pubkey, err := loadCertAndKey(tc.config.CertFile, tc.config.KeyFile)
	if err != nil {
		return err
	}
	tc.privateKey = prikey
	tc.publicKey = pubkey
	return nil
}

func (tc *tokenController) GenerateToken(claim *api.TokenClaim) (string, error) {
	_, sigAlg1, err := tc.privateKey.Sign(strings.NewReader("docker registry. co"), ALG)
	if err != nil {
		return "", err
	}

	// header = base64(json(header))
	header := token.Header{
		Type:       "JWT",
		SigningAlg: sigAlg1,
		KeyID:      tc.privateKey.KeyID(),
	}
	headeJSON, err := json.Marshal(header)
	if err != nil {
		return "", nil
	}

	// payload  = base64(json(payload_struct))
	now := time.Now().Unix()
	claims := &token.ClaimSet{
		Issuer:     tc.config.Issuer,
		Subject:    tc.config.Claim.Account,
		Audience:   claim.Service,
		Expiration: now + tc.config.Expiration,
		NotBefore:  now - 10,
		IssuedAt:   now,
		JWTID:      fmt.Sprintf("%s", rand.Int63()),
		Access:     []*token.ResourceActions{},
	}
	claims.Access = append(claims.Access, &token.ResourceActions{
		Type:    claim.Type,
		Name:    claim.Name,
		Actions: claim.Actions,
	})
	claimJSON, err := json.Marshal(claims)
	if err != nil {

	}

	payload := fmt.Sprintf("%s%s%s", encodeBase64(headeJSON), token.TokenSeparator, encodeBase64(claimJSON))
	sig, sigAlg2, err := tc.privateKey.Sign(strings.NewReader(payload), ALG)
	if err != nil || sigAlg1 != sigAlg2 {
		return "", nil
	}
	return fmt.Sprintf("%s%s%s", payload, token.TokenSeparator, encodeBase64(sig)), nil

}
