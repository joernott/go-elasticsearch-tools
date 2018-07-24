package elasticsearch

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

type ElasticsearchConnection struct {
	Protocol     string
	Server       string
	Port         int
	User         string
	Password     string
	ValidateSSL  bool
	Proxy        string
	ProxyIsSocks bool
	BaseURL      string
	Client       *http.Client
}

func NewElasticsearch(UseSSL bool, Server string, Port int, User string, Password string, ValidateSSL bool, Proxy string, ProxyIsSocks bool) *ElasticsearchConnection {
	var elasticsearch *ElasticsearchConnection
	var tr *http.Transport

	elasticsearch = new(ElasticsearchConnection)
	tr = &http.Transport{
		DisableKeepAlives:   false,
		IdleConnTimeout:     0,
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 100,
	}
	if !ValidateSSL {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if Proxy != "" {
		if ProxyIsSocks {
			dialer, err := proxy.SOCKS5("tcp", Proxy, nil, proxy.Direct)
			if err != nil {
				log.WithField("url", Proxy).Error("Can't connect to Socks5 proxy: " + err.Error())
			}
			tr.Dial = dialer.Dial
		} else {
			proxyUrl, err := url.Parse(Proxy)
			if err != nil {
				log.WithField("url", Proxy).Error(err)
			}
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	elasticsearch.Client = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10}
	if UseSSL {
		elasticsearch.Protocol = "https"
	} else {
		elasticsearch.Protocol = "http"
	}
	elasticsearch.BaseURL = elasticsearch.Protocol + "://"
	elasticsearch.Server = Server
	elasticsearch.Port = Port
	elasticsearch.User = User
	elasticsearch.Password = Password
	elasticsearch.ValidateSSL = ValidateSSL
	elasticsearch.Proxy = Proxy
	elasticsearch.ProxyIsSocks = ProxyIsSocks
	if User != "" {
		elasticsearch.BaseURL = elasticsearch.BaseURL + User + ":" + Password + "@"
	}
	elasticsearch.BaseURL = elasticsearch.BaseURL + Server + ":" + strconv.Itoa(Port)
	log.WithFields(log.Fields{
		"Protocol":     elasticsearch.Protocol,
		"Server":       elasticsearch.Server,
		"Port":         elasticsearch.Port,
		"User":         elasticsearch.User,
		"Password":     elasticsearch.Password,
		"ValidateSSL":  elasticsearch.ValidateSSL,
		"Proxy":        elasticsearch.Proxy,
		"ProxyIsSocks": elasticsearch.ProxyIsSocks,
		"BaseURL":      elasticsearch.BaseURL,
	}).Debug("Elasticsearch connection initialized")
	return elasticsearch
}

func (elasticsearch *ElasticsearchConnection) Get(endpoint string) (map[string]interface{}, error) {
	/*	var data interface{}

		target := elasticsearch.BaseURL + endpoint
		r, err := elasticsearch.Client.Get(target)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		response, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
	*/
	var x []byte

	response, err := elasticsearch.request("GET", endpoint, x)
	if err != nil {
		return nil, err
	}
	return toJSON(response)
}

func (elasticsearch *ElasticsearchConnection) GetRaw(endpoint string) ([]byte, error) {
	/*
		target := elasticsearch.BaseURL + endpoint
		r, err := elasticsearch.Client.Get(target)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		response, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		return response, nil
	*/
	var x []byte
	return elasticsearch.request("GET", endpoint, x)
}

func (elasticsearch *ElasticsearchConnection) Delete(endpoint string) (map[string]interface{}, error) {
	var x []byte

	response, err := elasticsearch.request("DELETE", endpoint, x)
	if err != nil {
		return nil, err
	}
	return toJSON(response)
}

func (elasticsearch *ElasticsearchConnection) DeleteRaw(endpoint string) ([]byte, error) {
	var x []byte
	return elasticsearch.request("DELETE", endpoint, x)
}

func (elasticsearch *ElasticsearchConnection) Post(endpoint string, jsonData []byte) (map[string]interface{}, error) {
	response, err := elasticsearch.request("POST", endpoint, jsonData)
	if err != nil {
		return nil, err
	}
	return toJSON(response)
}

func (elasticsearch *ElasticsearchConnection) PostRaw(endpoint string, jsonData []byte) ([]byte, error) {
	return elasticsearch.request("POST", endpoint, jsonData)
}

func (elasticsearch *ElasticsearchConnection) Put(endpoint string, jsonData []byte) (map[string]interface{}, error) {
	response, err := elasticsearch.request("PUT", endpoint, jsonData)
	if err != nil {
		return nil, err
	}
	return toJSON(response)
}

func (elasticsearch *ElasticsearchConnection) PutRaw(endpoint string, jsonData []byte) ([]byte, error) {
	return elasticsearch.request("PUT", endpoint, jsonData)
}

func (elasticsearch *ElasticsearchConnection) request(method string, endpoint string, jsonData []byte) ([]byte, error) {
	var req *http.Request
	var err error

	target := elasticsearch.BaseURL + endpoint
	switch method {
	case "GET", "DELETE":
		req, err = http.NewRequest(method, target, nil)
	default:
		req, err = http.NewRequest(method, target, bytes.NewBuffer(jsonData))
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	r, err := elasticsearch.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func toJSON(response []byte) (map[string]interface{}, error) {
	var data interface{}
	err := json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}
	m := data.(map[string]interface{})
	return m, nil
}
