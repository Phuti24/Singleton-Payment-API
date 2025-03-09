package security

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func EncryptWithHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error encrypting password: %v", err)
	}
	return string(hash), nil
}

func DecryptWithHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func EncryptHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := EncryptWithHash(request.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encrypting password: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"hashed_password": hashedPassword})
}

func DecryptHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Password       string `json:"password"`
		HashedPassword string `json:"hashed_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	isValid := DecryptWithHash(request.HashedPassword, request.Password)
	if isValid {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Password is valid"})
	} else {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
	}
}
