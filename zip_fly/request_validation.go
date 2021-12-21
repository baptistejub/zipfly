package zipfly

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func validateExpiration(expires string) bool {
	expiresTimestamp, err := strconv.ParseInt(expires, 10, 64)
	if err != nil || time.Now().Unix() > expiresTimestamp {
		return false
	}

	return true
}

func validateHMAC(message, signature, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	hexMac := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal(signature, []byte(hexMac))
}

func extractSignatureAndExpiresFromHeaders(req *http.Request) (string, string) {
	signature := req.Header.Get("x-zipfly-signature")
	expires := req.Header.Get("x-zipfly-expires")
	return signature, expires
}

func extractSignatureAndExpiresFromQueryString(req *http.Request) (string, string) {
	query := req.URL.Query()
	if query["signature"] == nil || query["expires"] == nil {
		return "", ""
	}

	return query["signature"][0], query["expires"][0]
}

func (s *Server) mustValidateRequestSignature() bool {
	return s.options.ValidateSignature || s.environment == "production"
}

func (s *Server) validateGetRequestSignature(req *http.Request) bool {
	if s.mustValidateRequestSignature() {
		signature, expires := extractSignatureAndExpiresFromQueryString(req)

		query := req.URL.Query()
		query.Del("signature")
		u, _ := url.Parse(s.options.PublicUrl)
		u.Path = req.URL.Path
		u.RawQuery = query.Encode()

		return s.validateSignature(signature, expires, u.String())
	}

	return true
}

func (s *Server) validatePostRequestSignature(req *http.Request, body []byte) bool {
	if s.mustValidateRequestSignature() {
		signature, expires := extractSignatureAndExpiresFromHeaders(req)

		message := expires + ":" + string(body)

		return s.validateSignature(signature, expires, message)
	}

	return true
}

func (s *Server) validateSignature(signature, expires, message string) bool {
	if signature == "" || expires == "" {
		return false
	}

	if !validateExpiration(expires) {
		return false
	}

	return validateHMAC([]byte(message), []byte(signature), []byte(s.options.SigningSecret))
}
