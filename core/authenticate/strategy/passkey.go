package strategy

import (
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

const (
	PasskeyAuthMethod string = "passkey"
)

type UserData struct {
	Email       string
	Id          string
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
	CreatedAt   time.Time
}

func NewPassKeyUser(email string) *UserData {
	user := &UserData{}
	user.Id = uuid.New().String()
	user.Email = email
	user.DisplayName = strings.Split(email, "@")[0]
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

type PasskeyStartRequest struct {
	Type    string
	Options any
}

func (p *PasskeyStartRequest) ForRegisteration(opt *protocol.CredentialCreation) {
	p.Type = "create"
	p.Options = opt
}

func (p *PasskeyStartRequest) ForLogin(opt *protocol.CredentialAssertion) {
	p.Type = "get"
	p.Options = opt
}
