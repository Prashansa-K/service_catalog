package controllers

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	constants "github.com/Prashansa-K/serviceCatalog/internal"
	api "github.com/Prashansa-K/serviceCatalog/internal/api/structs"
	"github.com/Prashansa-K/serviceCatalog/internal/models"
	"gorm.io/gorm"
)

func GetPaginatedServicesByFilters(db *gorm.DB, page int, sort, nameFilter, descriptionFilter string) (int64, []models.Service, error) {
	if nameFilter != "" {
		db = db.Where("LOWER(name) LIKE ?", "%"+nameFilter+"%")
	}

	if descriptionFilter != "" {
		db = db.Where("LOWER(description) LIKE ?", "%"+descriptionFilter+"%")
	}

	db = db.Model(&models.Service{})

	// Find the total count of all services with the above name and description
	var totalServices int64
	if err := db.Count(&totalServices).Error; err != nil {
		return -1, nil, err
	}

	totalPages := int(math.Ceil(float64(totalServices) / float64(constants.PAGE_SIZE)))
	if page > totalPages {
		return -1, nil, errors.New(constants.INVALID_PAGE_NUMBER)
	}

	offset := (page - 1) * constants.PAGE_SIZE
	sortByName := fmt.Sprintf("name %s", sort)

	var services []models.Service
	if err := db.Order(sortByName).Offset(offset).Limit(constants.PAGE_SIZE).Find(&services).Error; err != nil {
		return -1, nil, err
	}

	return totalServices, services, nil
}

func GetServiceByNameWithPaginatedVersions(db *gorm.DB, page int, serviceName string) (int64, *models.Service, error) {
	var totalVersions int64
	if err := db.Model(&models.Version{}).
		Joins("JOIN services ON versions.service_id = services.id").
		Where("services.name = ?", serviceName).
		Count(&totalVersions).Error; err != nil {
		return -1, nil, err
	}

	var service models.Service
	if err := db.Where("name = ?", serviceName).Preload("Versions", func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * constants.PAGE_SIZE).Limit(constants.PAGE_SIZE)
	}).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return -1, nil, errors.New(constants.SERVICE_RECORD_NOT_FOUND)
		}
		return -1, nil, err
	}

	return int64(totalVersions), &service, nil
}

func CreateService(db *gorm.DB, serviceRequest api.ServiceRequest) error {
	service := models.Service{
		Name:         serviceRequest.Name,
		Description:  serviceRequest.Description,
		VersionCount: 0, // Initial version count is 0
		CreatedAt:    time.Now(),
	}

	return db.Create(&service).Error
}

func CreateVersion(db *gorm.DB, versionRequest api.ServiceVersionRequest) error {
	var service models.Service
	if err := db.Where("name = ?", versionRequest.ServiceName).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(constants.SERVICE_RECORD_NOT_FOUND)
		}

		return err
	}

	version := models.Version{
		Name:        versionRequest.Name,
		ServiceID:   service.ID,
		Description: versionRequest.Description,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&version).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return errors.New(constants.DUPLICATE_VERSION_RECORD_ERROR)
		}

		return err
	}

	// Increment the version count for the service
	service.VersionCount++
	if err := db.Save(&service).Error; err != nil {
		return err
	}

	return nil
}

func DeleteService(db *gorm.DB, serviceName string) error {
	var service models.Service
	if err := db.Where("name = ?", serviceName).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(constants.SERVICE_RECORD_NOT_FOUND)
		}

		return err
	}

	// Soft deleting versions
	if err := db.Where("service_id = ?", service.ID).Delete(&models.Version{}).Error; err != nil {
		return err
	}

	// Soft delete the service
	if err := db.Delete(&service).Error; err != nil {
		return err
	}

	return nil
}

func DeleteVersion(db *gorm.DB, serviceName, versionName string) error {
	var service models.Service
	if err := db.Where("name = ?", serviceName).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(constants.SERVICE_RECORD_NOT_FOUND)
		}

		return errors.New(constants.ERROR_FETCHING_SERVICE)
	}

	var version models.Version
	if err := db.Where("service_id = ? AND name = ?", service.ID, versionName).First(&version).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(constants.VERSION_RECORD_NOT_FOUND)
		}

		return errors.New(constants.ERROR_FETCHING_SERVICE)
	}

	// Soft delete the version
	if err := db.Delete(&version).Error; err != nil {
		return err
	}

	// Decrement the version count for the service
	service.VersionCount--
	if err := db.Save(&service).Error; err != nil {
		return err
	}

	return nil
}

func UpdateService(db *gorm.DB, serviceRequest api.ServiceRequest) error {
	var service models.Service
	if err := db.Model(&models.Service{}).Where("id = ?", serviceRequest.ID).First(&service).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(constants.SERVICE_RECORD_NOT_FOUND)
		}

		return err
	}

	if serviceRequest.Name != "" {
		service.Name = serviceRequest.Name
	}

	if serviceRequest.Description != "" {
		service.Description = serviceRequest.Description
	}

	if err := db.Save(service).Error; err != nil {
		return err
	}

	return nil
}
