## ghastly raw

send a raw request and print the result

### Synopsis

send a raw request and print the result

```
ghastly raw <path> [flags]
```

### Options

```
  -a, --arg stringArray   arguments to send along with the request, key=value pairs. provide multiple times for multiple argments.
  -h, --help              help for raw
  -w, --websocket         if true, send request to the given websocket endpoint; otherwise, send a GET
```

### Options inherited from parent commands

```
      --loglevel string   log level; one of TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC (default "INFO")
      --server string     the URL used to access homeassistant. defaults to value of HASS_SERVER environment variable
      --token string      the bearer token used to authenticate to homeassistant. defaults to value of HASS_TOKEN environment variable
```

### SEE ALSO

* [ghastly](ghastly.md)	 - ghastly is a tool for interacting with homeassistant

###### Auto generated by spf13/cobra on 21-Apr-2020