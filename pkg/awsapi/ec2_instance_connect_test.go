package awsapi

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
	"github.com/google/go-cmp/cmp"
)

type mockEC2InstanceConnectiface struct {
	ec2instanceconnectiface.EC2InstanceConnectAPI

	Resp  ec2instanceconnect.SendSSHPublicKeyOutput
	Error error
}

func (m *mockEC2InstanceConnectiface) SendSSHPublicKey(p *ec2instanceconnect.SendSSHPublicKeyInput) (*ec2instanceconnect.SendSSHPublicKeyOutput, error) {
	return &m.Resp, m.Error
}

func TestSendSSHPubKey(t *testing.T) {
	m := NewEC2InstanceConnectClient(&mockEC2InstanceConnectiface{
		Resp: ec2instanceconnect.SendSSHPublicKeyOutput{
			RequestId: aws.String("1234567890"),
			Success:   aws.Bool(true),
		},
		Error: nil,
	})
	r, err := m.SendSSHPubKey(ec2instanceconnect.SendSSHPublicKeyInput{})
	if err != nil {
		t.Error("wrong result \n err is not nil")
	}
	if diff := cmp.Diff(r, true); diff != "" {
		t.Errorf("wrong result \n%s", diff)
	}
}

func TestSendSSHPubKeyWithError(t *testing.T) {
	m := NewEC2InstanceConnectClient(&mockEC2InstanceConnectiface{
		Resp: ec2instanceconnect.SendSSHPublicKeyOutput{
			RequestId: aws.String("1234567890"),
			Success:   aws.Bool(false),
		},
		Error: errors.New("error occured"),
	})
	r, err := m.SendSSHPubKey(ec2instanceconnect.SendSSHPublicKeyInput{})
	if err == nil {
		t.Error("wrong result \n err is nil")
	}
	if diff := cmp.Diff(r, false); diff != "" {
		t.Errorf("wrong result \n%s", diff)
	}
}
