# Go Account

To start your application in the dev profile, simply run:

```sh
docker-compose up -d --build
```

### Generate private key

```sh
openssl genpkey -algorithm RSA -out ./config/rsa2048_private.key -pkeyopt rsa_keygen_bits:2048
```

### Extract the public key from the private key

```sh
openssl rsa -pubout -in ./config/rsa2048_private.key -out ./config/rsa2048_public.key
```

```
/your-project
├── Dockerfile
├── README.md
├── cmd
│   └── server
│       └── server.go                       # Server initialization
├── config
│   └── config.go                           # Configuration files (YAML, JSON, etc.)
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal                                # Private application and library code
│   ├── handler                             # Http handling requests
│   │   ├── oauth_handler.go
│   │   └── user_handler.go
│   ├── domain                               # Database domains
│   │   ├── client.go
│   │   ├── oauth_access_token.go
│   │   ├── oauth_refresh_token.go
│   │   └── user.go
│   ├── repository                          # Data access layer
│   │   ├── access_token_repository.go
│   │   ├── client_repository.go
│   │   ├── refresh_token_repository.go
│   │   └── user_repository.go
│   └── usecase                             # Business logic
│       ├── oauth_usecase.go
│       └── user_usecase.go
├── main.go                                 # Application entry point
├── pkg                                     # Public library code (importable by other projects)
│   ├── database
│   │   ├── database.go
│   │   └── mongodb.go
│   ├── middleware                          # Middleware implementations
│   │   ├── authenticate.go
│   │   ├── authorize.go
│   │   ├── error_handler.go
│   │   └── token.go
│   └── utils                               # Utility functions
│       ├── custom_errors.go
│       ├── custom_validator.go             # Custom validators
│       ├── jwt_auth.go
│       └── utility_funcs.go
├── scripts                                 # Automation scripts (migrations, setup, etc.)
└── test
    ├── integration
    └── unit
```
