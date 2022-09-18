package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	extract "github.com/golang-jwt/jwt/v4"
	"llil.gq/go/database"
	"log"
	"net/http"
	"net/url"
)

type DataShortenHandler struct {
	database     database.Database
	jwtValidator *validator.Validator
	baseUrl      string
}

func FormatResponse(baseUrl string, shortURL string) []byte {
	response := make(map[string]string)
	shortURLResponse := fmt.Sprintf("%s/%s", baseUrl, shortURL)
	response["shortURL"] = shortURLResponse
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	return jsonResponse
}

func (h *DataShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	auth0Token := r.Header.Get("Authorization")
	auth0Sub := ""
	if auth0Token != "" {
		_, err := h.jwtValidator.ValidateToken(context.Background(), auth0Token)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		token, _, err := new(extract.Parser).ParseUnverified(auth0Token, extract.MapClaims{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if claims, ok := token.Claims.(extract.MapClaims); ok {
			auth0Sub = fmt.Sprint(claims["sub"])
		}
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var data longUrlPayload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = url.ParseRequestURI(data.LongURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data.Auth0Sub = auth0Sub
	shortURL := h.generateShortUrl(data)

	w.WriteHeader(http.StatusCreated)
	jsonResponse := FormatResponse(h.baseUrl, shortURL)
	w.Write(jsonResponse)
}

func (h *DataShortenHandler) generateShortUrl(data longUrlPayload) string {
	shortURL := computeShortURL(data.Auth0Sub + data.LongURL)
	result, err := database.SelectShortURL(h.database, shortURL)
	if err != nil {
		database.AddShortUrl(h.database, data.LongURL, shortURL, data.Auth0Sub)
	} else {
		if result.LongURL != data.LongURL {
			for err == nil {
				shortURL = computeShortURL(h.baseUrl + data.Auth0Sub + data.LongURL)
				result, err = database.SelectShortURL(h.database, shortURL)
			}
			database.AddShortUrl(h.database, data.LongURL, shortURL, data.Auth0Sub)
		}
	}
	return shortURL
}
