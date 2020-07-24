package token

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/docker/libtrust"
)

func encodeBase64(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

func loadCertAndKey(certFile, keyFile string) (libtrust.PrivateKey, libtrust.PublicKey, error) {
	if certFile == "" || keyFile == "" {
		return nil, nil, fmt.Errorf("loadCertAndKey: certFile or keyfile must be supplied")
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("load %s Failed: %s", certFile, err.Error())
	}
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, nil, fmt.Errorf("parse certififace %s Failed: %s", certFile, err.Error())
	}
	pubk, err := libtrust.FromCryptoPublicKey(x509Cert.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Gather Publickey from %s Failed: %s", keyFile, err.Error())
	}
	prik, err := libtrust.FromCryptoPrivateKey(cert.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Gather Private Key from %s Failed: %s", keyFile, err.Error())
	}
	return prik, pubk, nil
}
