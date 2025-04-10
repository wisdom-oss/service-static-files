package v1

import (
	"io"
	"log/slog"
	"net/http"
	"path"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/errgroup"

	objectStorage "microservice/internal/minio"
)

func Upload(c *gin.Context) {
	bucket := c.Param("bucket")
	basepath := c.Param("basepath")

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

	form, err := c.MultipartForm()
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	files, ok := form.File["file"]
	if !ok {
		// todo: send error
		return
	}

	var uploadGroup errgroup.Group
	var uploadResults []minio.UploadInfo
	var mutex sync.Mutex

	for _, file := range files {

		uploadGroup.Go(func() error {
			slog.Warn("uploading file", "filename", file.Filename)
			f, err := file.Open()
			if err != nil {
				return err
			}

			mime, err := mimetype.DetectReader(f)
			if err != nil {
				return err
			}

			_, err = f.Seek(0, io.SeekStart)
			if err != nil {
				return err
			}

			filename := path.Join(basepath, file.Filename)

			info, err := client.PutObject(c, bucket, filename, f, file.Size, minio.PutObjectOptions{
				ContentType: mime.String(),
			})

			if err != nil {
				return err
			}

			mutex.Lock()
			uploadResults = append(uploadResults, info)
			mutex.Unlock()
			return nil

		})
	}

	if err := uploadGroup.Wait(); err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusCreated)

}
