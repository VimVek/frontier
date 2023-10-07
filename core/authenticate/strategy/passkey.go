package strategy

import (
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

const (
	PasskeyAuthMethod string = "passkey"
)

type UserData struct {
	Id          string
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
	CreatedAt   time.Time
}

func NewPassKeyUser(id string) *UserData {
	user := &UserData{}
	user.Id = id
	user.Name = extractUsername(id)
	user.DisplayName = extractUsername(id)
	return user
}

func (u *UserData) WebAuthnID() []byte {
	return []byte(u.Id)
}

func (u *UserData) WebAuthnName() string {
	return u.Name
}

func (u *UserData) WebAuthnDisplayName() string {
	return u.DisplayName
}
func (u *UserData) WebAuthnIcon() string {
	return ""
}
func (u *UserData) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func extractUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}
