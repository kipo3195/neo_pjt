package repository

import (
	"context"
	"file/internal/consts"
	"file/internal/domain/fileUrl/entity"
	"file/internal/domain/fileUrl/repository"
	"log"
	"time"

	// OCI SDK 필수 패키지
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

type oraclefileUrlStorageRepositoryImpl struct {
	client     objectstorage.ObjectStorageClient
	namespace  string
	bucketName string
}

func NewFileUrlStorageRepository(client objectstorage.ObjectStorageClient, namespace string, bucketName string) repository.FileUrlStorageRepository {
	return &oraclefileUrlStorageRepositoryImpl{
		client:     client,
		namespace:  namespace,
		bucketName: bucketName,
	}
}

func (r *oraclefileUrlStorageRepositoryImpl) CreateFileUrl(ctx context.Context, en entity.CreateFileUrlEntity) ([]entity.CreateFileUrlResultEntity, error) {

	result := make([]entity.CreateFileUrlResultEntity, 0)

	for key := range en.FileInfoMap {

		fileInfoEntity := en.FileInfoMap[key]

		createdUrl, err := r.getUploadUrl(ctx, r.bucketName, fileInfoEntity.FileId, time.Minute*60)

		if err != nil {
			return nil, consts.ErrFileUrlCreateError
		}
		temp := entity.CreateFileUrlResultEntity{
			FileId:     key,
			FileName:   fileInfoEntity.FileName,
			CreatedUrl: createdUrl,
		}
		result = append(result, temp)
	}

	return result, nil
}

func (o *oraclefileUrlStorageRepositoryImpl) getUploadUrl(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {

	req := objectstorage.CreatePreauthenticatedRequestRequest{
		NamespaceName: &o.namespace,
		BucketName:    &bucketName,
		CreatePreauthenticatedRequestDetails: objectstorage.CreatePreauthenticatedRequestDetails{
			Name:        common.String("upload_" + objectName),
			AccessType:  objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectreadwrite,
			ObjectName:  &objectName,
			TimeExpires: &common.SDKTime{Time: time.Now().Add(expiry)},
		},
	}

	resp, err := o.client.CreatePreauthenticatedRequest(ctx, req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return "https://objectstorage.ap-chuncheon-1.oraclecloud.com" + *resp.AccessUri, nil
}
