package main

import (
	"db"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vault"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		creds := vault.VaultLogin()
		fmt.Fprintf(w, creds)
		result := db.DatabaseLogin(creds)
		fmt.Fprintf(w, strconv.FormatBool(result))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
