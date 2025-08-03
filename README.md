# Badger Database Web UI

A web interface for managing Badger key-value databases, built with Go, badger, mux, HTMX, and Tailwind.

> [!Note] Tested with badger v4.8.0

---
> [!NOTE]
> A vibe coding project.
> This codebase may prioritize rapid prototyping, creativity, and experimentation over strict conventions, exhaustive documentation, or production-grade robustness.
> Expect some rough edges, unconventional solutions, and evolving features. Use with curiosity and caution!

## 🚀 Features

- **Real-time Database Management**: Add, edit, delete, and search key-value pairs with live updates
- **HTMX Integration**: Smooth, dynamic interactions without full page reloads
- **RESTful API**: Complete REST API for programmatic access
- **Live Statistics**: Monitor database size and key count in real-time
- **Search Functionality**: Find keys quickly with live search
- **Responsive Design**: Works on desktop and mobile devices

---

## 🏁 Quick Start

### Prerequisites

- Go **1.21+**
- Git

### Installation

Clone the repository:

```bash
git clone https://github.com/jesusnoseq/badger-web-ui
cd badger-web-ui
```

Install dependencies:

```bash
make deps
```

### Run the application

```bash
make run
```

Open your browser and navigate to `http://localhost:8080`

---

## 🖥️ Usage

### Web Interface

- **Add Keys**: Use the "Add New Key" form to create new key-value pairs
- **Search**: Type in the search box to find keys in real-time
- **Edit**: Click the "Edit" button next to any key to modify its value
- **Delete**: Click the "Delete" button to remove a key (with confirmation)
- **Statistics**: View live database statistics in the header

### API Endpoints

- `GET /api/keys` - List all keys (with optional `?limit=N` parameter)
- `POST /api/keys` - Create a new key-value pair
- `GET /api/keys/{key}` - Get a specific key's value
- `PUT /api/keys/{key}` - Update a key's value
- `DELETE /api/keys/{key}` - Delete a key
- `GET /api/stats` - Get database statistics
- `GET /api/search?q={query}` - Search for keys

#### Example API Usage

```bash
# Add a new key
curl -X POST http://localhost:8080/api/keys \
  -H "Content-Type: application/json" \
  -d '{"key": "user:123", "value": "John Doe"}'

# Get a key
curl http://localhost:8080/api/keys/user:123

# Search for keys
curl http://localhost:8080/api/search?q=user

# Get statistics
curl http://localhost:8080/api/stats
```

---

## ⚙️ Configuration

### Environment Variables

You can customize the application's behavior using environment variables:

### Available Environment Variables

- `BADGER_DB_PATH`: Sets the path to the Badger database directory.
  - **Default:** `./badger-data`
- `BADGER_LOG`: Enables Badger logging if set to `true`.
  - **Default:** `false`
- `PORT`: Sets the port for the web server.
  - **Default:** `8080`

---

## 🐳 Docker Deployment

Build and run with Docker:

```bash
make docker-build
make docker-run
```

Or use docker-compose:

```bash
docker-compose up -d
```

---

## 🛠️ Development

Install air for auto-reload (optional):

```bash
make tools
```

Run with auto-reload:

```bash
make dev
```

Build binary:

```bash
make build
```

---

## 📁 Project Structure

```
├── main.go              # Main application file
├── templates/
│   └── index.html       # HTML template with HTMX
├── go.mod               # Go module file
├── go.sum               # Go dependencies
├── Dockerfile           # Docker configuration
├── Makefile             # Build automation
└── README.md            # This file
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

---

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## 🙏 Acknowledgments

- [Badger](https://github.com/dgraph-io/badger) - Fast key-value DB in Go
- [HTMX](https://htmx.org/) - High power tools for HTML
- [GorillaMux](https://github.com/gorilla/mux) - HTTP router and URL matcher
