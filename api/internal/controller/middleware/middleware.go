package middleware

import (
	"encoding/json"
	"net/http"
	"quic_upload/api/internal/view"
	"quic_upload/api/pkg/logging"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	logging.Init()
	logger := logging.GetLogger()
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			re, ok := err.(*requestError)
			if ok {
				logger.Errorf("code [%d] path [%s]: %v", re.StatusCode, r.RequestURI, err)
				resp := view.Response{
					Code:  responseCodeToString(re.StatusCode),
					Error: re.Error(),
				}
				w.WriteHeader(re.StatusCode)
				respJson, _ := json.Marshal(resp)
				w.Write(respJson)
				return
			}

			logger.Errorf("code [%d] path [%s]: %v", http.StatusInternalServerError, r.RequestURI, err)
			resp := view.Response{
				Code:  responseCodeToString(http.StatusInternalServerError),
				Error: re.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			respJson, _ := json.Marshal(resp)
			w.Write(respJson)
			return
		}
	}
}
