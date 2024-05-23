package v1

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	constants "github.com/Prashansa-K/serviceCatalog/internal"
	api "github.com/Prashansa-K/serviceCatalog/internal/api/structs"
	"github.com/Prashansa-K/serviceCatalog/internal/controllers"
	"github.com/Prashansa-K/serviceCatalog/internal/db"

	"github.com/labstack/echo/v4"
)

func GetServices(ctx echo.Context) error {
	db, err := db.GetDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	// get paging information
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
		err = nil
	}

	// get sorting information
	sort := strings.ToUpper(ctx.QueryParam("sort"))
	if sort != constants.ASC && sort != constants.DESC {
		sort = constants.ASC
	}

	// Per function tracing
	// var services []models.Service
	// var ok bool
	// rawData := tracer.TraceServiceFunc(ctx, controllers.GetPaginatedServicesByFilters, db, page, sort, ctx.QueryParam("name"), ctx.QueryParam("description"))
	// totalServices := rawData.([]reflect.Value)[0].Int()
	// if services, ok = rawData.([]reflect.Value)[1].Interface().([]models.Service); !ok {
	// 	services = nil
	// }
	// if rawData.([]reflect.Value)[2].Interface() != nil {
	// 	err = rawData.([]reflect.Value)[2].Interface().(error)
	// }

	totalServices, services, err := controllers.GetPaginatedServicesByFilters(db, page, sort, ctx.QueryParam("name"), ctx.QueryParam("description"))

	if err != nil {
		if err.Error() == constants.INVALID_PAGE_NUMBER {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var response []api.ServiceResponse
	for _, service := range services {
		response = append(response, api.ServiceResponse{
			ID:           service.ID,
			Name:         service.Name,
			Description:  service.Description,
			VersionCount: service.VersionCount,
		})
	}

	totalPages := int(math.Ceil(float64(totalServices) / float64(constants.PAGE_SIZE)))

	return ctx.JSON(http.StatusOK, api.ServicePaginationResponse{
		Services:     response,
		TotalPages:   totalPages,
		CurrentPage:  page,
		TotalRecords: totalServices,
	})
}

func GetService(ctx echo.Context) error {
	db, err := db.GetDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	// get paging information
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	totalVersions, service, err := controllers.GetServiceByNameWithPaginatedVersions(db, page, ctx.Param("serviceName"))
	if err != nil {
		if err.Error() == constants.SERVICE_RECORD_NOT_FOUND {
			return ctx.JSON(http.StatusNotFound, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var versions []api.ServiceVersion
	for _, version := range service.Versions {
		versions = append(versions, api.ServiceVersion{
			Name:        version.Name,
			Description: version.Description,
			CreatedAt:   version.CreatedAt,
		})
	}

	totalPages := int((totalVersions + constants.PAGE_SIZE - 1) / constants.PAGE_SIZE)

	return ctx.JSON(http.StatusOK, api.ServiceResponseWithVersionPagination{
		ID:                  service.ID,
		Name:                service.Name,
		Description:         service.Description,
		CreatedAt:           service.CreatedAt,
		VersionCount:        service.VersionCount,
		Versions:            versions,
		TotalPages:          totalPages,
		CurrentPage:         page,
		TotalVersionRecords: totalVersions,
	})
}

func CreateService(ctx echo.Context) error {
	db, err := db.GetDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var serviceRequest api.ServiceRequest
	if err := ctx.Bind(&serviceRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": constants.INVALID_REQUEST_BODY,
		})
	}

	if err := controllers.CreateService(db, serviceRequest); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, echo.Map{
		"message": constants.SERVICE_CREATED,
	})
}

func CreateVersion(ctx echo.Context) error {
	db, err := db.GetDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var versionRequest api.ServiceVersionRequest
	if err := ctx.Bind(&versionRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": constants.INVALID_REQUEST_BODY,
		})
	}

	if err := controllers.CreateVersion(db, versionRequest); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, echo.Map{
		"message": constants.SERVICE_VERSION_CREATED,
	})
}

func DeleteService(ctx echo.Context) error {
	db, err := db.GetDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := controllers.DeleteService(db, ctx.Param("serviceName")); err != nil {
		if err.Error() == constants.SERVICE_RECORD_NOT_FOUND {
			return ctx.JSON(http.StatusNotFound, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": constants.SERVICE_DELETED,
	})
}

func DeleteVersion(ctx echo.Context) error {
	db, err := db.GetDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := controllers.DeleteVersion(db, ctx.Param("serviceName"), ctx.Param("versionName")); err != nil {
		if err.Error() == constants.SERVICE_RECORD_NOT_FOUND {
			return ctx.JSON(http.StatusNotFound, err.Error())
		}

		if err.Error() == constants.VERSION_RECORD_NOT_FOUND {
			return ctx.JSON(http.StatusNotFound, err.Error())
		}

		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": constants.SERVICE_VERSION_DELETED,
	})
}

func UpdateService(ctx echo.Context) error {
	db, err := db.GetDB()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var serviceRequest api.ServiceRequest
	if err := ctx.Bind(&serviceRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": constants.INVALID_REQUEST_BODY,
		})
	}

	if err := controllers.UpdateService(db, serviceRequest); err != nil {
		if err.Error() == constants.SERVICE_RECORD_NOT_FOUND {
			return ctx.JSON(http.StatusNotFound, err.Error())
		}

		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": http.StatusAccepted,
	})
}
