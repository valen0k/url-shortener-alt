package url

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
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
	router.HandlerFunc(http.MethodPost, writeUrl, h.EncodeLink)
	router.HandlerFunc(http.MethodGet, readUrl, h.DecodeLink)
}

func (h *handler) EncodeLink(writer http.ResponseWriter, request *http.Request) {
	originalUrl := request.URL.Query().Get("url")
	hash := utils.GenerateHash(originalUrl)
	findByHash, err := h.repository.FindByHash(context.TODO(), hash)
	if err != nil {
		_, err := h.repository.Create(context.TODO(), originalUrl, hash, 0)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	fmt.Println(findByHash)

	writer.WriteHeader(http.StatusOK)
	writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"hash": hash})
}

func (h *handler) DecodeLink(writer http.ResponseWriter, request *http.Request) {
	hashUrl := httprouter.ParamsFromContext(request.Context()).ByName("code")
	url, err := h.repository.FindByHash(context.TODO(), hashUrl)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(writer, request, url, http.StatusFound)
}

func NewHandler(repository Repository) handlers.Handler {
	return &handler{repository: repository}
}
