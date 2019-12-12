package awsapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

// EC2Iface : ec2 interface
type EC2Iface interface {
	DescribeRunningEC2s() ([]EC2, error)
}

// EC2Instance : ec2 instance
type EC2Instance struct {
	client ec2iface.EC2API
}

// EC2 : required ec2 instance information
type EC2 struct {
	InstanceID       string
	PublicIPAddress  string
	PrivateIPAddress string
	InstanceType     string
	InstanceName     string
	AvailabilityZone string
}

// NewEC2Client : new ec2 client
func NewEC2Client(svc ec2iface.EC2API) EC2Iface {
	return &EC2Instance{
		client: svc,
	}
}

// DescribeRunningEC2s : get list of running ec2 instances
func (i *EC2Instance) DescribeRunningEC2s() ([]EC2, error) {
	// condition: running instance only
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}

	res, err := i.client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	e := []EC2{}
	for _, r := range res.Reservations {
		for _, i := range r.Instances {

			// public ip address
			if i.PublicIpAddress == nil {
				continue
			}
			publicIPAddress := *i.PublicIpAddress

			// tag:Name
			name := ""
			for _, t := range i.Tags {
				if *t.Key == "Name" {
					name = *t.Value
				}
			}

			// private ip address
			privateIPAddress := ""
			if i.PrivateIpAddress != nil {
				privateIPAddress = *i.PrivateIpAddress
			}

			e = append(e, EC2{
				InstanceID:       *i.InstanceId,
				InstanceType:     *i.InstanceType,
				PublicIPAddress:  publicIPAddress,
				PrivateIPAddress: privateIPAddress,
				InstanceName:     name,
				AvailabilityZone: *i.Placement.AvailabilityZone,
			})
		}
	}

	return e, nil
}

// FinderEC2 : find information of ec2 instance through fuzzyfinder
func FinderEC2(ec2List []EC2) (ec2 EC2, err error) {
	idx, err := fuzzyfinder.FindMulti(
		ec2List,
		func(i int) string {
			return fmt.Sprintf("[%s] %s (%s)",
				ec2List[i].InstanceName,
				ec2List[i].InstanceID,
				ec2List[i].InstanceType,
			)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf(
				"InstanceID: %s\ntag:Name: %s \nInstanceType: %s\nPublicIP: %s\nPrivateIP: %s",
				ec2List[i].InstanceID,
				ec2List[i].InstanceName,
				ec2List[i].InstanceType,
				ec2List[i].PublicIPAddress,
				ec2List[i].PrivateIPAddress,
			)
		}),
	)

	if err != nil {
		return ec2, err
	}

	for _, i := range idx {
		ec2 = ec2List[i]
	}

	return ec2, nil
}

// FinderUsername : find ssh username through fuzzyfinder
func FinderUsername(users []string) (user string, err error) {
	idx, err := fuzzyfinder.FindMulti(
		users,
		func(i int) string {
			return users[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return users[i]
		}),
	)

	if err != nil {
		return user, err
	}

	for _, i := range idx {
		user = users[i]
	}

	return user, nil
}
