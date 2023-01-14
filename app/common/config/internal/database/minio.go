package database

import (
	"douyin/app/common/config/internal/consts"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (g *Group) NewMinioClient() (*minio.Client, error) {
	if g.agollo == nil {
		return nil, consts.ErrEmptyConfigClient
	}

	v, err := g.agollo.GetViper(consts.MainNamespace)
	if err != nil {
		return nil, err
	}

	endpoint := v.GetString("Minio.Endpoint")
	username := v.GetString("Minio.Username")
	password := v.GetString("Minio.Password")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:           credentials.NewStaticV4(username, password, ""),
		Secure:          true,
		Transport:       nil,
		Region:          "",
		BucketLookup:    0,
		TrailingHeaders: false,
		CustomMD5:       nil,
		CustomSHA256:    nil,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
