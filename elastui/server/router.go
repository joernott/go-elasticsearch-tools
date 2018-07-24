package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joernott/elasticsearch-tools/elasticsearch"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

var Connection *elasticsearch.ElasticsearchConnection

func Router(elasticsearch *elasticsearch.ElasticsearchConnection) {

	Connection = elasticsearch
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/*path", ApiGET)
	router.PUT("/api/*path", ApiPUT)
	router.POST("/api/*path", ApiPOST)
	router.DELETE("/api/*path", ApiDELETE)
	log.Debug("Router initialized")
	router.ServeFiles("/static/*filepath", http.Dir("static/"))
	if _, err := os.Stat("server.crt"); err == nil {
		if _, err := os.Stat("server.key"); os.IsNotExist(err) {
			log.Error("Found server.crt but missing server.key")
			fmt.Println("Found server.crt but missing server.key")
			os.Exit(10)
		}
		log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", router))

	} else {
		log.Fatal(http.ListenAndServe(":8080", router))
	}
}
