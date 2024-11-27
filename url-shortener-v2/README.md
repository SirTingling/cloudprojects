# URL Shortener V2 (Lightweight)

A lightweight URL shortener built in Go using `gorilla/mux`. This application allows shortening URLs and redirecting short URLs to their original destinations.

---

## Features

- Shorten URLs with an MD5 hash-based key.
- Redirect short URLs to their original long URLs.
- Simple in-memory storage.

## Requirements

- Go 1.20+
- `gorilla/mux` package

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd url-shortener

### Installing Dependencies

```
go get -u github.com/gorilla/mux
```

### Running:
Start the application

```
go run main.go
```

The server runs at `http://localhost:8080`.


### Usage

- Send a `POST` request to `/create` with the `url` parameter

```bash
curl -X POST -d "url=https://example.com" http://localhost:8080/create
```

- Response:

```bash
http://localhost:8080/abcde
```

### Redirect to the Original URL

Access the shortened URL:
```bash
curl http://localhost:8080/abcde
```

This redirects to the original URL.

---
