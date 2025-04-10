package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"

	objectStorage "microservice/internal/minio"
)

func DeleteFile(c *gin.Context) {
	bucket := c.Param("bucket")
	fileName := c.Param("filename")

	_, err := objectStorage.Client().StatObject(c, bucket, fileName, defaultStatOpts)
	if err != nil {
		c.Abort()
		errCode := minio.ToErrorResponse(err).Code
		switch errCode {
		case ErrNoSuchKey:
			ErrFileNotFound.Emit(c)

		case ErrNoSuchBucket:
			ErrBucketNotFound.Emit(c)

		default:
			_ = c.Error(err)
		}

		return
	}

	err = objectStorage.Client().RemoveObject(c, bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
