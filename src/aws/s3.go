package aws

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func s3Session() *s3.S3 {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	log.Println(err)

	// Create s3 service client
	s3svc := s3.New(sess)

	return s3svc
}

// function for creating aws s3 bucket
func CreateBucket(bucket string) error {

	// starting s3 session
	s3svc := s3Session()

	// request for creating bucket
	_, err := s3svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	// checking error
	if err != nil {
		mapd := make(map[string]interface{})
		json.Unmarshal([]byte(err.Error()), &mapd)
		log.Println(mapd)
		log.Printf("Unable to create bucket %q, %v\n", bucket, err.Error())
		return err
	}

	// Wait until bucket is created before finishing
	log.Printf("Waiting for bucket %q to be created...\n", bucket)

	// waiting till bucket created
	err = s3svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Printf("Error occurred while waiting for bucket to be created, %v\n", bucket)
		return err
	}

	// when bucket created succesfully
	log.Printf("Bucket %q successfully created\n", bucket)
	return nil
}

//==================================================================================================//

func DeleteBucket(bucket string) error {
	// starting s3 session
	s3svc := s3Session()

	//list object in a bucket
	out, err := s3svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})

	// delete each object
	for i := 0; i < len(out.Contents); i++ {
		err = DeleteObject(bucket, aws.StringValue(out.Contents[i].Key))
		if err != nil {
			return err
		}
	}

	// request for creating bucket
	_, err = s3svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})

	// checking error
	if err != nil {
		log.Printf("Unable to delete bucket %q, %v\n", bucket, err)
		return err
	}

	// Wait until bucket is created before finishing
	log.Printf("Waiting for bucket %q to be deleted...\n", bucket)

	// waiting till bucket created
	err = s3svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Printf("Error occurred while waiting for bucket to be deleted, %v\n", bucket)
		return err
	}

	// when bucket deleted succesfully
	log.Printf("Bucket %q successfully deleted\n", bucket)
	return nil
}

//=================================================================================================

func DeleteObject(bucket, object string) error {
	// starting s3 session
	s3svc := s3Session()

	// request for creating bucket
	_, err := s3svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})

	// checking error
	if err != nil {
		log.Printf("Unable to delete object %q, %v\n", object, err)
		return err
	}

	// Wait until bucket is created before finishing
	log.Printf("Waiting for object %q to be deleted...\n", object)

	// waiting till bucket created
	err = s3svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})
	if err != nil {
		log.Printf("Error occurred while waiting for object to be deleted, %v\n", object)
		return err
	}

	// when bucket deleted succesfully
	log.Printf("Object %q successfully deleted\n", object)
	return nil
}
