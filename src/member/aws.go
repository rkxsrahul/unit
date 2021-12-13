package member

import (
	"log"
	"strings"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/aws"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

func awsWork(acc database.Accounts, wsID string) error {

	//create iam user
	_, err := aws.CreateIamUser(acc.Userid)
	if err != nil && !strings.Contains(err.Error(), "status code: 409") {
		log.Println(err)
		return err
	}

	//attach policy
	perr := aws.AttachPolicy("sl"+acc.Userid, "Enterprise", wsID)
	if perr != nil && !strings.Contains(perr.Error(), "status code: 409") {
		log.Println(perr)
		return perr
	}

	// create bucket
	berr := aws.CreateBucket(wsID + "-sl/" + "private-dataset-sl" + acc.Userid + "/")
	if berr != nil && !strings.Contains(berr.Error(), "status code: 409") {
		log.Println(berr)
		return berr
	}

	return nil
}

func awsDeleteWork(acc database.Accounts, wsID string) error {
	// delete bucket of user
	berr := aws.DeleteObject(wsID+"-sl", "private-dataset-sl"+acc.Userid+"/")
	if berr != nil && !strings.Contains(berr.Error(), "status code: 409") {
		log.Println(berr)
		return berr
	}

	//deattach enterprise policuy
	db := config.DB

	policies := []database.Policy{}
	db.Where("company=? AND p_type=?", wsID, "Enterprise").Find(&policies)
	for i := 0; i < len(policies); i++ {
		err := aws.DetachPolicy("sl"+acc.Userid, policies[i].Arn)
		if err != nil && !strings.Contains(err.Error(), "status code: 404") {
			log.Println(err)
			return err
		}
	}
	return nil
}
