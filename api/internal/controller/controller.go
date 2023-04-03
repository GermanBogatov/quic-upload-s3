package controller

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"quic_upload/api/internal/config"
	"quic_upload/api/internal/controller/middleware"
	"quic_upload/api/internal/service"
)

type Controller struct {
	service service.IService
	cfg     *config.Server
}

func NewController(service service.IService) *Controller {
	return &Controller{service: service}
}

type tests struct {
	Name    string
	Surname string
	Age     int
}

func (c *Controller) InitRoutes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/api/v1/test", middleware.Middleware(c.Test))
	router.HandlerFunc(http.MethodPost, "/api/v1/create", middleware.Middleware(c.Create))

	return router
}

func (c *Controller) validateXApiKeyHeader(r *http.Request) bool {
	fmt.Println("c.cfg.ApiKey:=", c.cfg.ApiKey)
	fmt.Println("r.Header.Get(X-API-KEY):=", r.Header.Get("X-API-KEY"))
	if c.cfg.ApiKey == r.Header.Get("X-API-KEY") {
		return true
	}
	return false

}

type testStruct struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var res testStruct

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return middleware.NewErr(http.StatusUnprocessableEntity, err)
	}
	fmt.Println("len():=", len(res.Data))
	if len(res.Data) == 0 {
		return middleware.NewErr(http.StatusUnprocessableEntity, errors.New("fail read data"))
	}

	_ = os.WriteFile(res.Name, res.Data, 0644)

	w.WriteHeader(http.StatusOK)
	return nil
}

func (c *Controller) Test(w http.ResponseWriter, r *http.Request) error {
	//if !c.validateXApiKeyHeader(r) {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return errors.New("fail api-key")
	//}

	w.Header().Set("Content-Type", "application/json")
	var me = tests{
		Name:    "German",
		Surname: "Bogatov from uploader",
		Age:     24,
	}
	notesBytes, err := json.Marshal(me)
	if err != nil {
		return middleware.NewErr(http.StatusUnprocessableEntity, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(notesBytes)
	return nil
}
