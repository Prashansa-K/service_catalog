## Echo framework
Choice of a web framework was done keeping the following in mind:
- Ease of use: To get onboarded quickly
- High performance
- Rich collection of middlewares
- HTTP/2 support: for future-proofing. Even if the current service version doesn't have the ability for the same, we can do a protocol upgrade in future without major changes
- Easy data binding: Request payloads can be easily binded for processing and support all common formats - json, xml or form-data.

## Relational Database
We had the following requirements:
- Adding, removing, fetching a service or version.
- Fetching all services.
- Searching services with filters (name, description) or directly via name.
- Returning paginated response.
- Sorting the response.

Along with that, following would be necessary too:
- A version should always be associated with a service.
- A service should not have >1 versions with the same name.

A relational database seemed like a straightforward choice here. Along with CRUD operations, filter, order by and uniqueness support, a relational database also promises low latency in these operations.

#### Future plans
- Indexing on the database can help us reduce the look-up times more.
- Caching get responses with appropriate TTLs can be used for faster look-ups.
- Connection pooling can be used to reduce DB connection setup time. Requests can reuse the free connections. This would enhance the performance further.

## ORM
gorm is used as the ORM. 
- An ORM can help to keep track of any changes made in DB models. If we use native queries, a manual search and update would be required. 
- ORM also helps to abstract the SQL queries, thus reducing any query inaccuracies.
- gorm ORM helps in direct soft-deletion too, without any additional checks.

## Middlewares
Except for the healthcheck route, all routes have the following middlewares:
- Rate Limiter: to protect against DDos attacks
- Authentication: to ensure that only authenticated users can access the APIs
- Removing Trailing Slash: to ensure that requests fail if an extra slash is added in the path
- Logger
- Tracer:
    Trace Function is provided at [./internal/api/tracer.go](./internal/api/tracer.go) which can be called to create spans for internal functions.
    [./internal/api/v1/services.go, line 36](./internal/api/v1/services.go) - commented out part shows how to use tracing for all functions.
- Metrics

## Security Features - Future plans
- Creating TLS certificates using Let's Encrypt and serving the APIs using HTTPS.
- Using a TLS connection with the database.

