package awsapi

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/google/go-cmp/cmp"
	"github.com/kenzo0107/omssh/pkg/utility"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/nsf/termbox-go"
)

var (
	testEC2Instances = []EC2Instance{
		EC2Instance{
			InstanceID:       "i-aaaaaa",
			PublicIPAddress:  "12.34.56.01",
			PrivateIPAddress: "192.168.10.1",
			InstanceType:     "t3.micro",
			AvailabilityZone: "ap-northeast-1a",
			InstanceName:     "hoge",
		},
		EC2Instance{
			InstanceID:       "i-bbbbbb",
			PublicIPAddress:  "12.34.56.02",
			PrivateIPAddress: "192.168.10.2",
			InstanceType:     "t3.small",
			AvailabilityZone: "ap-northeast-1c",
			InstanceName:     "moge",
		},
	}
)

type mockEC2Client struct {
	ec2iface.EC2API
}

func (m *mockEC2Client) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:       aws.String("i-aaaaaa"),
						InstanceType:     aws.String("t3.micro"),
						PublicIpAddress:  aws.String("12.34.56.01"),
						PrivateIpAddress: aws.String("192.168.10.1"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1a"),
						},
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("hoge"),
							},
						},
					},
				},
			},
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:       aws.String("i-bbbbbb"),
						InstanceType:     aws.String("t3.small"),
						PublicIpAddress:  aws.String("12.34.56.02"),
						PrivateIpAddress: aws.String("192.168.10.2"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1c"),
						},
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("moge"),
							},
						},
					},
				},
			},
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:       aws.String("i-cccccc"),
						InstanceType:     aws.String("t3.medium"),
						PrivateIpAddress: aws.String("192.168.10.3"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1c"),
						},
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("foo"),
							},
						},
					},
				},
			},
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:       aws.String("i-dddddd"),
						InstanceType:     aws.String("t3.large"),
						PrivateIpAddress: aws.String("192.168.10.4"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1c"),
						},
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("baz"),
							},
						},
					},
				},
			},
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:       aws.String("i-eeeeee"),
						InstanceType:     aws.String("t3.xlarge"),
						PrivateIpAddress: aws.String("192.168.10.5"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1c"),
						},
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("bar"),
							},
						},
					},
				},
			},
		},
	}, nil
}

type mockEC2ClientFaild struct {
	ec2iface.EC2API
}

func (m *mockEC2ClientFaild) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{},
	}, errors.New("Not Found")
}

func TestDescribeRunningEC2Instances(t *testing.T) {
	mockSvc := &mockEC2Client{}
	runningEC2Instances, err := DescribeRunningEC2Instances(mockSvc)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(testEC2Instances, runningEC2Instances); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}

func TestDescribeNotFoundRunningEC2Instances(t *testing.T) {
	mockEC2ClientFaild := &mockEC2ClientFaild{}
	runningEC2Instances, err := DescribeRunningEC2Instances(mockEC2ClientFaild)
	if err == nil {
		t.Error(err)
	}

	if diff := cmp.Diff([]EC2Instance(nil), runningEC2Instances); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}

func finderEC2Testing(t *testing.T, types string, tests []EC2Instance, expectedEC2 EC2Instance) {
	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)

	term.SetEvents(append(
		utility.TermboxKeys(types),
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

	actualEC2, err := FinderEC2Instance(tests)
	if err != nil {
		t.Error("cannot get profile")
	}
	if diff := cmp.Diff(expectedEC2, actualEC2); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}
}

func TestFinderEC2Instance(t *testing.T) {
	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)

	for _, testcase := range []struct {
		name string
		call func(t *testing.T)
	}{
		{
			"type i-a on terminal",
			func(t *testing.T) {
				finderEC2Testing(
					t,
					"i-a",
					testEC2Instances,
					EC2Instance{
						InstanceID:       "i-aaaaaa",
						PublicIPAddress:  "12.34.56.01",
						PrivateIPAddress: "192.168.10.1",
						InstanceType:     "t3.micro",
						AvailabilityZone: "ap-northeast-1a",
						InstanceName:     "hoge",
					},
				)
			},
		},
		{
			"type i-b on terminal",
			func(t *testing.T) {
				finderEC2Testing(t,
					"i-b",
					testEC2Instances,
					EC2Instance{
						InstanceID:       "i-bbbbbb",
						PublicIPAddress:  "12.34.56.02",
						PrivateIPAddress: "192.168.10.2",
						InstanceType:     "t3.small",
						AvailabilityZone: "ap-northeast-1c",
						InstanceName:     "moge",
					},
				)
			},
		},
		{
			"type foo - Not found Instance name on terminal",
			func(t *testing.T) {
				types := "foo"
				term := fuzzyfinder.UseMockedTerminal()
				term.SetSize(60, 10)

				term.SetEvents(append(
					utility.TermboxKeys(types),
					termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

				actual, err := FinderEC2Instance(testEC2Instances)
				if err == nil {
					t.Errorf("wrong result: err is not nil\n %s", err)
				}
				t.Logf("actual: %#v", actual)
			},
		},
	} {
		t.Run(testcase.name, testcase.call)
	}
}

func TestFinderUsername(t *testing.T) {
	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)

	for _, testcase := range []struct {
		name string
		call func(t *testing.T)
	}{
		{
			"type ubuntu on terminal",
			func(t *testing.T) {
				username := "ubuntu"
				expected := username
				term.SetEvents(append(
					utility.TermboxKeys(username),
					termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

				usernames := []string{"ubuntu", "ec2-user"}
				actual, err := FinderUsername(usernames)
				if err != nil {
					t.Error(err)
				}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
			},
		},
		{
			"type ec2-user on terminal",
			func(t *testing.T) {
				username := "ec2-user"
				expected := username
				term.SetEvents(append(
					utility.TermboxKeys(username),
					termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

				usernames := []string{"ubuntu", "ec2-user"}
				actual, err := FinderUsername(usernames)
				if err != nil {
					t.Error(err)
				}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
			},
		},
		{
			"type hoge on terminal",
			func(t *testing.T) {
				username := "hoge"
				expected := ""
				term.SetEvents(append(
					utility.TermboxKeys(username),
					termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

				usernames := []string{"ubuntu", "ec2-user"}
				actual, err := FinderUsername(usernames)
				if err == nil {
					t.Errorf("wrong result: err is not nil\n %s", err)
				}
				if diff := cmp.Diff(expected, actual); diff != "" {
					t.Errorf("wrong result: \n%s", diff)
				}
			},
		},
	} {
		t.Run(testcase.name, testcase.call)
	}
}
