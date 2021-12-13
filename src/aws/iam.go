package aws

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	authCore "git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
)

// types for creating policy
type StatementEntry struct {
	Effect   string
	Action   []string
	Resource interface{}
}

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

func iamSession() *iam.IAM {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	log.Println(err)

	// Create IAM service client
	iamsvc := iam.New(sess)

	return iamsvc
}

// CreateIamUser function for creating iam user
// @params userid -> id of user
func CreateIamUser(userid string) (string, error) {
	// starting iam service session
	iamsvc := iamSession()

	// creating iam user by passing params
	create, err := iamsvc.CreateUser(&iam.CreateUserInput{
		UserName: aws.String(strings.ToLower("sl" + userid)),
	})

	if err != nil {
		log.Println("CreateUser Error", err)
		return "", err
	} else {
		return aws.StringValue(create.User.UserName), nil
	}
}

//=========================================================================================================//

func DeleteIamUser(userid string) error {
	err := DetachPolicies("sl" + userid)
	if err != nil {
		return err
	}

	// starting iam service session
	iamsvc := iamSession()

	// deleting iam user by passing params
	_, err = iamsvc.DeleteUser(&iam.DeleteUserInput{
		UserName: aws.String(strings.ToLower("sl" + userid)),
	})
	if err != nil {
		log.Println("DeleteUser Error", err)
		return err
	} else {
		return nil
	}
}

//=========================================================================================================//

// AttachPolicy function for attaching policies to iam user
// @params username -> iam username of user
// @params policyType -> type of policy Individual, Enterprise  or Admin
// @params enterprise -> company name when type is Enterprise or Admin
func AttachPolicy(username, policyType, enterprise string) error {
	// checking username is blank
	if username == "" {
		log.Println("username is empty")
		return errors.New("username is empty")
	}

	// connecting to db
	db := config.DB

	// fetch policies from db on baisi of policy type
	var policy []database.Policy
	db.Where("p_type=? OR (p_type=? AND company=?)", "all", policyType, enterprise).Find(&policy)

	// creating aws iam session
	iamsvc := iamSession()

	// attaching policies to users
	for i := 0; i < len(policy); i++ {
		_, err := iamsvc.AttachUserPolicy(&iam.AttachUserPolicyInput{
			PolicyArn: aws.String(policy[i].Arn),
			UserName:  aws.String(username),
		})

		if err != nil {
			log.Println("GroupUser Error", err)
			return err
		}
	}
	return nil
}

//=========================================================================================================//

// list attach policies
func ListAttachPolicies(username string) ([]string, error) {
	// checking username is empty
	if username == "" {
		log.Println("username is empty")
		return nil, errors.New("username is empty")
	}

	// initialize slice of strings to save arns of policies
	arns := make([]string, 0)

	// creating aws iam session
	iamsvc := iamSession()
	// calling aws iam service function to list attached user policies by passing appropriate parameters
	err := iamsvc.ListAttachedUserPoliciesPages(&iam.ListAttachedUserPoliciesInput{
		// passing iam username of user
		UserName: aws.String(username),
	}, func(out *iam.ListAttachedUserPoliciesOutput, is bool) bool {
		// handle function for saving policy arn in above defined array
		for i := 0; i < len(out.AttachedPolicies); i++ {
			arns = append(arns, aws.StringValue(out.AttachedPolicies[i].PolicyArn))
		}
		return is
	})

	if err != nil {
		return nil, err
	}
	return arns, nil
}

//=========================================================================================================//

// DetachPolicies is a method to de-attach all policies corresponding to user
func DetachPolicies(username string) error {
	// checking username is empty
	if username == "" {
		log.Println("username is empty")
		return errors.New("username is empty")
	}

	// fetch attached policies to user
	arns, err := ListAttachPolicies(username)
	if err != nil {
		return err
	}

	// creating aws iam session
	iamsvc := iamSession()
	// detach policies one be one
	for i := 0; i < len(arns); i++ {
		_, err = iamsvc.DetachUserPolicy(&iam.DetachUserPolicyInput{
			UserName:  aws.String(username),
			PolicyArn: aws.String(arns[i]),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// DetachPolicy is a method to de-attach policies
func DetachPolicy(username, arn string) error {
	// checking username is empty
	if username == "" {
		log.Println("username is empty")
		return errors.New("username is empty")
	}

	// creating aws iam session
	iamsvc := iamSession()
	// detach policies one be one
	_, err := iamsvc.DetachUserPolicy(&iam.DetachUserPolicyInput{
		UserName:  aws.String(username),
		PolicyArn: aws.String(arn),
	})
	if err != nil {
		return err
	}

	return nil
}

//=========================================================================================================//

// CreatePolicy function for creating policy
// @params policyType -> type of policy Enterprise  or Admin
// @params enterprise -> company name
func CreatePolicy(enterprise, policyType string) error {

	// checking enterprise name is blank
	if enterprise == "" {
		log.Println("company name is empty")
		return errors.New("company name is empty")
	}

	// slugify company name
	company := enterprise

	// declaring bucket resource
	bucket := "arn:aws:s3:::" + company + "-sl"

	// calculating object resource
	var object string

	switch policyType {
	case "Enterprise":
		object = "arn:aws:s3:::" + company + "-sl/private-dataset-${aws:username}/*"
	case "Admin":
		object = "arn:aws:s3:::" + company + "-sl/*"
	}

	// calculating policy document to be created
	doc := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"s3:PutAnalyticsConfiguration",
					"s3:GetObjectVersionTagging",
					"s3:CreateBucket",
					"s3:ReplicateObject",
					"s3:GetObjectAcl",
					"s3:GetBucketObjectLockConfiguration",
					"s3:DeleteBucketWebsite",
					"s3:PutLifecycleConfiguration",
					"s3:GetObjectVersionAcl",
					"s3:DeleteObject",
					"s3:GetBucketPolicyStatus",
					"s3:GetObjectRetention",
					"s3:GetBucketWebsite",
					"s3:PutReplicationConfiguration",
					"s3:PutObjectLegalHold",
					"s3:GetObjectLegalHold",
					"s3:GetBucketNotification",
					"s3:PutBucketCORS",
					"s3:GetReplicationConfiguration",
					"s3:ListMultipartUploadParts",
					"s3:PutObject",
					"s3:GetObject",
					"s3:PutBucketNotification",
					"s3:PutBucketLogging",
					"s3:GetAnalyticsConfiguration",
					"s3:PutBucketObjectLockConfiguration",
					"s3:GetObjectVersionForReplication",
					"s3:GetLifecycleConfiguration",
					"s3:ListBucketByTags",
					"s3:GetInventoryConfiguration",
					"s3:GetBucketTagging",
					"s3:PutAccelerateConfiguration",
					"s3:DeleteObjectVersion",
					"s3:GetBucketLogging",
					"s3:ListBucketVersions",
					"s3:RestoreObject",
					"s3:ListBucket",
					"s3:GetAccelerateConfiguration",
					"s3:GetBucketPolicy",
					"s3:PutEncryptionConfiguration",
					"s3:GetEncryptionConfiguration",
					"s3:GetObjectVersionTorrent",
					"s3:AbortMultipartUpload",
					"s3:GetBucketRequestPayment",
					"s3:GetObjectTagging",
					"s3:GetMetricsConfiguration",
					"s3:DeleteBucket",
					"s3:PutBucketVersioning",
					"s3:GetBucketPublicAccessBlock",
					"s3:ListBucketMultipartUploads",
					"s3:PutMetricsConfiguration",
					"s3:GetBucketVersioning",
					"s3:GetBucketAcl",
					"s3:PutInventoryConfiguration",
					"s3:GetObjectTorrent",
					"s3:PutBucketWebsite",
					"s3:PutBucketRequestPayment",
					"s3:PutObjectRetention",
					"s3:GetBucketCORS",
					"s3:GetBucketLocation",
					"s3:ReplicateDelete",
					"s3:GetObjectVersion", // Allow for creating log groups
				},
				Resource: []string{
					bucket,
					object,
				},
			},
			{
				Effect: "Allow",
				// Allows for DeleteItem, GetItem, PutItem, Scan, and UpdateItem
				Action: []string{
					"s3:GetAccountPublicAccessBlock",
					"s3:ListAllMyBuckets",
					"s3:HeadBucket",
				},
				Resource: "*",
			},
		},
	}

	// marshaling above doc in bytes
	b, err := json.Marshal(&doc)
	if err != nil {
		log.Println("Error marshaling policy", err)
		return err
	}

	// creating aws iam session
	iamsvc := iamSession()

	// call create policy method of aws iam service
	result, err := iamsvc.CreatePolicy(&iam.CreatePolicyInput{
		PolicyName:     aws.String("stacklabs_policy_for_" + policyType + "_" + company),
		PolicyDocument: aws.String(string(b)),
	})

	// any error in creating policy
	if err != nil {
		log.Println(err)
		return err
	}

	// connecting to db
	db := config.DB

	// save policy in db
	db.Create(&database.Policy{
		Arn:     aws.StringValue(result.Policy.Arn),
		PType:   policyType,
		Company: enterprise,
	})

	return nil
}

//=========================================================================================================//

// DeletePolicies function for creating policy
// @params policyType -> type of policy Enterprise  or Admin
// @params enterprise -> company name
func DeletePolicies(enterprise string) error {
	// checking enterprise name is blank
	if enterprise == "" {
		log.Println("company name is empty")
		return errors.New("company name is empty")
	}

	// connecting to db
	db := config.DB

	var policy []database.Policy
	db.Where("company=?", enterprise).Find(&policy)

	// creating aws iam session
	iamsvc := iamSession()

	for i := 0; i < len(policy); i++ {
		// call delete policy method of aws iam service
		_, err := iamsvc.DeletePolicy(&iam.DeletePolicyInput{
			PolicyArn: aws.String(policy[i].Arn),
		})

		// any error in deleting policy
		if err != nil && !(strings.Contains(err.Error(), "status code: 404") || strings.Contains(err.Error(), "status code: 409")) {
			log.Println(err)
			return err
		}
	}
	return nil
}

//=========================================================================================================//

// MailAccessKeys function for genrating access keys of a user and send them on mail to user
func MailAccessKeys(user authCore.Accounts, reqtype string) (int, string) {

	// checking user has already have an access key generated in db
	// connecting to db
	db := config.DB

	// intialize variable with type access keys
	var userkeys []database.AccessKeys
	// fetching data on basis of userid
	db.Where("userid= ?", user.Userid).Find(&userkeys)
	if len(userkeys) != 0 {
		if reqtype != "fetch" {
			//send keys in mail
			go sendMail(user, userkeys[0].Key, userkeys[0].Secret)
			return 200, "We have sent credentials to your email, please check your email."
		}
		return 200, "Access keys generated"
	}

	// starting iam service session
	iamsvc := iamSession()

	//check is there any keys attached to that user
	res, err := iamsvc.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String("sl" + user.Userid),
	})
	if err != nil {
		log.Println(err)
		return 500, "Unable to generate credentials"
	}
	list := res.AccessKeyMetadata
	//delete previous keys if any
	for i := 0; i < len(list); i++ {
		_, err = iamsvc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: list[i].AccessKeyId,
			UserName:    list[i].UserName,
		})
		if err != nil {
			log.Println(err)
			return 500, "Unable to generate credentials"
		}
	}
	//generate new keys
	keysData, err := iamsvc.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String("sl" + user.Userid),
	})
	if err != nil {
		log.Println(err)
		return 500, "Unable to generate credentials."
	}
	keys := keysData.AccessKey

	db.Create(&database.AccessKeys{
		Userid: user.Userid,
		Key:    aws.StringValue(keys.AccessKeyId),
		Secret: aws.StringValue(keys.SecretAccessKey),
	})

	if reqtype != "fetch" {
		//send keys in mail
		go sendMail(user, aws.StringValue(keys.AccessKeyId), aws.StringValue(keys.SecretAccessKey))

		return 200, "We had sent credentials to your email, please check your email."
	}
	return 200, "Access keys generated"
}

func sendMail(user authCore.Accounts, access, secret string) {
	// map saving fname of user and aws access keys
	mapd := map[string]interface{}{
		"Name":      user.Name,
		"Accesskey": access,
		"Secretkey": secret,
		"Region":    "us-west-2",
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := mail.ReadToml("awskeys")

	// parse email template
	tmpl := mail.EmailTemplate(tmplPath, mapd)

	//now sending mail
	mail.SendMail(user.Email, subject, tmpl, images)
}
