# URL Shortener with GUI

This is a simple URL shortener application built in Go, featuring both a **REST API backend** and a **desktop GUI application**. The application allows users to shorten long URLs into short, easily shareable links, and also provides a graphical interface for user interaction.

---

## Features

- **REST API Backend**: Handles URL shortening and redirection using the Gin framework.
- **Desktop GUI**: Provides an easy-to-use graphical interface for shortening URLs, built using the Fyne framework.
- **Secure and Random**: Shortened URLs are generated using secure random values.
- **Minimalistic Design**: The GUI is clean and intuitive.

---

## Requirements

### Prerequisites
- Go 1.20+ installed on your system.
- `Gin` and `Fyne` Go packages (installation instructions below).

---

## Installation

### Clone the Repository
```bash
git clone <repository-url>
cd url-shortener
```

```
go get -u github.com/gin-gonic/gin
go get -u fyne.io/fyne/v2
```

### Running the Application

Start the Backend Server

The backend handles URL shortening and redirection. Run the server.go file:

```
go run server.go
```

The server will start on http://localhost:8080.

### Run the GUI Application

The GUI allows users to interact with the URL shortener via a desktop application. Run the gui.go file:

```
go run gui.go
```
This will open the GUI window.

### Usage
Backend API

The backend exposes two main endpoints:

Shorten a URL

- Endpoint: POST /shorten

Payload:
```json
{
  "url": "https://example.com"
}
```
Response:
```json
{
  "short_url": "http://localhost:8080/abc123"
}
```
Redirect to Original URL

- Endpoint: GET /:short

- Access a shortened URL like http://localhost:8080/abc123 to be redirected to the original URL.

### GUI
1) Enter a long URL into the input field.

2) Click the "Shorten URL" button.

3) The shortened URL will appear below the button. You can copy it and share it.

### Directory Structure
```bash
url-shortener/
├── server.go        # Backend server logic
├── gui.go           # GUI application logic
└── README.md        # Project documentation
```
Example Workflow
Commands to Run the Application

Start the backend server:
```
go run server.go
```

Output:

```
[GIN-debug] Listening and serving HTTP on :8080
```

Start the GUI application:

```
go run gui.go
```

Use the GUI to shorten URLs and retrieve the shortened links.

Test the shortened URL in your browser:

http://localhost:8080/<short_id>
