package service

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	vgcrypto "code.vegaprotocol.io/shared/libs/crypto"
	vgrand "code.vegaprotocol.io/shared/libs/rand"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

var ErrSessionNotFound = errors.New("session not found")

//go:generate go run github.com/golang/mock/mockgen -destination mocks/rsa_store_mock.go -package mocks code.vegaprotocol.io/vegawallet/service RSAStore
type RSAStore interface {
	GetRsaKeys() (*RSAKeys, error)
}

type auth struct {
	log *zap.Logger
	// sessionID -> wallet name
	sessions    map[string]string
	privKey     *rsa.PrivateKey
	pubKey      *rsa.PublicKey
	tokenExpiry time.Duration

	mu sync.Mutex
}

func NewAuth(log *zap.Logger, cfgStore RSAStore, tokenExpiry time.Duration) (*auth, error) {
	keys, err := cfgStore.GetRsaKeys()
	if err != nil {
		return nil, err
	}
	priv, err := jwt.ParseRSAPrivateKeyFromPEM(keys.Priv)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse private RSA key: %w", err)
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(keys.Pub)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse public RSA key: %w", err)
	}

	return &auth{
		sessions:    map[string]string{},
		privKey:     priv,
		pubKey:      pub,
		log:         log,
		tokenExpiry: tokenExpiry,
	}, nil
}

type Claims struct {
	jwt.StandardClaims
	Session string
	Wallet  string
}

func (a *auth) NewSession(name string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	expiresAt := time.Now().Add(a.tokenExpiry)

	session := genSession()

	claims := &Claims{
		Session: session,
		Wallet:  name,
		StandardClaims: jwt.StandardClaims{
			// these are seconds
			ExpiresAt: jwt.NewTime((float64)(expiresAt.Unix())),
			Issuer:    "vega wallet",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	ss, err := token.SignedString(a.privKey)
	if err != nil {
		a.log.Error("unable to sign token", zap.Error(err))
		return "", err
	}

	a.sessions[session] = name
	return ss, nil
}

// VerifyToken returns the wallet name associated for this session.
func (a *auth) VerifyToken(token string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	claims, err := a.parseToken(token)
	if err != nil {
		return "", err
	}

	w, ok := a.sessions[claims.Session]
	if !ok {
		return "", ErrSessionNotFound
	}

	return w, nil
}

func (a *auth) Revoke(token string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	claims, err := a.parseToken(token)
	if err != nil {
		return err
	}

	_, ok := a.sessions[claims.Session]
	if !ok {
		return ErrSessionNotFound
	}
	delete(a.sessions, claims.Session)
	return nil
}

func (a *auth) parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return a.pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// ExtractToken this is public for testing purposes.
func ExtractToken(f func(string, http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			writeError(w, ErrInvalidOrMissingToken, http.StatusBadRequest)
			return
		}
		splitToken := strings.Split(token, "Bearer")
		if len(splitToken) != 2 || len(splitToken[1]) == 0 {
			writeError(w, ErrInvalidOrMissingToken, http.StatusBadRequest)
			return
		}
		f(strings.TrimSpace(splitToken[1]), w, r, ps)
	}
}

func genSession() string {
	return hex.EncodeToString(vgcrypto.Hash(vgrand.RandomBytes(10)))
}

func writeError(w http.ResponseWriter, e error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	buf, err := json.Marshal(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
