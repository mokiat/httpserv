# httpserv

I find myself needing a simple HTTP file server quite often and `python -m SimpleHTTPServer` just doesn't cut it for me anymore, since it always runs on `0.0.0.0`.

The `httpserv` tool takes things just one idea further, allowing the most basic aspects (host, port, directory) to be configured.

Use the following command to download the server (requires a working Go environment):

```sh
GO111MODULE=on go get github.com/mokiat/httpserv
```
