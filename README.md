# go2web

CLI that fetches web pages and search results using manual HTTP over TCP.

## Usage

```bash
Options:
  -u, -url <URL>               # make an HTTP request to the specified URL and print the response
  -s, -search <search-term>    # make an HTTP request to search the term using your favorite search engine and print top 10 results
  -h, -help                    # show this help
```

## Prerequisites
- Go version `1.25.4`+
- Terminal with a monospace nerd font

## Build

```bash
go build -o go2web ./cmd/go2web
```

Run:
```
./go2web -h
```
## Demo
- URL mode:
<img src="docs/demo-url.gif">

- Search mode:
<img src="docs/demo-search.gif">

- Help:
<img src="docs/demo-help.gif">
