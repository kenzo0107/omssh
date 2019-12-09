package awsapi

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
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
	testEC2Instances = []EC2Info{
		{
			InstanceID:       "i-aaaaaa",
			PublicIPAddress:  "12.34.56.01",
			InstanceType:     "t3.micro",
			AvailabilityZone: "ap-northeast-1a",
		},
		{
			InstanceID:       "i-bbbbbb",
			PublicIPAddress:  "12.34.56.02",
			InstanceType:     "t3.small",
			AvailabilityZone: "ap-northeast-1c",
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
						InstanceId:      aws.String("i-aaaaaa"),
						InstanceType:    aws.String("t3.micro"),
						PublicIpAddress: aws.String("12.34.56.01"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1a"),
						},
					},
				},
			},
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:      aws.String("i-bbbbbb"),
						InstanceType:    aws.String("t3.small"),
						PublicIpAddress: aws.String("12.34.56.02"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1c"),
						},
					},
				},
			},
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:   aws.String("i-cccccc"),
						InstanceType: aws.String("t3.medium"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1c"),
						},
					},
				},
			},
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:   aws.String("i-dddddd"),
						InstanceType: aws.String("t3.large"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1c"),
						},
					},
				},
			},
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:   aws.String("i-eeeeee"),
						InstanceType: aws.String("t3.xlarge"),
						Placement: &ec2.Placement{
							AvailabilityZone: aws.String("ap-northeast-1c"),
						},
					},
				},
			},
		},
	}, nil
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

func finderEC2Testing(t *testing.T, types string, expectedEC2 EC2Info) {
	term := fuzzyfinder.UseMockedTerminal()
	term.SetSize(60, 10)

	term.SetEvents(append(
		utility.TermboxKeys(types),
		termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter})...)

	actualEC2, err := FinderEC2Instance(testEC2Instances)
	if err != nil {
		t.Error("cannot get profile")
	}
	if diff := cmp.Diff(expectedEC2, actualEC2); diff != "" {
		t.Errorf("wrong result: \n%s", diff)
	}

	actual := term.GetResult()

	g := fmt.Sprintf("finder_ec2_%s_ui.golden", types)
	fname := filepath.Join("..", "..", "testdata", g)
	// ioutil.WriteFile(fname, []byte(actual), 0644)
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("failed to load a golden file: %s", err)
	}
	expected := string(b)
	if runtime.GOOS == "windows" {
		expected = strings.Replace(expected, "\r\n", "\n", -1)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
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
				expected := EC2Info{
					InstanceID:       "i-aaaaaa",
					PublicIPAddress:  "12.34.56.01",
					InstanceType:     "t3.micro",
					AvailabilityZone: "ap-northeast-1a",
				}
				finderEC2Testing(t, "i-a", expected)
			},
		},
		{
			"type i-b on terminal",
			func(t *testing.T) {
				expected := EC2Info{
					InstanceID:       "i-bbbbbb",
					PublicIPAddress:  "12.34.56.02",
					InstanceType:     "t3.small",
					AvailabilityZone: "ap-northeast-1c",
				}
				finderEC2Testing(t, "i-b", expected)
			},
		},
	} {
		t.Run(testcase.name, testcase.call)
	}
}
