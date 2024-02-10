package service

import (
	"encoding/json"
	"fmt"
	"github.com/dark-enstein/vault/internal/model"
	"github.com/dark-enstein/vault/internal/store"
	"net/http"
)

var (
	Tokenize     = "/tokenize"
	Introduction = "/new"
)

var (
	newVaultHandleFunc = func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Welcome to Data Vault")
	}
	tokenize = func(w http.ResponseWriter, r *http.Request) {
		var resp model.Response
		var token model.Tokenize

		w.Header().Set("Content-Type", "application/json")
		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()
		defer r.Body.Close()

		// Check that json is a valid model.Tokenize structure
		if err := jsonDecoder.Decode(&token); err != nil {
			resp.Error = err.Error()
			// return 400 status codes
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Connect to store interface
		store.NewSyncMap()

		// tokenize logic

		// store in store

		// generate response

		// set header and return
	}
)

type VaultHandler map[string]func(w http.ResponseWriter, r *http.Request)

func NewVaultHandler() *VaultHandler {
	vh := make(VaultHandler, 10)
	vh[Introduction] = newVaultHandleFunc
	vh[Tokenize] = newVaultHandleFunc
	//vh[Introduction] = newVaultHandleFunc
	return &vh
}
