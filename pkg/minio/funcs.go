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

func (c *Client) SetupBucket(ctx context.Context, bucketName string, ops minio.MakeBucketOptions) error {
	ok, err := c.conn.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("mc.conn.BucketExists:%w", err)
	}

	if !ok {
		c.conn.MakeBucket(ctx, bucketName, ops)
	}

	return nil
}

func (c *Client) Upload(ctx context.Context, bucketName string, object Object, ops minio.PutObjectOptions) (minio.UploadInfo, error) {
	ops.ContentType = object.ContentType
	info, err := c.conn.PutObject(ctx, bucketName, object.Name, object.File, object.Size, ops)
	if err != nil {
		return minio.UploadInfo{}, fmt.Errorf("mc.conn.PutObject:%w", err)
	}

	return info, nil
}

func (c *Client) Download(ctx context.Context, bucketName string, objectName string, ops minio.GetObjectOptions) (*Object, error) {
	obj, err := c.conn.GetObject(ctx, bucketName, objectName, ops)
	if err != nil {
		return nil, fmt.Errorf("c.conn.GetObject: %w", err)
	}

	return &Object{Name: objectName, File: obj}, nil
}
