package internal

const (
	// Service related constants
	PAGE_SIZE = 2
	ASC       = "ASC"
	DESC      = "DESC"

	// 200
	SUCCESS                 = "Success"
	SERVICE_DELETED         = "Service Deleted Successfully"
	SERVICE_VERSION_DELETED = "Service Version Deleted Successfully"

	// 201
	SERVICE_CREATED         = "Service Created Successfully"
	SERVICE_VERSION_CREATED = "Version Created Successfully"

	// 202
	ACCEPTED = "Accepted"

	// 4xx
	INVALID_REQUEST_BODY           = "invalid request body"
	INVALID_PAGE_NUMBER            = "invalid page number"
	SERVICE_RECORD_NOT_FOUND       = "service not found"
	VERSION_RECORD_NOT_FOUND       = "version not found"
	DUPLICATE_VERSION_RECORD_ERROR = "version with the same name already exists for this service"

	//5xx
	INTERNAL_SERVER_ERROR  = "internal server error"
	ERROR_FETCHING_SERVICE = "error fetching service"
)
