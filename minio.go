package object

import (
	"context"
	"errors"
	"io"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

type MinioStore struct {
	Client *minio.Client
}

func (s *MinioStore) ExistBucket(ctx context.Context, bucketName string) (bool, error) {
	return s.Client.BucketExists(ctx, bucketName)
}

func (s *MinioStore) CreateBucket(ctx context.Context, bucketName string) (bool, error) {
	makeBucketOption := minio.MakeBucketOptions{}

	if bucketName == "" {
		return false, errors.New(bucketName + " is not set")
	}

	if err := s.Client.MakeBucket(ctx, bucketName, makeBucketOption); err != nil {
		return false, err
	}

	return true, nil
}

func (s *MinioStore) ListObjects(ctx context.Context, object Object) ListObjectsResult {
	result := ListObjectsResult{BucketName: object.BucketName}

	listObjectsOption := minio.ListObjectsOptions{}
	for o := range s.Client.ListObjects(ctx, object.BucketName, listObjectsOption) {
		result.Append(
			Object{
				BucketName:  object.BucketName,
				Name:        o.Key,
				ObjectKey:   filepath.Join(object.Prefix, o.Key),
				Size:        o.Size,
				ContentType: o.ContentType,
				Error:       o.Err != nil,
			})
	}
	return result
}

func (s *MinioStore) StatObjectInfo(ctx context.Context, object Object) (*Object, error) {
	info, err := s.Client.StatObject(ctx, object.BucketName, object.Prefix, minio.StatObjectOptions{})

	if err != nil {
		return nil, err
	}

	result := &Object{
		BucketName:  object.BucketName,
		Name:        info.Key,
		Size:        info.Size,
		ContentType: info.ContentType,
		Error:       info.Err != nil,
		Prefix:      object.Prefix,
		ObjectKey:   filepath.Join(object.Prefix, info.Key),
	}

	return result, nil
}

func (s *MinioStore) DeleteObject(ctx context.Context, object Object) error {
	deleteObjectOption := minio.RemoveObjectOptions{}

	if err := s.Client.RemoveObject(ctx, object.BucketName, object.ObjectKey, deleteObjectOption); err != nil {
		return err
	}

	return nil
}

func (s *MinioStore) DeleteObjects(ctx context.Context, object Object) {
	deleted := make([]string, 0)

	objectCh := make(chan minio.ObjectInfo)
	removeCh := minio.RemoveObjectsOptions{}

	go func() {
		defer close(objectCh)
		listObjectsOption := minio.ListObjectsOptions{Prefix: object.ObjectKey, Recursive: true}

		for object := range s.Client.ListObjects(ctx, object.BucketName, listObjectsOption) {
			deleted = append(deleted, object.Key)
			objectCh <- object
		}
	}()

	s.Client.RemoveObjects(ctx, object.BucketName, objectCh, removeCh)
}

func (s *MinioStore) UploadObject(ctx context.Context, option UploadObjectOption) error {
	putObjectOption := minio.PutObjectOptions{ContentType: option.ContentType}
	_, err := s.Client.PutObject(ctx, option.BucketName, option.Prefix, option.File, option.Size, putObjectOption)
	return err
}

func (s *MinioStore) GetObject(ctx context.Context, object Object) (io.ReadCloser, error) {
	return s.Client.GetObject(ctx, object.BucketName, object.ObjectKey, minio.GetObjectOptions{})
}
