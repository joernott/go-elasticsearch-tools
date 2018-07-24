# Elasticsearch-tools Tools for working with elasticsearch

## Purpose
This repository collects tools interacting with the elasticsearch API.
It is WIP and far from complete. Main features are the ability to connect
to elasticseach through a proxy and handle self signed certificates gracefully.

Feel free to contribute and improve.

## License
BSD 3-Clause

## Tools
### Elasticui
A local webserver backend which connects to an elasticsearch server via proxy.
It can show and delete indices and execute API calls.

#### Usage
  elasticui [flags]

##### Flags:
|Short | Long          | Type  | Purpose                                       |
|------|---------------|-------|-----------------------------------------------|
| -h   | --help        |       | help for gobana                               |
| -c   | --config      |string | Configuration file (default elastui.yml in the current working directory)|
| -u   | --user        |string | Username for Elasticsearch                    |
| -p   | --password    |string | Password for the Elasticsearch user           |
| -s   | --ssl         |bool   | Use SSL                                       |
| -v   | --validatessl |bool   | Validate SSL certificate (default true)       |
| -H   | --host        |string | Hostname of the server (default "localhost")  |
| -P   | --port        |int    | Network port (default 9200)                   |
| -y   | --proxy       |string | Proxy (defaults to none)                      |
| -Y   | --socks       |bool   | This is a SOCKS proxy                         |
| -l   | --loglevel    |int    | Log level (default 5)                         |
| -L   | --logfile     |string | Log file (defaults to stdout)                 |


### Gobana
A commandline tool to query elasticsearch. It is a commandline 'Kibana' to
retrieve data form elasticsearch.

Using cobra/viper, Gobana uses a configuration file gobana.yml and allows for
commandline parameters.

#### Usage
  gobana [flags]

##### Flags:
|Short | Long          | Type  | Purpose                                       |
|------|---------------|-------|-----------------------------------------------|
| -h   | --help        |       | help for gobana                               |
| -c   | --config      |string | Configuration file (default gobana.yml in the current working directory)|
| -u   | --user        |string | Username for Elasticsearch                    |
| -p   | --password    |string | Password for the Elasticsearch user           |
| -s   | --ssl         |bool   | Use SSL                                       |
| -v   | --validatessl |bool   | Validate SSL certificate (default true)       |
| -H   | --host        |string | Hostname of the server (default "localhost")  |
| -P   | --port        |int    | Network port (default 9200)                   |
| -y   | --proxy       |string | Proxy (defaults to none)                      |
| -Y   | --socks       |bool   | This is a SOCKS proxy                         |
| -l   | --loglevel    |int    | Log level (default 5)                         |
| -L   | --logfile     |string | Log file (defaults to stdout)                 |
| -e   | --endpoint    |string | API endpoint (default "_search")              |
| -q   | --query       |string | Query to pass along                           |
| -Q   | --queryfile   |string | File containing a query                       |
| -t   | --toml        |bool   | Use TOML template parsing on query file       |
| -d   | --data        |strings| Pass fields to template parsing, use key=value and use it in the template with {{ .Key }}, this flag can be used multiple times|
| -J   | --jsonoutput  |string | Output the result json into this file         |
| -S   | --singlevalue |string | Output one single result value from the hits  |
| -A   | --aggregation |string | Output one aggregation value                  |
| -V   | --valueonly   |bool   | Output only the value                         |

