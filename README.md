# httpserv

I find myself needing a simple HTTP file server quite often and `python -m SimpleHTTPServer` just doesn't cut it for me anymore, since it always runs on `0.0.0.0` and the actual command differs with the Python versions.

The `httpserv` aims to solve this problem, especially for Go developers. You can easily install it using the following shell command:

```sh
go install github.com/mokiat/httpserv@latest
```

Use this next command to get available configuration flags:

```sh
httpserv --help
```

To run an HTTP server that serves the files and folders located in the current directory, just run:

```sh
httpserv
```

If you want to use `httpserv` but you already have an executable called `httpserv` on your system or you just don't want to install it, you can run it directly:

```sh
go run github.com/mokiat/httpserv@latest
```
