# Service Catalog

## Overview
### Terminologies
- `Service` : A service can have a name, optional description and some versions associated with it.
- `Version` : A version is always associated with a service and can have a name and optional description.

### Requirements
An end-user should be able see an overview of services in their organization using the UI. For the same, the front-end would request the backend for the following:
- Fetching all services
- Fetching a specific service, along with all its versions
- Searching for a specific service based on filters (using name, description, etc)


## APIs
[OpenAPI Specification](./openapi_spec.yaml)

| Method | API                                           | Description                                                                                   |
|--------|-----------------------------------------------|-----------------------------------------------------------------------------------------------|
| GET    | /ping                                         | Healthcheck endpoint                                                                          |
| GET    | /metrics                                      | Shows the service's metrics - by default runs on different port                               |
| GET    | /v1/services                                  | Fetches all services, paginated and arranged in ascending order by name by default.           |
| GET    | /v1/service/:serviceName                      | Fetches the specific service information, along with all its versions - paginated by default. |
| POST   | /v1/service                                   | Creates a service using the information passed in request body.                               |
| POST   | /v1/service/version                           | Creates a service version using the information passed in request body.                       |
| PATCH  | /v1/service                                   | Updates a service's name, description using its id                                            |
| DELETE | /v1/service/:serviceName                      | Deletes a service, along with all its versions                                                |
| DELETE | /v1/service/:serviceName/version/:versionName | Deletes a service's version                                                                   |

### Future plans
Along with the above APIs, we can add Bulk APIs too for service and version creations or deletions. This API can take multiple inputs at once and process them asyncronously.
We can return a `202 Accepted` response to the user along with a requestId. They can later check the status of creation with the requestId.

Async processing here would help in keeping the UX clean, where the user doesn't have to wait for the updates and can come back later and check.

## How to run the project locally?
### Pre-requisties:
- Golang v1.22.2
- Docker

### Project Run
The project can be initiated and run with a single command:
- `make run`

This will initialise the project dependencies:
- A Postgres DB
- A Jaegar Instance

Both of these will be started as docker containers.
For the same the following ports are needed to be free: 5432, 5775, 6831, 6832, 5778, 9411, 16686, 14268

### How to access?
- Primary server would start on port 8080 by default. You can access the APIs via the url: http://localhost:8080/
- Additionally, a metrics server will begin on port 8081. Access it via http://localhost:8081/metrics

## Service Features

### Relational Database - Postgres
The service_catalog service uses a postgres database by default.
Initial schemas are linked [here](./migrations/0.sql).

DB configuration parameters can be passed using environment variables:
- DB_HOST
- DB_PORT
- DB_USER
- DB_PASSWORD
- DB_NAME

Default properties are added [here.](./config/db.go)

### Soft deletion
By default, DELETE APIs soft-delete the DB records. All GET requests ensure that soft-deleted records are not fetched.
Soft-deletion helps in recovering accidentally deleted services or versions.

A manual clean-up job can be set to run on a certain frequency - a week or a month. Post this, no recovery would be possible.

### Search Filters in APIs
The GET response of /services can be filtered via name or description. This can help in searching for a service.

### Paginated Response
The GET response of /services and /service/:serviceName is paginated by default.
Pagination helps in chunking a huge response into multiple smaller responses, thus saving network bandwidth as well as making the UX better.

### Sorted Response
The GET response of /services is sorted by name in ascending order by default. The user can choose to sort in descending order too using query parameters.

### Authentication
Except the /ping API, all service operation APIs have API key based authentication enabled.
API_AUTH_KEY can be passed as an environment variable for the setting the same.

### Rate-limiting
All service APIs are rate-limited, with the following configuration: 
- RPS            = 5
- BURST_REQUESTS = 10

This can be changed from [./internal/routes/middlewares.go](./internal/routes/middlewares.go)

### Observability
Observability is added in the service in the following ways:
- All request logs are added to `./log` folder. From here, we can send the logs to an external server periodically.
- All metrics are exposed on `http://localhost:8081/metrics` from where these can be scraped via Prometheus.
- Tracing is enabled and sent to a Jaegar instance. By default, we send it to a local docker instance on `localhost:6831` by UDP. This can be changed using environment variables.
- Uptime monitor can be set up on the /ping route.

### Testing
Unit tests are added for main controller functions. For the same, DB is mocked.

#### Future Test Plan

- Integration Tests can be added by mocking the web app and adding APIs

## Architecture Choices
Please check [Architecture.md](./Architecture.md)

## APIs functioning
Please check [Functioning.md](./Functioning.md) for working and screenshots.