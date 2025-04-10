package v1

import (
	"net/http"

	"github.com/wisdom-oss/common-go/v3/types"
)

// This file contains the commonly generated errors of the routes.

var ErrFileNotFound = types.ServiceError{
	Type:   "",
	Status: http.StatusNotFound,
	Title:  "File not found",
	Detail: "The file does not exist on the server. Please check your request",
}
