package awsapi

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
	"github.com/google/go-cmp/cmp"
)

type mockEC2InstanceConnectClient struct {
	ec2instanceconnectiface.EC2InstanceConnectAPI
}

func (m *mockEC2InstanceConnectClient) SendSSHPublicKey(input *ec2instanceconnect.SendSSHPublicKeyInput) (*ec2instanceconnect.SendSSHPublicKeyOutput, error) {
	return &ec2instanceconnect.SendSSHPublicKeyOutput{
		RequestId: aws.String("123456789"),
		Success:   aws.Bool(true),
	}, nil
}

// SendSSHPubKey : send ssh public key to using ec2 instance api
func TestSendSSHPubKey(t *testing.T) {
	mockSvc := &mockEC2InstanceConnectClient{}
	r, err := SendSSHPubKey(mockSvc, "ec2-user", "t3.micro", "", "ap-northeast-1a")
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(r, true); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}
