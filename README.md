# Go Account

To start your application in the dev profile, simply run:
```sh
> docker-compose up -d --build
```

```
/your-project
│
├── /cmd/                      # Command-line executables
│   └── /your-app/             # Main application entry point
│       └── main.go
│
├── /configs/                  # Configuration files (YAML, JSON, etc.)
│   └── config.yaml
│
├── /internal/                 # Private application and library code
│   ├── /server/               # Server initialization
│   ├── /versioning/           # API versioning setup
│   └── /utils/                # Utility functions
│
├── /pkg/                      # Public library code (importable by other projects)
│   └── /middlewares/          # Middleware implementations
│
├── /api/                      # API-specific code
│   ├── /v1/                   # Version 1 API
│   │   ├── /controllers/      # Controllers for handling requests
│   │   ├── /models/           # Database models
│   │   ├── /repositories/     # Data access layer
│   │   └── /services/         # Business logic
│   ├── /v2/                   # Version 2 API
│   │   ├── /controllers/      # Controllers for handling requests
│   │   ├── /models/           # Database models
│   │   ├── /repositories/     # Data access layer
│   │   └── /services/         # Business logic
│   └── router.go              # Central routing logic
│
├── /scripts/                  # Automation scripts (migrations, setup, etc.)
│   └── migrate.sh
│
├── /test/                     # Test files
│   └── /unit/                 # Unit tests
│   └── /integration/          # Integration tests
│
├── Dockerfile                 # Dockerfile for containerizing the app
├── go.mod                     # Go module file
└── go.sum                     # Go dependencies checksum file

```