package wallet

import (
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type auth struct {
	log *zap.Logger
	// sessionID -> wallet name
	sessions    map[string]string
	privKey     *rsa.PrivateKey
	pubKey      *rsa.PublicKey
	tokenExpiry time.Duration

	mu sync.Mutex
}

func NewAuth(log *zap.Logger, rootPath string, tokenExpiry time.Duration) (*auth, error) {
	// get rsa keys
	pubBuf, privBuf, err := readRsaKeys(rootPath)
	if err != nil {
		return nil, err
	}
	priv, err := jwt.ParseRSAPrivateKeyFromPEM(privBuf)
	if err != nil {
		return nil, err
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubBuf)
	if err != nil {
		return nil, err
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

func (a *auth) NewSession(walletname string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	expiresAt := time.Now().Add(a.tokenExpiry)

	session := genSession()
	// Create the Claims
	claims := &Claims{
		Session: session,
		Wallet:  walletname,
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

	// all good up to now, insert the new session
	a.sessions[session] = walletname
	return ss, nil
}

// VerifyToken returns the walletname associated for this session
func (a *auth) VerifyToken(token string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// first parse the token
	claims, err := a.parseToken(token)
	if err != nil {
		return "", err
	}

	wallet, ok := a.sessions[claims.Session]
	if !ok {
		return "", ErrSessionNotFound
	}

	return wallet, nil
}

func (a *auth) Revoke(token string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	claims, err := a.parseToken(token)
	if err != nil {
		return err
	}

	// extract session from the token
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

// ExtractToken this is public for testing purposes
func ExtractToken(f func(string, http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")
		if len(token) <= 0 {
			// invalid token, return an error here
			writeError(w, ErrInvalidOrMissingToken, http.StatusBadRequest)
			return
		}
		splitToken := strings.Split(token, "Bearer")
		if len(splitToken) != 2 || len(splitToken[1]) <= 0 {
			// invalid token, return an error here
			writeError(w, ErrInvalidOrMissingToken, http.StatusBadRequest)
			return
		}
		// then call the function
		f(strings.TrimSpace(splitToken[1]), w, r, ps)
	}
}

func genSession() string {
	hasher := sha3.New256()
	hasher.Write([]byte(randSeq(10)))
	return hex.EncodeToString(hasher.Sum(nil))
}

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
