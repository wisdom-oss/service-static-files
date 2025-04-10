package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"

	"microservice/internal"
	objectStorage "microservice/internal/minio"
)

// GetFile retrieves the provided file from the default bucket ("main directory")
// and returns it to the user.
func GetFile(c *gin.Context) {
	bucket := c.Param("bucket")
	fileName := c.Param("filename")

	if strings.TrimSpace(bucket) == "" {
		bucket = internal.DefaultMinioPublicBucketName
	}

	info, err := objectStorage.Client().StatObject(c, bucket, fileName, defaultStatOpts)
	if err != nil {
		c.Abort()
		errCode := minio.ToErrorResponse(err).Code
		if errCode == ErrNoSuchKey {
			ErrFileNotFound.Emit(c)
			return
		}
		_ = c.Error(err)
		return
	}

	obj, err := objectStorage.Client().GetObject(c, bucket, fileName, defaultGetOpts)
	if err != nil {
		c.Abort()
		resp := minio.ToErrorResponse(err)
		if resp.Code == ErrNoSuchKey {
			ErrFileNotFound.Emit(c)
			return
		}
		_ = c.Error(err)
		return
	}

	c.Header("ETag", info.ETag)
	c.Header("Last-Modified", info.LastModified.Format(time.RFC1123))

	c.DataFromReader(http.StatusOK, -1, info.ContentType, obj, nil)

}
