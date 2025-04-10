package v1

import (
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"

	objectStorage "microservice/internal/minio"
)

func Upload(c *gin.Context) {
	bucket := c.Param("bucket")
	filename := c.Param("filename")

	client := objectStorage.Client()
	bucketExists, err := client.BucketExists(c, bucket)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	if !bucketExists {
		err := client.MakeBucket(c, bucket, minio.MakeBucketOptions{})
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}
	}

	uploadFile, err := c.FormFile("file")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	f, err := uploadFile.Open()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	mime, err := mimetype.DetectReader(f)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	info, err := client.PutObject(c, bucket, filename, f, uploadFile.Size, minio.PutObjectOptions{
		ContentType: mime.String(),
	})
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.Header("ETag", info.ETag)
	c.Status(http.StatusCreated)
}
