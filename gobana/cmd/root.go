package cmd

import (
	"fmt"
	"os"

	"bufio"
	"path/filepath"

	_ "github.com/davecgh/go-spew/spew"
	"github.com/joernott/elasticsearch-tools/gobana/handler"
	"github.com/joernott/lra"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Connection *lra.Connection

var rootCmd = &cobra.Command{
	Use:   "gobana",
	Short: "Gobana is a commandline kibana",
	Long:  `A commandline kibana written in go`,
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		err := HandleConfigFile()
		if err != nil {
			panic(err)
		}
		err = InitLogging()
		if err != nil {
			fmt.Println("Error configuring logging")
			os.Exit(10)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		var result *handler.ElasticsearchResult

		g, err := handler.NewGobana(
			viper.GetBool("ssl"),
			viper.GetString("host"),
			viper.GetInt("port"),
			viper.GetString("user"),
			viper.GetString("password"),
			viper.GetBool("validatessl"),
			viper.GetString("proxy"),
			viper.GetBool("socks"),
			viper.GetString("query"),
			viper.GetString("queryfile"),
			viper.GetBool("toml"),
			viper.GetStringSlice("data"),
			viper.GetString("endpoint"))
		if err != nil {
			os.Exit(20)
		}
		result, err = g.Execute()
		if err != nil {
			os.Exit(21)
		}
		jsonFile := viper.GetString("jsonoutput")
		if jsonFile != "" {
			err := result.WriteFile(jsonFile)
			if err != nil {
				os.Exit(22)
			}
		}
		fieldName := viper.GetString("singlevalue")
		if fieldName != "" {
			result.SingleValue(fieldName)
		}
		fieldName = viper.GetString("aggregation")
		if fieldName != "" {
			result.GetAggregation(fieldName)
		}
	},
}

var ConfigFile string
var UseSSL bool
var ValidateSSL bool
var Host string
var Port int
var User string
var Password string
var LogLevel int
var LogFile string
var Proxy string
var ProxyIsSocks bool
var Query string
var QueryFile string
var Toml bool
var Data []string
var Endpoint string
var JsonOutputFile string
var SingleValue string
var Aggregation string
var ValueOnly bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	pwd, _ := os.Getwd()
	cfgpath := pwd + string(os.PathSeparator) + "gobana.yml"
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", cfgpath, "Configuration file")
	rootCmd.PersistentFlags().BoolVarP(&UseSSL, "ssl", "s", false, "Use SSL")
	rootCmd.PersistentFlags().BoolVarP(&ValidateSSL, "validatessl", "v", true, "Validate SSL certificate")
	rootCmd.PersistentFlags().StringVarP(&Host, "host", "H", "localhost", "Hostname of the server")
	rootCmd.PersistentFlags().IntVarP(&Port, "port", "P", 9200, "Network port")
	rootCmd.PersistentFlags().StringVarP(&User, "user", "u", "", "Username for Elasticsearch")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "Password for the Elasticsearch user")
	rootCmd.PersistentFlags().IntVarP(&LogLevel, "loglevel", "l", 5, "Log level")
	rootCmd.PersistentFlags().StringVarP(&LogFile, "logfile", "L", "", "Log file (defaults to stdout)")
	rootCmd.PersistentFlags().StringVarP(&Proxy, "proxy", "y", "", "Proxy (defaults to none)")
	rootCmd.PersistentFlags().BoolVarP(&ProxyIsSocks, "socks", "Y", false, "This is a SOCKS proxy")
	rootCmd.PersistentFlags().StringVarP(&Query, "query", "q", "", "Query to pass along")
	rootCmd.PersistentFlags().StringVarP(&QueryFile, "queryfile", "Q", "", "File containing a query")
	rootCmd.PersistentFlags().BoolVarP(&Toml, "toml", "t", false, "Use TOML template parsing on query file")
	rootCmd.PersistentFlags().StringSliceVarP(&Data, "data", "d", []string{}, "Data to pass to template parsing, use key=value.")
	rootCmd.PersistentFlags().StringVarP(&Endpoint, "endpoint", "e", "_search", "API endpoint")
	rootCmd.PersistentFlags().StringVarP(&JsonOutputFile, "jsonoutput", "J", "", "Output the result json into this file")
	rootCmd.PersistentFlags().StringVarP(&SingleValue, "singlevalue", "S", "", "Output one single result value from the hits")
	rootCmd.PersistentFlags().StringVarP(&Aggregation, "aggregation", "A", "", "Output one aggregation value")
	rootCmd.PersistentFlags().BoolVarP(&ValueOnly, "valueonly", "V", false, "Output only the value")

	viper.SetDefault("ssl", false)
	viper.SetDefault("validatessl", true)
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 9200)
	viper.SetDefault("user", "")
	viper.SetDefault("password", "")
	viper.SetDefault("loglevel", 5)
	viper.SetDefault("logfile", "")
	viper.SetDefault("proxy", "")
	viper.SetDefault("socks", false)
	viper.SetDefault("query", "")
	viper.SetDefault("queryfile", "")
	viper.SetDefault("toml", false)
	viper.SetDefault("data", []string{})
	viper.SetDefault("endpoint", "_search")
	viper.SetDefault("jsonoutput", "")
	viper.SetDefault("singlevalue", "")
	viper.SetDefault("aggregation", "")
	viper.SetDefault("valueonly", false)

	viper.BindPFlag("ssl", rootCmd.PersistentFlags().Lookup("ssl"))
	viper.BindPFlag("validatessl", rootCmd.PersistentFlags().Lookup("validatessl"))
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	viper.BindPFlag("proxy", rootCmd.PersistentFlags().Lookup("proxy"))
	viper.BindPFlag("socks", rootCmd.PersistentFlags().Lookup("socks"))
	viper.BindPFlag("query", rootCmd.PersistentFlags().Lookup("query"))
	viper.BindPFlag("queryfile", rootCmd.PersistentFlags().Lookup("queryfile"))
	viper.BindPFlag("toml", rootCmd.PersistentFlags().Lookup("toml"))
	viper.BindPFlag("data", rootCmd.PersistentFlags().Lookup("data"))
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))
	viper.BindPFlag("jsonoutput", rootCmd.PersistentFlags().Lookup("jsonoutput"))
	viper.BindPFlag("singlevalue", rootCmd.PersistentFlags().Lookup("singlevalue"))
	viper.BindPFlag("aggregation", rootCmd.PersistentFlags().Lookup("aggregation"))
	viper.BindPFlag("valueonly", rootCmd.PersistentFlags().Lookup("valueonly"))
}

func HandleConfigFile() error {
	if ConfigFile != "" {
		log.Debug("Read config from " + ConfigFile)
		viper.SetConfigFile(ConfigFile)
	} else {
		log.Debug("Read config from home directory")
		home, err := homedir.Dir()
		if err != nil {
			log.Error(err)
			return err
		}
		viper.AddConfigPath(home)
		ex, err := os.Executable()
		if err != nil {
			log.Error(err)
			return err
		}
		pwd := filepath.Dir(ex)
		viper.AddConfigPath(pwd)
		viper.SetConfigName("gobana")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Can't read config: " + err.Error())
		return err
	}

	return nil
}

func InitLogging() error {
	LogFile = viper.GetString("logfile")
	LogLevel = viper.GetInt("loglevel")
	if LogFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		f, err := os.Create(LogFile)
		if err != nil {
			fmt.Println("Could not create logfile '" + LogFile + "'")
			return err
		}
		w := bufio.NewWriter(f)
		log.SetOutput(w)
	}
	switch LogLevel {
	case 0:
		log.SetLevel(log.PanicLevel)
	case 1:
		log.SetLevel(log.FatalLevel)
	case 2:
		log.SetLevel(log.ErrorLevel)
	case 3:
		log.SetLevel(log.WarnLevel)
	case 4:
		log.SetLevel(log.InfoLevel)
	case 5:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
	log.WithFields(log.Fields{
		"LogFile":  LogFile,
		"LogLevel": LogLevel,
	}).Debug("Logging configured")
	log.Debug("PersistentPreRun finished")
	return nil
}
