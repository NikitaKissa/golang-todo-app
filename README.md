# Golang Todo App

A learning-oriented Todo REST API built with Go and PostgreSQL.

The project was created as a practical exercise in building a backend application in Go using the standard `net/http` package, PostgreSQL, migrations, middleware, layered architecture, validation, structured logging, and Swagger documentation.

The application provides three main features:

* **Users** — create, retrieve, update, and delete users.
* **Tasks** — create, retrieve, update, and delete tasks associated with users.
* **Statistics** — get aggregated task statistics, including completion rate and average completion time.

> This is an educational project. Authentication, authorization, rate limiting, and other production-level concerns are intentionally out of scope.

---

## Features

### Users

The users feature provides a basic CRUD API:

* create a user;
* get a user by ID;
* get a paginated list of users;
* update a user;
* delete a user.

A user contains:

* `id`;
* `full_name`;
* `phone_number`;
* `version`.

The API validates user input, including the length of the name and phone number.

---

### Tasks

Tasks belong to an author user and support the following operations:

* create a task;
* get a task by ID;
* get a paginated list of tasks;
* filter tasks by user;
* update a task;
* delete a task;
* mark a task as completed or incomplete.

A task contains:

* `id`;
* `author_user_id`;
* `title`;
* `description`;
* `completed`;
* `created_at`;
* `completed_at`;
* `version`.

Task creation requires an `author_user_id` and a `title`. The description is optional and has a maximum length of 1000 characters.

Task listing supports pagination using `limit` and `offset`, as well as filtering by `user_id`.

---

### Statistics

The statistics endpoint provides aggregated information about tasks.

Currently available metrics include:

* total number of created tasks;
* number of completed tasks;
* completion rate;
* average task completion time.

Statistics can be filtered by:

* user;
* date range.

The endpoint accepts `user_id`, `from`, and `to` query parameters.

The response contains fields such as `tasks_created`, `tasks_completed`, `tasks_completed_rate`, and `tasks_average_completion_time`.

---

## Tech Stack

* **Go 1.26.6**
* **net/http** — HTTP server and routing
* **PostgreSQL** — primary database
* **pgx/v5** — PostgreSQL driver
* **golang-migrate** — database migrations
* **go-playground/validator** — request validation
* **google/uuid** — UUID generation where required
* **envconfig** — environment configuration
* **Uber Zap** — structured logging
* **Swaggo / HTTP Swagger** — API documentation
* **Docker / Docker Compose** — local development environment
* **Makefile** — project development commands

The project deliberately does not use an ORM or query builder. Database access is implemented using `pgx`, with an additional application-level abstraction built around interfaces.

---

## Architecture

The project follows a layered and feature-oriented architecture.

```text
.
├── cmd
│   └── todoapp
│       ├── Dockerfile
│       └── main.go
├── internal
│   ├── core
│   │   ├── config
│   │   ├── domain
│   │   ├── errors
│   │   ├── logger
│   │   ├── repository
│   │   └── transport
│   └── features
│       ├── statistics
│       ├── tasks
│       └── users
├── migrations
├── docs
├── public
├── docker-compose.yaml
├── Makefile
├── go.mod
└── README.md
```

### `cmd`

Contains the application entry point.

```text
cmd/todoapp/main.go
```

This is where the application is assembled and started.

---

### `internal/core`

Contains shared application infrastructure and abstractions.

```text
core/
├── config
├── domain
├── errors
├── logger
├── repository
└── transport
```

The `core` layer contains common concepts that are shared between different features, such as configuration, domain models, repository abstractions, errors, logging, and HTTP transport components.

---

### `internal/features`

Application functionality is split by business feature:

```text
features/
├── statistics
├── tasks
└── users
```

This keeps feature-specific code isolated and makes it easier to reason about individual parts of the application.

Each feature is responsible for its own business logic and transport/repository integration.

---

### Repository abstraction

The project uses `pgx` for communication with PostgreSQL, but database access is not coupled directly to the business logic.

Repository interfaces provide an additional abstraction layer between the application and the concrete PostgreSQL implementation.

This approach makes the boundaries between business logic and infrastructure explicit and provides a cleaner architecture without introducing an ORM.

---

## Middleware

The HTTP server uses middleware for common cross-cutting concerns.

The project includes middleware for:

* request logging;
* panic recovery;
* CORS;
* request ID / request tracing where applicable.

Rate limiting and content-type middleware are intentionally not implemented.

Authentication and authorization are also not part of the project.

---

## API

The API is versioned under:

```text
/api/v1
```

The default local server address is:

```text
127.0.0.1:8080
```

The API is documented using Swagger. The repository contains generated Swagger/OpenAPI files under:

```text
docs/
├── docs.go
├── swagger.json
└── swagger.yaml
```

The Swagger specification describes the available endpoints, request parameters, request bodies, responses, and validation rules.

---

## Endpoints

### Users

| Method   | Endpoint                  | Description               |
| -------- | ------------------------- | ------------------------- |
| `GET`    | `/api/v1/users`           | Get users with pagination |
| `POST`   | `/api/v1/users`           | Create a user             |
| `GET`    | `/api/v1/users/{user_id}` | Get a user by ID          |
| `PATCH`  | `/api/v1/users/{user_id}` | Update a user             |
| `DELETE` | `/api/v1/users/{user_id}` | Delete a user             |

The users collection supports pagination using:

```text
?limit=<limit>&offset=<offset>
```

User creation requires `full_name` and optionally accepts `phone_number`.

The update endpoint uses partial-update semantics. For nullable fields such as `phone_number`, the API distinguishes between an omitted field, a provided value, and `null`. `full_name` cannot be set to `null`.

---

### Tasks

| Method   | Endpoint                  | Description               |
| -------- | ------------------------- | ------------------------- |
| `GET`    | `/api/v1/tasks`           | Get tasks with pagination |
| `POST`   | `/api/v1/tasks`           | Create a task             |
| `GET`    | `/api/v1/tasks/{task_id}` | Get a task by ID          |
| `PATCH`  | `/api/v1/task/{task_id}`  | Update a task             |
| `DELETE` | `/api/v1/tasks/{task_id}` | Delete a task             |

> Note: the update endpoint is currently documented as `/api/v1/task/{task_id}`, while the other task endpoints use `/api/v1/tasks/...`. The Swagger specification is the source of truth for the current API contract.

Task listing supports:

```text
?user_id=<user_id>&limit=<limit>&offset=<offset>
```

A task can be updated using partial-update semantics. The API supports changing the title, description, and completion status.

---

### Statistics

| Method | Endpoint             | Description         |
| ------ | -------------------- | ------------------- |
| `GET`  | `/api/v1/statistics` | Get task statistics |

Optional query parameters:

```text
user_id
from
to
```

Example:

```text
GET /api/v1/statistics?user_id=1&from=2026-08-01&to=2026-08-19
```

The date parameters use the `YYYY-MM-DD` format.

---

## Example Requests

### Create a user

```http
POST /api/v1/users
Content-Type: application/json

{
  "full_name": "Jack Daniels",
  "phone_number": "+48123123123"
}
```

---

### Create a task

```http
POST /api/v1/tasks
Content-Type: application/json

{
  "author_user_id": 1,
  "title": "Make homework",
  "description": "Exercise 15 p. 12"
}
```

The API requires `author_user_id` and `title` when creating a task.

---

### Update a task

```http
PATCH /api/v1/task/1
Content-Type: application/json

{
  "completed": true
}
```

---

### Get tasks for a specific user

```http
GET /api/v1/tasks?user_id=1&limit=20&offset=0
```

---

### Get statistics

```http
GET /api/v1/statistics?user_id=1
```

Example response:

```json
{
  "tasks_created": 34,
  "tasks_completed": 29,
  "tasks_completed_rate": 85.2941,
  "tasks_average_completion_time": "15m34s"
}
```

The response format corresponds to the generated API schema.

---

## Database & Migrations

PostgreSQL is used as the primary persistent storage.

Database schema changes are managed using `golang-migrate`.

Migrations are located in:

```text
migrations/
├── 000001_init.up.sql
└── 000001_init.down.sql
```

The project keeps both `up` and `down` migrations, allowing the schema to be applied and rolled back in a controlled way.

### Running migrations

The project exposes migration-related commands through the `Makefile`.

Replace the placeholders below with the corresponding targets from your `Makefile`:

```bash
# Apply all pending migrations
make migrate-up

# Roll back the latest migration
make migrate-down
```

---

## Configuration

Application configuration is loaded from environment variables.

Create a local environment file according to the configuration expected by the application.

Example:

```env
LOGGER_LEVEL=DEBUG
TIME_ZONE=UTC

POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_TIMEOUT=

HTTP_ADDRESS=:8080
HTTP_TIMEOUT=30s
```

> Replace the example variables above with the exact variables used by the application.

For local development, PostgreSQL can be started using Docker Compose.

---

## Running the Project

### Prerequisites

Make sure the following are installed:

* Go 1.26+
* Docker
* Docker Compose
* Make

### 1. Clone the repository

```bash
git clone <https://github.com/NikitaKissa/golang-todo-app>
cd golang-todo-app
```

### 2. Start PostgreSQL

The project includes a `docker-compose.yaml` file for the local development environment.

```bash
docker compose up -d
```

### 3. Configure environment variables

Set the environment variables required by the application.

```bash
cp .env.example .env
```

> If `.env.example` is not included in the repository, create the required environment variables manually according to the application's configuration package.

### 4. Run migrations

```bash
make migrate-up
```

### 5. Start the application

The application is normally started through the project's `Makefile`:

```bash
make todoapp-run
```

The server is configured to run on:

```text
127.0.0.1:8080
```

---

## Makefile

The project uses a `Makefile` as the main entry point for common development operations.

Add the actual targets from your `Makefile` here.

For example:

```bash
make todoapp-deploy
make todoapp-run
make migrate-up
make migrate-down
make swagger-gen
```

This keeps the development workflow simple and avoids having to remember long commands for common operations.

---

## Docker

The project includes Docker configuration for local development and application containerization.

Relevant files:

```text
docker-compose.yaml
cmd/todoapp/Dockerfile
```

Docker Compose is used to run the PostgreSQL dependency locally, while the application has its own Dockerfile under `cmd/todoapp`.

---

## Swagger Documentation

The project uses Swaggo to generate API documentation.

Generated documentation is stored in:

```text
docs/
├── docs.go
├── swagger.json
└── swagger.yaml
```

The Swagger specification contains the complete API contract, including:

* endpoints;
* query parameters;
* path parameters;
* request bodies;
* response schemas;
* HTTP status codes;
* validation constraints.

For example, the generated schema documents both successful responses and common errors such as `400`, `404`, `409`, and `500`.

If Swagger generation is exposed through the `Makefile`, document the corresponding command here:

```bash
make swagger-gen
```

---

## Error Handling

The API uses a common error response structure:

```json
{
  "error": "system error message: error type",
  "message": "human readable message"
}
```

The `error` field is intended for the technical error information, while `message` contains a human-readable description.

Typical HTTP responses include:

* `200 OK` — successful read/update;
* `201 Created` — successful creation;
* `204 No Content` — successful deletion;
* `400 Bad Request` — invalid input;
* `404 Not Found` — requested entity does not exist;
* `409 Conflict` — conflicting operation;
* `500 Internal Server Error` — unexpected server-side error.

---

## Concurrency & Versioning

Users and tasks contain a `version` field.

The version is returned as part of the API representation and can be used as a basis for optimistic concurrency control.

For example:

```json
{
  "id": 1,
  "title": "Make homework",
  "completed": false,
  "version": 1
}
```

This is one of the project concepts that goes beyond a minimal CRUD implementation and demonstrates how resource versioning can be incorporated into a backend API.

---

## Validation

Request validation is implemented using `go-playground/validator`.

Examples of validation constraints documented by the API include:

### User

```text
full_name
- required
- minimum length: 3
- maximum length: 100

phone_number
- minimum length: 10
- maximum length: 15
```

### Task

```text
title
- required
- minimum length: 1
- maximum length: 100

description
- optional
- minimum length: 1
- maximum length: 1000
```

## These constraints are part of the generated API contract.

## Logging

The application uses Uber Zap for structured logging.

Logs are written to the project's output directory during local development:

```text
out/
└── logs/
```

The repository currently contains generated log files from development runs. In a production repository, generated runtime logs should generally not be committed to version control.

---

## Testing

Automated tests have not been implemented yet.

This is an intentional limitation of the current learning project and a natural next step for further development.

Potential areas for future tests include:

* HTTP handlers;
* services/business logic;
* repository implementations;
* validation;
* statistics calculations;
* middleware;
* integration tests with PostgreSQL.

---

## Project Goals

The main goal of this project is not to build a production-ready Todo application.

Instead, it serves as a practical Go backend project for learning and practicing:

* Go project structure;
* `net/http`;
* REST API design;
* PostgreSQL integration;
* `pgx`;
* repository abstractions;
* interfaces;
* separation of concerns;
* feature-oriented architecture;
* middleware;
* request validation;
* structured logging;
* database migrations;
* Docker;
* environment-based configuration;
* Swagger/OpenAPI;
* pagination;
* filtering;
* partial updates;
* error handling;
* resource versioning;
* basic statistics and aggregation.

The project intentionally avoids adding unnecessary complexity such as authentication, authorization, rate limiting, background workers, message brokers, or microservices.

---

## Project Structure

A simplified overview:

```text
golang-todo-app/
│
├── cmd/
│   └── todoapp/
│       ├── Dockerfile
│       └── main.go
│
├── internal/
│   ├── core/
│   │   ├── config/
│   │   ├── domain/
│   │   ├── errors/
│   │   ├── logger/
│   │   ├── repository/
│   │   └── transport/
│   │
│   └── features/
│       ├── statistics/
│       ├── tasks/
│       └── users/
│
├── migrations/
│   ├── 000001_init.down.sql
│   └── 000001_init.up.sql
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── public/
├── docker-compose.yaml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## What Is Not Included

This project intentionally does **not** implement:

* authentication;
* authorization;
* user sessions;
* JWT;
* OAuth;
* password management;
* rate limiting;
* message brokers;
* background jobs;
* distributed architecture;
* automated tests.

These features can be added later as separate learning exercises.

---

## Possible Next Steps

Some natural extensions for the project are:

1. Add unit tests for services and handlers.
2. Add integration tests using a real PostgreSQL instance.
3. Add authentication and authorization.
4. Add task filtering by completion status.
5. Add sorting for task lists.
6. Improve pagination metadata.
7. Add more detailed statistics.
8. Add CI with GitHub Actions.
9. Add linting and static analysis.
10. Improve Docker setup for production-like environments.
11. Add API contract tests.
12. Add graceful shutdown and more advanced HTTP server configuration.

---

## Learning Takeaways

After working through this project, the main concepts practiced are:

```text
HTTP
  ↓
Middleware
  ↓
Transport / Handlers
  ↓
Feature Logic
  ↓
Repository Interfaces
  ↓
PostgreSQL / pgx
```

The important part of the project is not the Todo domain itself, but the engineering patterns around it: how HTTP requests are validated and handled, how business logic is separated from infrastructure, how database access is abstracted, how migrations are managed, and how a small Go application can be organized without introducing unnecessary frameworks.

---

## License

No license has been specified for this project yet.
