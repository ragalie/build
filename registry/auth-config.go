package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type HostHeaders map[string]http.Header

type AuthConfig []AuthFile

type AuthFile struct {
	Domains     []string    `json:"domains"`
	Credentials Credentials `json:"credentials"`
}

type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func ReadAuthConfig(directory string) (*AuthConfig, error) {
	authConfig := AuthConfig{}

	err := filepath.Walk(directory, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !validAuthFile(file) {
			return nil
		}

		jsonContents, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("auth-config: %s", err)
		}

		authFile := AuthFile{}
		err = json.Unmarshal(jsonContents, &authFile)
		if err != nil {
			return fmt.Errorf("auth-config: %s", err)
		}

		authConfig = append(authConfig, authFile)

		return nil
	})

	return &authConfig, err
}

func validAuthFile(file os.FileInfo) bool {
	if !file.Mode().IsRegular() {
		return false
	}

	if filepath.Ext(file.Name()) != ".json" {
		return false
	}

	return true
}

func (ac *AuthConfig) HostHeaders() HostHeaders {
	hostHeaders := HostHeaders{}
	for _, authFile := range *ac {
		fakeRequest := http.Request{Header: http.Header{}}
		fakeRequest.SetBasicAuth(authFile.Credentials.User, authFile.Credentials.Password)

		for _, domain := range authFile.Domains {
			hostHeaders[domain] = fakeRequest.Header
		}
	}

	return hostHeaders
}
