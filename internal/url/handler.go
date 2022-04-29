package url

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"math/rand"
	"net/http"
	"time"
	"url-shortener-alt/internal/handlers"
	"url-shortener-alt/pkg/utils"
)

const (
	writeUrl = "/a/"
	readUrl  = "/s/:code"
)

type handler struct {
	repository Repository
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, writeUrl, h.EncodeLink)
	router.HandlerFunc(http.MethodGet, readUrl, h.DecodeLink)
}

func (h *handler) EncodeLink(writer http.ResponseWriter, request *http.Request) {
	originalUrl := request.URL.Query().Get("url")

	findByOriginal, err := h.repository.FindByOriginal(context.TODO(), originalUrl)
	if err != nil && findByOriginal.ID == "nil" {
		log.Println(err)
		writer.WriteHeader(http.StatusTeapot)
		return
	}
	if err == nil {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Add("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"hash": findByOriginal.HashUrl})
		return
	}

	hashUrl := utils.RandomString(8)
	for h.repository.CheckHash(context.TODO(), hashUrl) {
		hashUrl = utils.RandomString(8)
	}

	url := Url{
		OriginalUrl: originalUrl,
		HashUrl:     hashUrl,
	}

	_, err = h.repository.Create(context.TODO(), url)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"hash": url.HashUrl})
}

func (h *handler) DecodeLink(writer http.ResponseWriter, request *http.Request) {
	hashUrl := httprouter.ParamsFromContext(request.Context()).ByName("code")
	url, err := h.repository.FindByHash(context.TODO(), hashUrl)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(writer, request, url.OriginalUrl, http.StatusFound)
}

func NewHandler(repository Repository) handlers.Handler {
	rand.Seed(time.Now().UnixNano())
	return &handler{repository: repository}
}
