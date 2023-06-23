package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	extract "github.com/golang-jwt/jwt/v4"
	"github.com/samber/lo"
	"on3z.cc/go/database"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type DataShortenHandler struct {
	database     database.Database
	jwtValidator *validator.Validator
	baseUrl      string
}

func FormatResponsePost(baseUrl string, shortURL string) []byte {
	response := make(map[string]string)
	shortURLResponse := fmt.Sprintf("%s/%s", baseUrl, shortURL)
	response["shortURL"] = shortURLResponse
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	return jsonResponse
}

func FormatResponseGet() []byte {
	response := make(map[string]string)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	return jsonResponse
}

func (h *DataShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
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

	if !lo.Contains[string]([]string{"POST", "GET", "PUT"}, r.Method) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// ONLY POST request can be anonymous
	if lo.Contains[string]([]string{"GET", "PUT"}, r.Method) && auth0Sub == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method == "POST" {
		var data POSTPayload
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
		jsonResponse := FormatResponsePost(h.baseUrl, shortURL)
		w.Write(jsonResponse)
	}
	if r.Method == "GET" {
		pageLengthStr := r.URL.Query().Get("length")
		offsetStr := r.URL.Query().Get("offset")
		if pageLengthStr == "" || offsetStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			response := make(map[string]string)
			response["error"] = "Missing GET parameters, required: length, offset"
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResponse)
			return
		}
		pageLength, err := strconv.Atoi(pageLengthStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response := h.listShortUrlByUser(auth0Sub, pageLength, offset)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
	if r.Method == "PUT" {
		var data PUTPayload
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
		// Check shortURL already exists in db for this user
		err = h.updateShortUrl(data, auth0Sub)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func (h *DataShortenHandler) generateShortUrl(data POSTPayload) string {
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

func (h *DataShortenHandler) listShortUrlByUser(auth0Sub string, pageLength int, offset int) []byte {
	result, err := database.ListShortURLByUser(h.database, auth0Sub, pageLength, offset)
	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	return response
}

func (h *DataShortenHandler) updateShortUrl(data PUTPayload, auth0Sub string) error {
	_, err := database.SelectShortURLByUserByShortURL(h.database, data.ShortURL, auth0Sub)
	if err != nil {
		return err
	}
	database.UpdateShortUrl(h.database, data.LongURL, data.ShortURL, auth0Sub)
	return nil
}
