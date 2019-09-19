package omssh

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
)

// EC2Client : ec2 client
type EC2Client struct {
	sess *session.Session
	svc  *ec2.EC2
}

// EC2Info : required ec2 instance information
type EC2Info struct {
	InstanceID       string
	PublicIPAddress  string
	PrivateIPAddress string
	InstanceType     string
	InstanceName     string
	AvailabilityZone string
}

// NewEC2 : new ec2 client
func NewEC2(sess *session.Session) *EC2Client {
	svc := ec2.New(sess)
	return &EC2Client{sess, svc}
}

// GetEC2List : get list of ec2 instances
func (d *EC2Client) GetEC2List() ([]EC2Info, error) {
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

	res, err := d.svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	e := []EC2Info{}
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

			e = append(e, EC2Info{
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

// FinderEC2Info : find information of ec2 instance through fuzzyfinder
func FinderEC2Info(ec2List []EC2Info) (ec2Info EC2Info, err error) {
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
		log.Fatal(err)
		return ec2Info, err
	}

	for _, i := range idx {
		ec2Info = ec2List[i]
	}

	return ec2Info, nil
}

// FinderUsername : find ssh username through fuzzyfinder
func FinderUsername() (user string, err error) {
	users := []string{
		"ubuntu",
		"ec2-user",
	}
	idx, err := fuzzyfinder.FindMulti(
		users,
		func(i int) string {
			return fmt.Sprintf("%s",
				users[i],
			)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf(
				"%s",
				users[i],
			)
		}),
	)

	if err != nil {
		log.Fatal(err)
		return user, err
	}

	for _, i := range idx {
		user = users[i]
	}

	return user, nil
}
