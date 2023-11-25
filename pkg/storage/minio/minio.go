package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github/yogabagas/join-app/config"
	"io"
	"log"
	"net/url"
	"os"
	"time"
)

type Minio interface {
	SetBucket(bucketName string) *MinioImpl
	SetLocation(location string) *MinioImpl
	WithContext(ctx context.Context) *MinioImpl
	MakeBucket() *MinioImpl
	UploadObject(objectName string, filePath string, contentType string) (bool, error)
	BucketExist() bool
	ListBuckets() []minio.BucketInfo
	RemoveBucket() bool
	ListObjects() (response map[string]interface{})
	ListIncompleteUploads(objectPrefix string, isRecursive bool) (response map[string]interface{})
	SetBucketPolicy(policy string) bool
	FGetObject(objectName string, filePath string) bool
	GetObject(objectName string, objectOptions minio.GetObjectOptions, localPath string) bool
	PutObject(objectName string, filePath string) bool
	StatObject(objectName string) (objectInfo minio.ObjectInfo)
	RemoveObject(objectName string) bool
	RemoveIncompleteUpload(objectName string) bool
	SignUrl(objectName string) *url.URL
}

type MinioImpl struct {
	Client   *minio.Client
	Bucket   string
	Location string
	ctx      context.Context
}

var MinioClient *MinioImpl

func NewMinio() *MinioImpl {
	minioClient, err := minio.New(config.GlobalCfg.Storage.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.GlobalCfg.Storage.Minio.AccessKey, config.GlobalCfg.Storage.Minio.SecretKey, config.GlobalCfg.Storage.Minio.Token),
		Secure: config.GlobalCfg.Storage.Minio.Secure,
	})

	if err != nil {
		log.Fatalln(err)
	}

	MinioClient = &MinioImpl{Client: minioClient}
	return &MinioImpl{Client: minioClient}
}

func (m *MinioImpl) SetBucket(bucketName string) *MinioImpl {
	m.Bucket = bucketName
	return m
}

func (m *MinioImpl) SetLocation(location string) *MinioImpl {
	m.Location = location
	return m
}

func (m *MinioImpl) WithContext(ctx context.Context) *MinioImpl {
	m.ctx = ctx
	return m
}

func (m *MinioImpl) MakeBucket() *MinioImpl {
	err := m.Client.MakeBucket(m.ctx, m.Bucket, minio.MakeBucketOptions{Region: m.Location})
	if err != nil {
		fmt.Println(err)
		return m
	}
	log.Printf("Successfully created %s\n", m.Bucket)
	return m
}

// UploadObject API Reference : File Object Operations
// FPutObject uploads objects that are less than 128MiB in a single PUT operation.
// For objects that are greater than the 128MiB in size, FPutObject seamlessly
// uploads the object in chunks of 128MiB or more depending on the actual file size.
// The max upload size for an object is 5TB.
func (m *MinioImpl) UploadObject(objectName string, filePath string, contentType string) (bool, error) {
	info, err := m.Client.FPutObject(m.ctx, m.Bucket, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return false, err
	}
	fmt.Println("Successfully uploaded / of size", objectName, info)
	return true, err
}

func (m *MinioImpl) BucketExist() bool {
	exist, err := m.Client.BucketExists(m.ctx, m.Bucket)
	if err != nil && !exist {
		return false
	}
	return true
}

func (m *MinioImpl) ListBuckets() []minio.BucketInfo {
	buckets, err := m.Client.ListBuckets(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	return buckets
}

func (m *MinioImpl) RemoveBucket() bool {
	err := m.Client.RemoveBucket(m.ctx, m.Bucket)
	if err != nil {
		// fmt.Println(err)
		return false
	}
	return true
}

func (m *MinioImpl) ListObjects() (response map[string]interface{}) {
	contextCancel, cancel := context.WithCancel(m.ctx)

	defer cancel()

	objectSource := m.Client.ListObjects(contextCancel, m.Bucket, minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    "",
	})

	for object := range objectSource {
		if object.Err != nil {
			// fmt.Println(object.Err)
			return response
		}
		response["listObjects"] = object
	}
	return response
}

func (m *MinioImpl) ListIncompleteUploads(objectPrefix string, isRecursive bool) (response map[string]interface{}) {
	multiPartObjectUpload := m.Client.ListIncompleteUploads(context.Background(), m.Bucket, objectPrefix, isRecursive)
	for multiPartObject := range multiPartObjectUpload {
		if multiPartObject.Err != nil {
			// fmt.Println(multiPartObject.Err)
			return response
		}
		response["ListInComplete"] = multiPartObject
	}
	return response
}

// API Reference : Bucket policy Operations
// policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::my-bucketname/*"],"Sid": ""}]}`
func (m *MinioImpl) SetBucketPolicy(policy string) bool {
	err := m.Client.SetBucketPolicy(context.Background(), m.Bucket, policy)
	if err != nil {
		// fmt.Println(err, bucketName)
		return false
	}
	return true
}

func (m *MinioImpl) GetBucketPolicy() string {
	policy, err := m.Client.GetBucketPolicy(m.ctx, m.Bucket)
	if err != nil {
		return ""
	}
	return policy
}

// Get Object file to local system
func (m *MinioImpl) FGetObject(objectName string, filePath string) bool {
	err := m.Client.FGetObject(m.ctx, m.Bucket, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		// fmt.Println(err)
		return false
	}
	return true
}

// Returns a stream of the object data. Most of the common errors occur when reading the stream.
// localPath := "/tmp/local-file.jpg"
func (m *MinioImpl) GetObject(objectName string, objectOptions minio.GetObjectOptions, localPath string) bool {
	object, err := m.Client.GetObject(m.ctx, m.Bucket, objectName, objectOptions)
	if err != nil {
		fmt.Println(err)
		return false
	}

	localFile, err := os.Create(localPath)
	if err != nil {
		// fmt.Println(err)
		return false
	}

	if _, err = io.Copy(localFile, object); err != nil {
		// fmt.Println(err)
		return false
	}
	return true
}

// Uploads objects that are less than 128MiB in a single PUT operation. For objects that are greater than 128MiB
// in size, PutObject seamlessly uploads the object as parts of 128MiB or more depending on the actual file size.
// The max upload size for an object is 5TB.
func (m *MinioImpl) PutObject(objectName string, filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return false
	}

	_, err = m.Client.PutObject(m.ctx, m.Bucket, objectName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return false
	}
	return true
}

func (m *MinioImpl) StatObject(objectName string) (objectInfo minio.ObjectInfo) {
	objInfo, err := m.Client.StatObject(m.ctx, m.Bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return objInfo
	}
	return objInfo
}

func (m *MinioImpl) CopyObject(bucketNameSource string, objectNameSource string, bucketNameDestination string, objectNameDestination string) bool {
	// Use-case 1: Simple copy object with no conditions.
	// Source object
	srcOpts := minio.CopySrcOptions{
		Bucket: bucketNameSource,
		Object: objectNameSource,
	}

	// Destination object
	dstOpts := minio.CopyDestOptions{
		Bucket: bucketNameDestination,
		Object: objectNameDestination,
	}

	// Copy object call
	_, err := m.Client.CopyObject(m.ctx, dstOpts, srcOpts)
	if err != nil {
		return false
	}
	return true
}

func (m *MinioImpl) RemoveObject(objectName string) bool {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
		// VersionID:        versionId, //Version ID of the object to delete
	}
	err := m.Client.RemoveObject(m.ctx, m.Bucket, objectName, opts)
	if err != nil {
		return false
	}
	return true
}

func (m *MinioImpl) RemoveIncompleteUpload(objectName string) bool {
	err := m.Client.RemoveIncompleteUpload(m.ctx, m.Bucket, objectName)
	if err != nil {
		return false
	}
	return true
}

//	Presigned URL operations
//
// Set request parameters for content-disposition.
func (m *MinioImpl) SignUrl(objectName string) *url.URL {
	reqParams := make(url.Values)
	// Generates a presigned url which expires in a day.
	presignedURL, err := m.Client.PresignedGetObject(m.ctx, m.Bucket, objectName, time.Duration(1000)*time.Second, reqParams)
	if err != nil {
		// fmt.Println(err)
		return presignedURL
	}
	return presignedURL
}
