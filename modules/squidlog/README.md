# squid_log

This module parses [`Squid`](http://www.squid-cache.org/) caching proxy logs.

## Charts

Module produces following charts:

-   Total Requests in `requests/s`
-   Excluded Requests in `requests/s`
-   Requests By Type in `requests/s`
-   Responses By HTTP Status Code Class in `responses/s`
-   Responses By HTTP Status Code in `responses/s`
-   Bandwidth in `kilobits/s`
-   Response Time in `milliseconds/s`
-   Unique Clients in `clients/s`
-   Requests By Cache Result Code in `requests/s`
-   Requests By Cache Result Delivery Transport Tag in `requests/s`
-   Requests By Cache Result Handling Tag in `requests/s`
-   Requests By Cache Result Produced Object Tag in `requests/s`
-   Requests By Cache Result Load Source Tag in `requests/s`
-   Requests By Cache Result Errors Tag in `requests/s`
-   Requests By HTTP Method in `requests/s`
-   Requests By MIME Type in `requests/s`
-   Requests By Hierarchy Code in `requests/s`
-   Forwarded Requests By Server Address in `requests/s`

## Log Parsers

Squidlog supports 3 different log parsers:

-   CSV
-   [LTSV](http://ltsv.org/)
-   RegExp

RegExp is the slowest among them but it is very likely you will need to use it if your log format is not default.

## Known Fields

These are [Squid](http://www.squid-cache.org/Doc/config/logformat/) log format codes.

Squidlog is aware how to parse and interpret following codes:

| field                   | squid format code | description                                                            |
|-------------------------|-------------------|------------------------------------------------------------------------|
| resp_time               | %tr               | Response time (milliseconds).
| client_address          | %>a               | Client source IP address.
| client_address          | %>A               | Client FQDN.
| cache_code              | %Ss               | Squid request status (TCP_MISS etc).
| http_code               | %>Hs              | The HTTP response status code from Content Gateway to client.
| resp_size               | %<st              | Total size of reply sent to client (after adaptation).
| req_method              | %rm               | Request method (GET/POST etc).
| hier_code               | %Sh               | Squid hierarchy status (DEFAULT_PARENT etc).
| server_address          | %<a               | Server IP address of the last server or peer connection.
| server_address          | %<A               | Server FQDN or peer name.
| mime_type               | %mt               | MIME content type.

In addition, to make `Squid` [native log format](https://wiki.squid-cache.org/Features/LogFormat#Squid_native_access.log_format_in_detail) csv parsable,
Squidlog understands these groups of codes:

| field                   | squid format code | description                                                            |
|-------------------------|-------------------|------------------------------------------------------------------------|
| result_code             | %Ss/%>Hs          | Cache code and http code.
| hierarchy               | %Sh/%<a           | Hierarchy code and server address.


## Custom Log Format

Custom log format is easy. Use [known fields](#known-fields) to construct your log format.

-   If using CSV parser

```yaml
jobs:
  - name: squid_log_custom_csv_exampla
    path: /var/log/squid/access.log
    log_type: csv
    csv_config:
      format: '- resp_time client_address result_code resp_size req_method - - hierarchy mime_type'
```

Copy your current log format. Replace all known squid format codes with appropriate [known](#known-fields) fields. Replaces others with "-".

-   If using LTSV parser

Provide fields mapping. You need to map your label names to [known](#known-fields) fields.

```yaml
  - name: squid_log_custom_ltsv_exampla
    path: /var/log/squid/access.log
    log_type: ltsv
    ltsv_config:
      mapping:
        label1: resp_time
        label2: client_address
        ...
```

-   If using RegExp parser

Use pattern with subexpressions names. These names should be [known](#known-fields) by squidlog.

```yaml
jobs:
  - name: squid_log_custom_regexp_exampla
    path: /var/log/squid/access.log
    log_type: regexp
    regexp_config:
      format: '^[0-9.]+\s+(?P<resp_time>[0-9]+) (?P<client_address>[\da-f.:]+) (?P<cache_code>[A-Z_]+)/(?P<http_code>[0-9]+) (?P<resp_size>[0-9]+) (?P<req_method>[A-Z]+) [^ ]+ [^ ]+ (?P<hier_code>[A-Z_]+)/[\da-z.:-]+ (?P<mime_type>[A-Za-z-]+)'
```

## Configuration

This module needs only `path` to log file if you use [native log format](https://wiki.squid-cache.org/Features/LogFormat#Squid_native_access.log_format_in_detail).
If you use custom log format you need [to set it manually](#custom-log-format). 

```yaml
jobs:
  - name: squid
    path: /var/log/squid/access.log
    log_type: csv
    csv_config
      format: '- - %h - - %t \"%r\" %>s %b'
```
 
For all available options, please see the module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/squid_log.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m squid_log