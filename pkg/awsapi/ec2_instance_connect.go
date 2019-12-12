package awsapi

import (
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
)

// EC2InstanceConnectIface : ec2 instance connect interface
type EC2InstanceConnectIface interface {
	SendSSHPubKey(ec2instanceconnect.SendSSHPublicKeyInput) (bool, error)
}

// EC2InstanceConnectInstance : ec2 instance connect instance
type EC2InstanceConnectInstance struct {
	client ec2instanceconnectiface.EC2InstanceConnectAPI
}

// NewEC2InstanceConnectClient : new ec2 instance connect client
func NewEC2InstanceConnectClient(svc ec2instanceconnectiface.EC2InstanceConnectAPI) EC2InstanceConnectIface {
	return &EC2InstanceConnectInstance{
		client: svc,
	}
}

// SendSSHPubKey : send ssh public key to using ec2 instance api
func (i *EC2InstanceConnectInstance) SendSSHPubKey(p ec2instanceconnect.SendSSHPublicKeyInput) (bool, error) {
	r, err := i.client.SendSSHPublicKey(&p)
	if err != nil {
		return false, err
	}
	return *r.Success, nil
}
