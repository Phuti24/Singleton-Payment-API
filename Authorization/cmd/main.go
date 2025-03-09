package main

import (
	config "LockBox/api/api_config"
	"LockBox/api/credentials"
	"LockBox/api/credentialsdb"
	"LockBox/api/security"
	"fmt"
	"log"
	"net/http"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configs: %v", err)
	}

	log.Printf(conf.Database.Name, "8080")

	http.HandleFunc("/auth/register", credentials.RegisterHandler)
	http.HandleFunc("/auth/login", credentials.LoginHandler)

	http.HandleFunc("/store/create", credentialsdb.CreateHandler)
	http.HandleFunc("/store/delete", credentialsdb.DeleteHandler)
	http.HandleFunc("/store/update", credentialsdb.UpdateHandler)
	http.HandleFunc("/store/get", credentialsdb.FetchWithIdHandler)

	http.HandleFunc("/security/encrypt", security.EncryptHandler)
	http.HandleFunc("/security/decrypt", security.DecryptHandler)

	fmt.Printf("Server Listening On Port:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
