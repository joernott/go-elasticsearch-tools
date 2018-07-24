package cmd

import (
	"fmt"
	"os"

	"bufio"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/joernott/elasticsearch-tools/elasticsearch"
	"github.com/joernott/elasticsearch-tools/elastui/server"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Connection *elasticsearch.ElasticsearchConnection

var rootCmd = &cobra.Command{
	Use:   "elasticui",
	Short: "ElasticUI manages elasticsearch indices",
	Long:  `A web interface for managing indices in elasticsearch written in go/javascript`,
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		if LogFile == "" {
			log.SetOutput(os.Stdout)
		} else {
			f, err := os.Create(LogFile)
			if err != nil {
				fmt.Println("Could not create logfile '" + LogFile + "'")
				os.Exit(10)
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
		spew.Dump(LogLevel)
		log.WithFields(log.Fields{
			"LogFile":  LogFile,
			"LogLevel": LogLevel,
		}).Debug("Logging configured")

		if ConfigFile != "" {
			log.Debug("Read config from " + ConfigFile)
			viper.SetConfigFile(ConfigFile)
		} else {
			log.Debug("Read config from home directory")
			home, err := homedir.Dir()
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			viper.AddConfigPath(home)
			ex, err := os.Executable()
			if err != nil {
				log.Error(err)
				panic(err)
			}
			pwd := filepath.Dir(ex)
			viper.AddConfigPath(pwd)
			viper.SetConfigName("elastui")
		}

		if err := viper.ReadInConfig(); err != nil {
			log.Error("Can't read config" + err.Error())
			os.Exit(1)
		}
		log.Debug("PersistentPreRun finished")
	},
	Run: func(cmd *cobra.Command, args []string) {
		Connection = elasticsearch.NewElasticsearch(UseSSL, Host, Port, User, Password, ValidateSSL, Proxy, ProxyIsSocks)
		server.Router(Connection)
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", "./elastui.yml", "config file (default is ./elastui.yml)")
	rootCmd.PersistentFlags().BoolVarP(&UseSSL, "ssl", "s", false, "Use SSL (defaults to no)")
	rootCmd.PersistentFlags().BoolVarP(&ValidateSSL, "validatessl", "v", true, "Validate SSL certificate (defaults to yes)")
	rootCmd.PersistentFlags().StringVarP(&Host, "host", "H", "localhost", "Hostname of the server (defaults to localhost)")
	rootCmd.PersistentFlags().IntVarP(&Port, "port", "P", 9200, "Network port (defaults to 9200)")
	rootCmd.PersistentFlags().StringVarP(&User, "user", "u", "", "user (default is empty)")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "password for the elasticsearch user")
	rootCmd.PersistentFlags().IntVarP(&LogLevel, "loglevel", "l", 5, "log level (defaults to 4 (Info))")
	rootCmd.PersistentFlags().StringVarP(&LogFile, "logfile", "L", "", "logfile (defaults to stdout)")
	rootCmd.PersistentFlags().StringVarP(&Proxy, "proxy", "y", "", "proxy (defaults to none)")
	rootCmd.PersistentFlags().BoolVarP(&ProxyIsSocks, "socks", "Y", false, "This is a SOCKS proxy (defaults to no)")

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
}
