package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/joernott/lra"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Gobana struct {
	Connection *lra.Connection
	Endpoint   string
	Query      string
}

type ElasticsearchResult struct {
	Took         int                          `json:"took"`
	TimedOut     bool                         `json:"timed_out"`
	Shards       ElasticsearchShardResult     `json:"_shards"`
	Hits         ElasticsearchHitResult       `json:"hits"`
	Error        string                       `json:"error"`
	Status       int                          `json:"status`
	Aggregations map[string]AggregationResult `json:"aggregations"`
}

type ElasticsearchShardResult struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type ElasticsearchHitResult struct {
	Total    int64                  `json:"total"`
	MaxScore float64                `json:"max_score"`
	Hits     []ElasticsearchHitList `json:"hits"`
}

type ElasticsearchHitList struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	Id     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

type AggregationResult map[string]interface{}

func NewGobana(UseSSL bool, Server string, Port int, User string, Password string, ValidateSSL bool, Proxy string, ProxyIsSocks bool, Query string, Queryfile string, Toml bool, Data []string, Endpoint string) (*Gobana, error) {
	var g *Gobana
	var err error

	logger := log.WithField("func", "NewGobana")
	g = new(Gobana)

	hdr := make(lra.HeaderList)
	hdr["Content-Type"] = "application/json"
	Connection, err := lra.NewConnection(UseSSL,
		Server,
		Port,
		User,
		Password,
		"",
		ValidateSSL,
		Proxy,
		ProxyIsSocks,
		hdr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	g.Connection = Connection
	g.Endpoint = Endpoint
	if Queryfile != "" {
		if Query != "" {
			logger.WithFields(log.Fields{
				"Query":     Query,
				"QueryFile": Queryfile,
			}).Error("Can't use exclusive parameters query and queryfile at the same time")
			err := errors.New("Can't use exclusive parameters query and queryfile at the same time")
			return nil, err
		}
		Query, err = getQueryFromQueryfile(Queryfile)
		if err != nil {
			return nil, err
		}
	}
	if Toml {
		Query, err = parseToml(Query, Data)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}
	g.Query = Query
	logger.Debug(g.Query)
	return g, nil
}

func getQueryFromQueryfile(Queryfile string) (string, error) {
	logger := log.WithField("func", "getQueryFromQueryfile")
	logger.Debug("Reading query from file '" + Queryfile + "'.")
	Queryfile, err := filepath.Abs(Queryfile)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	logger.WithField("Path", Queryfile).Debug("Expand path")
	buf, err := ioutil.ReadFile(Queryfile)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	logger.Info("Read query from'" + Queryfile + "'.")
	return string(buf), nil
}

func parseToml(Query string, Data []string) (string, error) {
	var Fields map[string]string
	logger := log.WithField("func", "parseToml")

	Fields = make(map[string]string)
	for _, d := range Data {
		x := strings.SplitN(d, "=", 2)
		Fields[x[0]] = x[1]
	}
	logger.Debug("Parsing TOML")
	t := template.New("query")
	t, err := t.Parse(Query)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, Fields)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return buf.String(), nil
}

func (gobana *Gobana) Execute() (*ElasticsearchResult, error) {
	var ResultJson *ElasticsearchResult

	ResultJson = new(ElasticsearchResult)

	logger := log.WithField("func", "Gobana.Execute")
	logger.WithFields(log.Fields{"query": gobana.Query, "endpoint": gobana.Endpoint}).Debug("Execute")
	result, err := gobana.Connection.Post("/"+gobana.Endpoint, []byte(gobana.Query))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info("Successfully executed query")
	err = json.Unmarshal(result, ResultJson)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	ll := viper.GetInt("loglevel")
	if ll == 5 {
		debugOut, err := json.MarshalIndent(ResultJson, "", "  ")
		if err != nil {
			logger.Error(err)
		}
		log.Debug(string(debugOut))
		spew.Dump(ResultJson)
	}
	return ResultJson, nil
}

func (result *ElasticsearchResult) WriteFile(FileName string) error {
	logger := log.WithField("func", "ElasticsearchResult.WriteFile")
	output, err := json.Marshal(result)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = ioutil.WriteFile(FileName, output, 0644)
	return nil
}

func (result *ElasticsearchResult) SingleValue(FieldName string) {
	var found int64
	var ok bool

	logger := log.WithField("func", "ElasticsearchResult.SingleValue")
	logger.Info("Collecting single values")
	for _, hit := range result.Hits.Hits {
		s := hit.Source
		if result, ok := s[FieldName]; ok {
			if !viper.GetBool("valueonly") {
				fmt.Print(hit.Id + ":")
			}
			fmt.Println(result)
			found++
		}
		log.WithFields(log.Fields{
			"Index":    hit.Index,
			"Type":     hit.Type,
			"HasField": ok,
			"Value":    result}).Debug("Processed " + hit.Id)
	}
	logger.WithField("Found", found).Info("Finished collecting single values")
}

func (result *ElasticsearchResult) GetAggregation(FieldName string) {
	logger := log.WithField("func", "ElasticsearchResult.GetAggregation")
	logger.Info("Collecting aggregation")
	if a, ok := result.Aggregations[FieldName]; ok {
		for k, v := range a {
			if !viper.GetBool("valueonly") {
				fmt.Printf("%v:", k)
			}
			fmt.Printf("%v\n", v)
			log.WithFields(log.Fields{
				"Key":   k,
				"Value": v}).Debug("Processed Aggregation " + FieldName)
		}

	} else {
		logger.WithField("aggregation", FieldName).Warn("Aggregation not found")
	}

}
