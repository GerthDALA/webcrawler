# Web Crawler

A production-ready web crawler built with Go, following Domain-Driven Design (DDD) architecture principles.

## Features

- Concurrent webpage crawling with configurable concurrency
- Domain-Driven Design architecture
- Result-type error handling (similar to Rust's Result)
- Machine learning for content analysis:
  - Text vectorization and nearest neighbor search
  - Named Entity Recognition (NER)
  - Topic modeling
  - Content classification
- PostgreSQL database with pgvector for vector embeddings
- REST API for controlling the crawler and accessing content
- Command-line interface
- Docker and Docker Compose for easy deployment

## Architecture

The project follows a Domain-Driven Design (DDD) architecture with the following layers:

- **Domain Layer**: Contains the core domain models, entities, and business logic
- **Application Layer**: Orchestrates the domain layer and provides services to the interfaces
- **Infrastructure Layer**: Implements the interfaces defined in the domain layer
- **Interface Layer**: Exposes the application via REST API and CLI

## Directory Structure

```
webcrawler/
├── cmd/
│   └── crawler/         # Application entry point
├── internal/
│   ├── domain/          # Domain models and interfaces
│   ├── application/     # Application services
│   ├── infrastructure/  # Implementation of interfaces
│   └── interfaces/      # User interfaces (REST API, CLI)
├── pkg/                 # Public packages
├── deployments/         # Deployment configurations
├── test/                # Test helpers and fixtures
├── Makefile             # Build automation
└── README.md            # Project documentation
```

## Prerequisites

- Go 1.20+
- PostgreSQL 14+
- Docker and Docker Compose (for containerized deployment)
- pgvector extension for PostgreSQL

## Installation

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/webcrawler.git
   cd webcrawler
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Build the application:
   ```bash
   make build
   ```

4. Set up the database:
   ```bash
   make migrate-up
   ```

5. Run the application:
   ```bash
   make run
   ```

### Docker Deployment

1. Build and run with Docker Compose:
   ```bash
   make docker-run
   ```

## Usage

### Command Line Interface

The crawler provides a command-line interface for common operations:

```bash
# Start the crawler with a seed URL
./webcrawler crawl --seed https://example.com --concurrency 10 --depth 3

# Run the API server
./webcrawler server --host localhost --port 8080

# Show crawler statistics
./webcrawler stats

# Search for content
./webcrawler search --query "keyword" --limit 10

# Analyze content
./webcrawler analyze content --id <content-id>
./webcrawler analyze text --text "Text to analyze"
```

### REST API

The crawler exposes a REST API for controlling the crawler and accessing content:

```bash
# Add a seed URL
curl -X POST http://localhost:8080/api/crawler/seed -H "Content-Type: application/json" -d '{"url": "https://example.com"}'

# Start the crawler
curl -X POST http://localhost:8080/api/crawler/start

# Get crawler statistics
curl http://localhost:8080/api/crawler/stats

# Search for content
curl http://localhost:8080/api/content/search?q=keyword&limit=10

# Get content by ID
curl http://localhost:8080/api/content/{id}

# Analyze text
curl -X POST http://localhost:8080/api/analysis/text -H "Content-Type: application/json" -d '{"text": "Text to analyze"}'
```

## Configuration

The application is configured using a YAML configuration file located at `config.yaml`. You can also override configuration settings using environment variables.

Example configuration:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: webcrawler
  sslmode: disable

crawler:
  user_agent: "WebCrawler/1.0 (+https://example.com/bot)"
  max_depth: 3
  concurrency: 10
  politeness_delay: 1000  # milliseconds
  timeout: 30             # seconds
  max_redirects: 5
  follow_redirects: true
  allowed_domains: []     # empty means all domains
  allowed_extensions: ["html", "htm", "php", "asp", "aspx", "jsp"]
  disallowed_paths: ["/admin", "/login", "/logout", "/register", "/cart", "/checkout"]
  allowed_content_types: ["text/html", "application/xhtml+xml"]
  max_url_length: 2048
  retry_count: 3
  retry_delay: 5000       # milliseconds

ml:
  vector_dimensions: 384
  min_term_frequency: 2
  max_features: 20000
  num_topics: 10

api:
  host: localhost
  port: 8080
```

## Development

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Code Formatting

```bash
make fmt
```

### Generating Documentation

```bash
make doc
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.