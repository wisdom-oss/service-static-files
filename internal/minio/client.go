package objectStorage

import (
	"errors"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"microservice/internal"
)

const (
	KeyHost         = internal.ConfigKey_Minio_Host
	KeyClientID     = internal.ConfigKey_Minio_ClientID
	KeyClientSecret = internal.ConfigKey_Minio_ClientSecret
	KeyUseSSL       = internal.ConfigKey_Minio_UseSSL
)

const healthcheckDuration = 15 * time.Second

var c *minio.Client

var (
	errNoMinioEndpoint     = errors.New("no minio endpoint set")
	errNoMinioClientID     = errors.New("no minio client id set")
	errNoMinioClientSecret = errors.New("no minio client secret set")
	errMinioOffline        = errors.New("unable to connect to minio")
)

func Connect() (err error) {
	slog.Info("initializing object storage connection")

	config := internal.Configuration()

	if !config.IsSet(KeyHost) {
		return errNoMinioEndpoint
	}
	if !config.IsSet(KeyClientID) {
		return errNoMinioClientID
	}
	if !config.IsSet(KeyClientSecret) {
		return errNoMinioClientSecret
	}

	endpoint := config.GetString(KeyHost)
	clientID := config.GetString(KeyClientID)
	clientSecret := config.GetString(KeyClientSecret)
	useSSL := config.GetBool(KeyUseSSL)

	opts := minio.Options{
		Creds:  credentials.NewStaticV4(clientID, clientSecret, ""),
		Secure: useSSL,
	}

	c, err = minio.New(endpoint, &opts)
	if err != nil {
		return err
	}

	_, _ = c.HealthCheck(healthcheckDuration)
	if c.IsOffline() {
		return errMinioOffline
	}

	return nil
}

func Client() *minio.Client {
	return c
}
