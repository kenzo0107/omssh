package awsapi

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
)

// EC2InstanceConnectClient : ec2 instance connect client
type EC2InstanceConnectClient struct {
	sess *session.Session
	svc  *ec2instanceconnect.EC2InstanceConnect
}

// NewEC2InstanceConnect : new ec2 instance connect client
func NewEC2InstanceConnect(sess *session.Session) *EC2InstanceConnectClient {
	svc := ec2instanceconnect.New(sess)
	return &EC2InstanceConnectClient{sess, svc}
}

// SendSSHPubKey : send ssh public key to using ec2 instance api
func (d *EC2InstanceConnectClient) SendSSHPubKey(user, instanceID, publicKey, availabilityZone string) error {
	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(availabilityZone),
		InstanceId:       aws.String(instanceID),
		InstanceOSUser:   aws.String(user),
		SSHPublicKey:     aws.String(publicKey),
	}
	_, err := d.svc.SendSSHPublicKey(input)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
