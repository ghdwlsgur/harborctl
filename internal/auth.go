package internal

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

var (
	homeDir, _            = os.UserHomeDir()
	_credentialParentFile = filepath.Join(homeDir + "/.harborctl")
	_credentialFile       = _credentialParentFile + "/credentials"
)

type User struct {
	Username string
	Password string
}

type Auth interface {
	setUsername(username string) error
	setPassword(password string) error
	GetUsername() string
	GetPassword() string
	createCredentialParentFile() error
	createCredentialFile() error
	setCredential() error
	Verify() bool
}

var _ Auth = (*User)(nil)

func (u *User) setUsername(username string) error {
	if username == "" {
		return errors.New("failed to set username because username is empty")
	}
	u.Username = username
	return nil
}

func (u *User) setPassword(password string) error {
	if password == "" {
		return errors.New("failed to set password because password is empty")
	}
	u.Password = password
	return nil
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetPassword() string {
	return u.Password
}

func (u User) createCredentialParentFile() error {
	if _, err := os.Stat(_credentialParentFile); os.IsNotExist(err) {
		err := os.Mkdir(_credentialParentFile, 0755)
		if err != nil {
			return errors.New("createCredentialParentFile")
		}
	}
	return nil
}

func (u User) createCredentialFile() error {
	if _, err := os.Stat(_credentialFile); os.IsNotExist(err) {
		f, err := os.Create(_credentialFile)
		if err != nil {
			return errors.New("createCredentialFile")
		}
		f.Close()
	}
	return nil
}

func (u *User) setCredential() error {
	if err := u.createCredentialParentFile(); err != nil {
		return fmt.Errorf("setCredential: %w", err)
	}

	if err := u.createCredentialFile(); err != nil {
		return fmt.Errorf("setCredential: %w", err)
	}

	f, err := os.Create(_credentialFile)
	if err != nil {
		return fmt.Errorf("os.Create - setCredential: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "username=%s\n", u.GetUsername())
	fmt.Fprintf(f, "password=%s\n", u.GetPassword())

	return nil
}

func (u User) Verify() bool {
	if _, err := os.Stat(_credentialFile); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func (u *User) Parsing() (*User, error) {
	credential, err := os.ReadFile(_credentialFile)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile - Parsing: %w", err)
	}

	lines := strings.Split(string(credential), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "=")

		if len(parts) > 1 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if key == "username" {
				err := u.setUsername(value)
				if err != nil {
					return nil, fmt.Errorf("Parsing: %w", err)
				}
			}

			if key == "password" {
				err := u.setPassword(value)
				if err != nil {
					return nil, fmt.Errorf("Parsing: %w", err)
				}
			}
		}
	}

	return u, nil
}

func (u *User) GetBasicAuth() string {
	s := u.GetUsername() + ":" + u.GetPassword()
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (u *User) Login() error {

	var username string
	promptUsername := &survey.Input{
		Message: "Please enter your harbor username",
	}
	survey.AskOne(promptUsername, &username)
	err := u.setUsername(username)
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	var password string
	promptPassword := &survey.Password{
		Message: "Please enter your harbor password",
	}
	survey.AskOne(promptPassword, &password)
	err = u.setPassword(password)
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}

	err = u.setCredential()
	if err != nil {
		return fmt.Errorf("Login: %w", err)
	}
	return nil
}

func GetHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
