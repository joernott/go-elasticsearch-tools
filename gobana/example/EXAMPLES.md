# Example usages

# Retrieve an aggregation with commandline query 
Retrieve the aggregation only, the details are irrelevant (size=0)
```
./gobana.exe -l1 -e '_search?size=0' -q '{"aggs":{"max_lpt":{"max":{"field":"log_processing_time"}}}}' -A max_lpt
```

# Retrieve an aggregation from file
```
./gobana.exe -l1 -e '_search?size=0' -Q 'example\agg.json' -A max_lpt
```

# Use toml template
```
./gobana.exe -l1 -e '_search?size=0' -q '{"query":{"range":{"@timestamp":{"gte":"{{ .Von }}","lt":"{{ .Bis }}"}}},"aggs":{"max_lpt":{"max":{"field":"log_processing_time"}}}}' -A max_lpt -t -d "Von=now-30d/d" -d "Bis=now/d"
```
