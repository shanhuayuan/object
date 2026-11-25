package object

import (
	"context"
	"io"
	"mime/multipart"
)

type StoreConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Token     string
	UseSsl    bool
}

type Object struct {
	BucketName  string `json:"bucketName"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	Error       bool   `json:"error"`
	Prefix      string `json:"prefix"`
	ObjectKey   string `json:"objectKey"`
}

type Bucket struct {
	Name string `json:"name"`
}

type ListObjectsResult struct {
	BucketName string   `json:"bucketName"`
	Objects    []Object `json:"objects"`
	Total      int      `json:"total"`
}

type DeleteObjectsResult struct {
	Objects []Object
	Total   int
}

type UploadObjectOption struct {
	BucketName  string `form:"bucketName" binding:"required"`
	Prefix      string `form:"prefix" binding:"required"`
	File        multipart.File
	Size        int64
	ContentType string
}

func (r *ListObjectsResult) Append(object Object) {
	r.Objects = append(r.Objects, object)
	r.Total += 1
}

type Store interface {
	ExistBucket(ctx context.Context, bucketName string) (bool, error)
	CreateBucket(ctx context.Context, bucketName string) (bool, error)
	ListObjects(ctx context.Context, object Object) ListObjectsResult
	UploadObject(ctx context.Context, option UploadObjectOption) error
	StatObjectInfo(ctx context.Context, object Object) (*Object, error)
	GetObject(ctx context.Context, object Object) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, object Object) error
	DeleteObjects(ctx context.Context, object Object)
}
