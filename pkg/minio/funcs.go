package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type Object struct {
	Name        string
	ContentType string
	Size        int64
	File        io.Reader
}

func (mc *Client) SetupBucket(ctx context.Context, bucketName string, ops minio.MakeBucketOptions) error {
	ok, err := mc.conn.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("mc.conn.BucketExists:%w", err)
	}

	if !ok {
		mc.conn.MakeBucket(ctx, bucketName, ops)
	}

	return nil
}

func (mc *Client) Upload(ctx context.Context, bucketName string, object Object, ops minio.PutObjectOptions) (minio.UploadInfo, error) {
	ops.ContentType = object.ContentType
	info, err := mc.conn.PutObject(ctx, bucketName, object.Name, object.File, object.Size, ops)
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("mc.conn.PutObject:%w", err)
	}

	return info, nil
}
