package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dark-enstein/vault/internal/model"
	"github.com/dark-enstein/vault/internal/tokenize"
	"net/http"
)

const (
	CodeSuccess = iota
	CodeInternalServerError
	CodeInvalidRequest
	CodeMethodNotAllowed
	CodeRequestTimeout
)

var (
	KeyDelimiter = tokenize.KeyDelimiter
)

var (
	Tokenize     = "/tokenize"
	Detokenize   = "/detokenize"
	GetTokens    = "/alltokens"
	Introduction = "/new"
)

var (
	ErrMethodNotAllowed = "method not allowed"
)

type VaultHandler map[string]func(w http.ResponseWriter, r *http.Request)

func NewVaultHandler(ctx context.Context, srv *Service) *VaultHandler {
	vh := make(VaultHandler, 10)
	vh[Introduction] = VaultHandlerFunc(srv)
	vh[Tokenize] = TokenizeHandlerFunc(srv)
	vh[Detokenize] = DetokenizeHandlerFunc(srv)
	vh[GetTokens] = GetTokensHandler(srv)
	//vh[Introduction] = newVaultHandleFunc
	return &vh
}

func VaultHandlerFunc(srv *Service) func(w http.ResponseWriter, r *http.Request) {
	log := srv.log
	return func(w http.ResponseWriter, r *http.Request) {
		log.Logger().Info().Msg(fmt.Sprintf("received a request on %s", Introduction))
		fmt.Fprintf(w, "Welcome to Data Vault")
		log.Logger().Info().Msgf("VaultHandlerFunc completed with no errors")
	}
}

func GetTokensHandler(srv *Service) func(w http.ResponseWriter, r *http.Request) {
	log := srv.log
	return func(w http.ResponseWriter, r *http.Request) {
		log.Logger().Info().Msg(fmt.Sprintf("received a request on %s", GetTokens))
		var resp model.Response
		var err error

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			resp.Error = append(resp.Error, ErrMethodNotAllowed)
			log.Logger().Error().Msg(ErrMethodNotAllowed)
			resp.Code = CodeMethodNotAllowed
			json.NewEncoder(w).Encode(resp)
			return
		}

		//reqCtx := context.Background()

		// tokenize logic
		manager := srv.manager

		// user request valid, not proceed to process
		tokens, err := manager.GetAllTokens()
		if err != nil {
			resp.Error = append(resp.Error, err.Error())
			log.Logger().Error().Msg(err.Error())
			resp.Code = CodeInternalServerError
			// return 400 status codes
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// generate response
		tokenStruct := &model.All{
			tokens,
		}
		resp.Resp = tokenStruct
		resp.Code = CodeSuccess

		// set header and return
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
}

func DetokenizeHandlerFunc(srv *Service) func(w http.ResponseWriter, r *http.Request) {
	log := srv.log
	return func(w http.ResponseWriter, r *http.Request) {
		log.Logger().Info().Msg(fmt.Sprintf("received a request on %s", Detokenize))
		var resp model.Response
		var detoken model.Detokenize
		var err error

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			resp.Error = append(resp.Error, ErrMethodNotAllowed)
			log.Logger().Error().Msg(ErrMethodNotAllowed)
			resp.Code = CodeMethodNotAllowed
			json.NewEncoder(w).Encode(resp)
			return
		}

		//reqCtx := context.Background()

		w.Header().Set("Content-Type", "application/json")
		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()
		defer r.Body.Close()

		// Check that json is a valid model.Tokenize structure
		if err = jsonDecoder.Decode(&detoken); err != nil {
			resp.Error = append(resp.Error, err.Error())
			log.Logger().Error().Msg(err.Error())
			resp.Code = CodeInvalidRequest
			// return 400 status codes
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		var decryptedStr string
		var children []*model.ChildReceipt

		// tokenize logic
		manager := srv.manager

		// user request valid, not proceed to process
		parentKey := detoken.ID
		for i := 0; i < len(detoken.Data); i++ {
			var found bool
			childKey := detoken.Data[i].Key
			combinedKeyName := tokenize.GetCombinedKey(parentKey, childKey)
			found, decryptedStr, err = manager.Detokenize(combinedKeyName, detoken.Data[i].Value)
			if err != nil || !found {
				resp.Error = append(resp.Error, fmt.Sprintf("error with key %s.%s: %s", parentKey, childKey, err.Error()))
				log.Logger().Error().Msg(fmt.Sprintf("error with key %s.%s: %s", parentKey, childKey, err.Error()))
				resp.Code = CodeInternalServerError
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}
			children = append(children, &model.ChildReceipt{
				Key: childKey,
				Value: &model.ChildResp{
					Found: found,
					Datum: decryptedStr,
				},
			})
		}

		// generate response
		tokenStruct := &model.DetokenizeResponse{
			ID:   detoken.ID,
			Data: children,
		}
		resp.Resp = tokenStruct
		resp.Code = CodeSuccess

		// set header and return
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
}

func TokenizeHandlerFunc(srv *Service) func(w http.ResponseWriter, r *http.Request) {
	log := srv.log
	return func(w http.ResponseWriter, r *http.Request) {
		log.Logger().Info().Msg(fmt.Sprintf("received a request on %s", Tokenize))
		var resp model.Response
		var token model.Tokenize
		var err error

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			resp.Error = append(resp.Error, ErrMethodNotAllowed)
			log.Logger().Error().Msg(ErrMethodNotAllowed)
			resp.Code = CodeMethodNotAllowed
			json.NewEncoder(w).Encode(resp)
			return
		}

		//reqCtx := context.Background()

		w.Header().Set("Content-Type", "application/json")
		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()
		defer r.Body.Close()

		// Check that json is a valid model.Tokenize structure
		if err = jsonDecoder.Decode(&token); err != nil {
			resp.Error = append(resp.Error, err.Error())
			log.Logger().Error().Msg(err.Error())
			resp.Code = CodeInvalidRequest
			// return 400 status codes
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		var tokenStr string
		var children []model.Child

		// tokenize logic
		manager := srv.manager

		// ensure user request parameter is correct and valid
		validationResp, ok := manager.Validate(&token)
		if !ok {
			for i := 0; i < len(validationResp); i++ {
				resp.Error = append(resp.Error, fmt.Sprintf("error with key %s: %s", validationResp[i].Key, validationResp[i].Err))
				log.Logger().Error().Msg(fmt.Sprintf("error with key %s: %s", validationResp[i].Key, validationResp[i].Err))
			}
			resp.Code = CodeInvalidRequest
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// user request valid, not proceed to process
		parentKey := token.ID
		for i := 0; i < len(token.Data); i++ {
			childKey := token.Data[i].Key
			combinedKeyName := tokenize.GetCombinedKey(parentKey, childKey)
			tokenStr, err = manager.Tokenize(combinedKeyName, token.Data[i].Value)
			if err != nil {
				resp.Error = append(resp.Error, fmt.Sprintf("error with key %s.%s: %s", parentKey, childKey, err.Error()))
				log.Logger().Error().Msg(fmt.Sprintf("error with key %s.%s: %s", parentKey, childKey, err.Error()))
				resp.Code = CodeInternalServerError
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}
			children = append(children, model.Child{
				Key:   childKey,
				Value: tokenStr,
			})
		}

		// generate response
		tokenStruct := &model.TokenizeResponse{
			ID:   token.ID,
			Data: children,
		}
		resp.Resp = tokenStruct
		resp.Code = CodeSuccess

		// set header and return
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}
}
