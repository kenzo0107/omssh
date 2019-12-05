package awsapi

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
)

// SendSSHPubKey : send ssh public key to using ec2 instance api
func SendSSHPubKey(svc ec2instanceconnectiface.EC2InstanceConnectAPI, user, instanceID, publicKey, availabilityZone string) (bool, error) {
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(availabilityZone),
		InstanceId:       aws.String(instanceID),
		InstanceOSUser:   aws.String(user),
		SSHPublicKey:     aws.String(publicKey),
	}
	r, err := svc.SendSSHPublicKey(input)
	if err != nil {
		log.Fatal(err)
		return *r.Success, err
	}
	return *r.Success, nil
}
