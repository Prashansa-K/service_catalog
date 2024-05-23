package controllers

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	constants "github.com/Prashansa-K/serviceCatalog/internal"
	api "github.com/Prashansa-K/serviceCatalog/internal/api/structs"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormMockDB *gorm.DB
var mock sqlmock.Sqlmock

func initMockDB() error {
	db, sqlmocker, err := sqlmock.New()
	if err != nil {
		return err
	}

	mock = sqlmocker

	gormMockDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	return nil
}

func TestGetServiceByNameWithPaginatedVersions_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "versions" JOIN services ON versions.service_id = services.id WHERE services.name = $1 AND "versions"."deleted_at" IS NULL`)).
		WithArgs("test-service").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE name = $1 AND "services"."deleted_at" IS NULL ORDER BY "services"."id" LIMIT $2`)).
		WithArgs("test-service", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow("123", "test-service", "Test service", constants.PAGE_SIZE))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "versions" WHERE "versions"."service_id" = $1 AND "versions"."deleted_at" IS NULL LIMIT $2`)).
		WithArgs(int64(123), constants.PAGE_SIZE).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "service_id", "description"}).
			AddRow("1", "v1", "123", "Version 1").
			AddRow("2", "v2", "123", "Version 2"))

	// Create the controller and call the method
	totalVersions, service, err := GetServiceByNameWithPaginatedVersions(gormMockDB, 1, "test-service")

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalVersions)
	assert.Equal(t, "test-service", service.Name)
	assert.Equal(t, "Test service", service.Description)
	assert.Len(t, service.Versions, 2)
	assert.Equal(t, "v1", service.Versions[0].Name)
	assert.Equal(t, "v2", service.Versions[1].Name)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPaginatedServicesByFilters_NoFilters_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	// Expect the query to be executed
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "services" WHERE "services"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE "services"."deleted_at" IS NULL ORDER BY name ASC LIMIT $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow(123, "test-service-1", "Test check 1", 2).
			AddRow(456, "test-service-2", "Test check 2", 1))

	// Create the controller and call the method
	totalServices, services, err := GetPaginatedServicesByFilters(gormMockDB, 1, constants.ASC, "", "")

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalServices)
	assert.Len(t, services, 2)
	assert.Equal(t, uint(123), services[0].ID)
	assert.Equal(t, "test-service-1", services[0].Name)
	assert.Equal(t, "Test check 1", services[0].Description)
	assert.Equal(t, uint(456), services[1].ID)
	assert.Equal(t, "test-service-2", services[1].Name)
	assert.Equal(t, "Test check 2", services[1].Description)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPaginatedServicesByFilters_NameFilters_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	// Expect the query to be executed
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "services" WHERE LOWER(name) LIKE $1 AND "services"."deleted_at" IS NULL`)).
		WithArgs(`%test%`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE LOWER(name) LIKE $1 AND "services"."deleted_at" IS NULL ORDER BY name ASC LIMIT $2`)).
		WithArgs(`%test%`, constants.PAGE_SIZE).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow(123, "test-service-1", "Test check 1", 2).
			AddRow(456, "test-service-2", "Test check 2", 1))

	// Create the controller and call the method
	totalServices, services, err := GetPaginatedServicesByFilters(gormMockDB, 1, constants.ASC, "test", "")

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalServices)
	assert.Len(t, services, 2)
	assert.Equal(t, uint(123), services[0].ID)
	assert.Equal(t, "test-service-1", services[0].Name)
	assert.Equal(t, "Test check 1", services[0].Description)
	assert.Equal(t, uint(456), services[1].ID)
	assert.Equal(t, "test-service-2", services[1].Name)
	assert.Equal(t, "Test check 2", services[1].Description)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPaginatedServicesByFilters_DescriptionFilters_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	// Expect the query to be executed
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "services" WHERE LOWER(description) LIKE $1 AND "services"."deleted_at" IS NULL`)).
		WithArgs(`%check%`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE LOWER(description) LIKE $1 AND "services"."deleted_at" IS NULL ORDER BY name ASC LIMIT $2`)).
		WithArgs(`%check%`, constants.PAGE_SIZE).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow(123, "test-service-1", "Test check 1", 2).
			AddRow(456, "test-service-2", "Test check 2", 1))

	// Create the controller and call the method
	totalServices, services, err := GetPaginatedServicesByFilters(gormMockDB, 1, constants.ASC, "", "check")

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalServices)
	assert.Len(t, services, 2)
	assert.Equal(t, uint(123), services[0].ID)
	assert.Equal(t, "test-service-1", services[0].Name)
	assert.Equal(t, "Test check 1", services[0].Description)
	assert.Equal(t, uint(456), services[1].ID)
	assert.Equal(t, "test-service-2", services[1].Name)
	assert.Equal(t, "Test check 2", services[1].Description)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPaginatedServicesByFilters_AllFilters_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	// Expect the query to be executed
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "services" WHERE LOWER(name) LIKE $1 AND LOWER(description) LIKE $2 AND "services"."deleted_at" IS NULL`)).
		WithArgs(`%test%`, `%check%`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(2))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE LOWER(name) LIKE $1 AND LOWER(description) LIKE $2 AND "services"."deleted_at" IS NULL ORDER BY name ASC LIMIT $3`)).
		WithArgs(`%test%`, `%check%`, constants.PAGE_SIZE).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow(123, "test-service-1", "Test check 1", 2).
			AddRow(456, "test-service-2", "Test check 2", 1))

	// Create the controller and call the method
	totalServices, services, err := GetPaginatedServicesByFilters(gormMockDB, 1, constants.ASC, "test", "check")

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalServices)
	assert.Len(t, services, 2)
	assert.Equal(t, uint(123), services[0].ID)
	assert.Equal(t, "test-service-1", services[0].Name)
	assert.Equal(t, "Test check 1", services[0].Description)
	assert.Equal(t, uint(456), services[1].ID)
	assert.Equal(t, "test-service-2", services[1].Name)
	assert.Equal(t, "Test check 2", services[1].Description)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateService_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	// Expect the query to be executed
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "services" ("name","description","created_at","version_count") VALUES ($1,$2,$3,$4) RETURNING "deleted_at","id"`)).
		WithArgs("test-service", "Test service", sqlmock.AnyArg(), 0).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow("123"))
	mock.ExpectCommit()

	serviceRequest := api.ServiceRequest{
		Name:        "test-service",
		Description: "Test service",
	}
	err := CreateService(gormMockDB, serviceRequest)

	// Assert the results
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateVersion_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE name = $1 AND "services"."deleted_at" IS NULL ORDER BY "services"."id" LIMIT $2`)).
		WithArgs("test-service", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow("123", "test-service", "Test service", 1))

	// Expect the query to be executed
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "versions" ("service_id","name","created_at","description") VALUES ($1,$2,$3,$4) RETURNING "deleted_at","description","id"`)).
		WithArgs(123, "v1", sqlmock.AnyArg(), "Version 1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "services" SET "name"=$1,"description"=$2,"created_at"=$3,"deleted_at"=$4,"version_count"=$5 WHERE "services"."deleted_at" IS NULL AND "id" = $6`)).
		WithArgs("test-service", "Test service", sqlmock.AnyArg(), nil, 2, 123).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	versionRequest := api.ServiceVersionRequest{
		Name:        "v1",
		ServiceName: "test-service",
		Description: "Version 1",
	}
	err := CreateVersion(gormMockDB, versionRequest)

	// Assert the results
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateVersion_ServiceNotFound(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE name = $1 AND "services"."deleted_at" IS NULL ORDER BY "services"."id" LIMIT $2`)).
		WithArgs("test-non-existing-service", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	versionRequest := api.ServiceVersionRequest{
		Name:        "v1",
		ServiceName: "test-non-existing-service",
		Description: "Version 1",
	}

	err := CreateVersion(gormMockDB, versionRequest)

	// Assert the results
	assert.Error(t, err)
	assert.Equal(t, "service not found", err.Error())

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateVersion_DuplicateVersions(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE name = $1 AND "services"."deleted_at" IS NULL ORDER BY "services"."id" LIMIT $2`)).
		WithArgs("test-service", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow("123", "test-service", "Test service", 1))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "versions" ("service_id","name","created_at","description") VALUES ($1,$2,$3,$4) RETURNING "deleted_at","description","id"`)).
		WithArgs(123, "v1", sqlmock.AnyArg(), "Version 1").
		WillReturnError(errors.New("duplicate key value violates unique constraint \"versions_service_id_name_key\""))

	versionRequest := api.ServiceVersionRequest{
		Name:        "v1",
		ServiceName: "test-service",
		Description: "Version 1",
	}

	err := CreateVersion(gormMockDB, versionRequest)

	// Assert the results
	assert.Error(t, err)
	assert.Equal(t, "version with the same name already exists for this service", err.Error())

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateVersion_ErrorInVersionCreation(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE name = $1 AND "services"."deleted_at" IS NULL ORDER BY "services"."id" LIMIT $2`)).
		WithArgs("test-service", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow("123", "test-service", "Test service", 1))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "versions" ("service_id","name","created_at","description") VALUES ($1,$2,$3,$4) RETURNING "deleted_at","description","id"`)).
		WithArgs(123, "v1", sqlmock.AnyArg(), "Version 1").
		WillReturnError(errors.New("some other error"))

	versionRequest := api.ServiceVersionRequest{
		Name:        "v1",
		ServiceName: "test-service",
		Description: "Version 1",
	}

	err := CreateVersion(gormMockDB, versionRequest)

	// Assert the results
	assert.Error(t, err)
	assert.Equal(t, "error creating version", err.Error())

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteService_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE name = $1 AND "services"."deleted_at" IS NULL ORDER BY "services"."id" LIMIT $2`)).
		WithArgs("test-service", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow("123", "test-service", "Test service", 1))

	// Expect the query to be executed
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "versions" SET "deleted_at"=$1 WHERE service_id = $2 AND "versions"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 123).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "services" SET "deleted_at"=$1 WHERE "services"."id" = $2 AND "services"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 123).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Create the controller and call the method

	err := DeleteService(gormMockDB, "test-service")

	// Assert the results
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteVersionInDB_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE name = $1 AND "services"."deleted_at" IS NULL ORDER BY "services"."id" LIMIT $2`)).
		WithArgs("test-service", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow("123", "test-service", "Test service", 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "versions" WHERE (service_id = $1 AND name = $2) AND "versions"."deleted_at" IS NULL ORDER BY "versions"."id" LIMIT $3`)).
		WithArgs(123, "v1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow("123", "test-service", "Test service", 1))

	// Expect the query to be executed
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "versions" SET "deleted_at"=$1 WHERE "versions"."id" = $2 AND "versions"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), 123).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "services" SET "name"=$1,"description"=$2,"created_at"=$3,"deleted_at"=$4,"version_count"=$5 WHERE "services"."deleted_at" IS NULL AND "id" = $6`)).
		WithArgs("test-service", "Test service", sqlmock.AnyArg(), nil, 0, 123).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Create the controller and call the method

	err := DeleteVersion(gormMockDB, "test-service", "v1")

	// Assert the results
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateService_Success(t *testing.T) {
	if (gormMockDB == nil) || (mock == nil) {
		// Setup the mock DB
		err := initMockDB()
		assert.NoError(t, err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "services" WHERE id = $1 AND "services"."deleted_at" IS NULL ORDER BY "services"."id" LIMIT $2`)).
		WithArgs(123, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "version_count"}).
			AddRow("123", "test-service", "Test service", 1))

	// Expect the query to be executed
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "services" SET "id"=$1,"name"=$2,"description"=$3,"created_at"=$4,"deleted_at"=$5,"version_count"=$6 WHERE "services"."deleted_at" IS NULL AND "id" = $7`)).
		WithArgs(123, "test-service-2", "Test service 2", sqlmock.AnyArg(), nil, 1, 123).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	serviceRequest := api.ServiceRequest{
		ID:          123,
		Name:        "test-service-2",
		Description: "Test service 2",
	}
	err := UpdateService(gormMockDB, serviceRequest)

	// Assert the results
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
