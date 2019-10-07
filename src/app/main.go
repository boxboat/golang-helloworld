package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		creds := VaultLogin()
		fmt.Fprintf(w, creds)
		result := DatabaseLogin(creds)
		fmt.Fprintf(w, strconv.FormatBool(result))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}

type vaultAuthSend struct {
	Jwt  string `json:"jwt"`
	Role string `json:"role"`
}

type vaultAuthReturn struct {
	Auth struct {
		Accessor      string `json:"accessor"`
		ClientToken   string `json:"client_token"`
		LeaseDuration int    `json:"lease_duration"`
		Metadata      struct {
			Role                     string `json:"role"`
			ServiceAccountName       string `json:"service_account_name"`
			ServiceAccountNamespace  string `json:"service_account_namespace"`
			ServiceAccountSecretName string `json:"service_account_secret_name"`
			ServiceAccountUID        string `json:"service_account_uid"`
		} `json:"metadata"`
		Policies  []string `json:"policies"`
		Renewable bool     `json:"renewable"`
	} `json:"auth"`
}

type vaultCreds struct {
	Data struct {
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"data"`
}

func VaultLogin() string {
	jwt := retrieveJwt()
	token := loginVault(jwt)
	creds := dbCred(token)
	return creds
}

func retrieveJwt() string {
	f, _ := ioutil.ReadFile("")
	return string(f)
}

func loginVault(jwt string) string {
	url := "http://vault.vault.svc.cluster.local:8200/v1/auth/kubernetes/login"
	toSendData := vaultAuthSend{jwt, "HelloWorld"}
	data, _ := json.Marshal(toSendData)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	client := &http.Client{Timeout: time.Second * 10}
	resp, _ := client.Do(req)
	respUnmarshal := vaultAuthReturn{}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	json.Unmarshal(bodyBytes, &respUnmarshal)

	token := respUnmarshal.Auth.ClientToken
	return token
}

func dbCred(token string) string {
	url := "http://vault.vault.svc.cluster.local:8200/v1/database/creds/mysql/test"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Vault-Token", token)
	client := &http.Client{Timeout: time.Second * 10}
	resp, _ := client.Do(req)
	respUnmarshal := vaultCreds{}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	json.Unmarshal(bodyBytes, &respUnmarshal)

	return respUnmarshal.Data.Username + ":" + respUnmarshal.Data.Password
}

func DatabaseLogin(creds string) bool {
	db, err := sql.Open("mysql", creds+"@tcp(mysql.mysql.svc.cluster.local)/main")
	defer db.Close()
	if err != nil {
		return false
	} else {
		return true
	}
}
