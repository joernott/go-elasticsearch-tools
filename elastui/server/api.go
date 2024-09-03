package server

import (
	_ "encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/static/index.html", http.StatusPermanentRedirect)
}

func ApiGET(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	path := params.ByName("path")
	query := r.URL.RawQuery
	logger := log.WithFields(log.Fields{
		"func":    "ApiGET",
		"BaseURL": Connection.BaseURL,
		"Path":    path,
		"Query":   query,
	})

	logger.Debug("GET")
	data, err := Connection.Get(params.ByName("path") + "?" + query)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		//w.Header().Set("Content-Type", "application/json")
		logger.Debug("GET successful")
		w.Write([]byte(data))
	}
}

func ApiPOST(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	path := params.ByName("path")
	query := r.URL.RawQuery
	logger := log.WithFields(log.Fields{
		"func":    "ApiPOST",
		"BaseURL": Connection.BaseURL,
		"Path":    path,
		"Query":   query,
	})

	logger.Debug("POST")
	inputdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	data, err := Connection.Post(params.ByName("path")+"?"+query, inputdata)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		logger.Debug("POST successful")
		w.Write([]byte(data))
	}
}

func ApiPUT(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	path := params.ByName("path")
	query := r.URL.RawQuery
	logger := log.WithFields(log.Fields{
		"func":    "ApiPUT",
		"BaseURL": Connection.BaseURL,
		"Path":    path,
		"Query":   query,
	})

	logger.Debug("PUT")
	inputdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	data, err := Connection.Put(params.ByName("path")+"?"+query, inputdata)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		logger.Debug("PUT successful")
		w.Write([]byte(data))
	}
}

func ApiDELETE(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	path := params.ByName("path")
	query := r.URL.RawQuery
	logger := log.WithFields(log.Fields{
		"func":    "ApiDELETE",
		"BaseURL": Connection.BaseURL,
		"Path":    path,
		"Query":   query,
	})

	logger.Debug("DELETE")
	data, err := Connection.Delete(params.ByName("path")+"?"+query, []byte("{}"))
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		//w.Header().Set("Content-Type", "application/json")
		logger.Debug("DELETE successful")
		w.Write([]byte(data))
	}
}
