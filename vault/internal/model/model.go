package model

type Tokenize struct {
	ID   string `json:"id"`
	Data []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"data"`
}

type TokenizeResponse struct {
	ID   string `json:"id"`
	Data []struct {
		Key   string `json:"key"`
		Token string `json:"token"`
	} `json:"data"`
}

type Detokenize struct {
	ID   string `json:"id"`
	Data []struct {
		Key   string `json:"key"`
		Token string `json:"token"`
	} `json:"data"`
}

type DetokenizeResponse struct {
	ID   string `json:"id"`
	Data []struct {
		Key   string `json:"key"`
		Value struct {
			Found bool   `json:"found"`
			Datum string `json:"satum"`
		} `json:"value"`
	} `json:"data"`
}

type Resp interface {
}

type Response struct {
	Resp  `json:"resp"`
	Code  int    `json:"code"`
	Error string `json:"error"`
}
