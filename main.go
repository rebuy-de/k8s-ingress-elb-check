package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

const StateInService = "InService"

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <load-balancer-name>", os.Args[0])
		os.Exit(2)
	}

	loadBalancerName := os.Args[1]

	sess := session.Must(session.NewSession())

	instanceID, err := getInstanceID(sess)
	err = checkInstance(sess, loadBalancerName, instanceID)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Instance '%s' of LoadBalancer '%s' is InService.\n", instanceID, loadBalancerName)
}

func getInstanceID(sess *session.Session) (string, error) {
	meta := ec2metadata.New(sess)
	identity, err := meta.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}

	return identity.InstanceID, nil
}

func checkInstance(sess *session.Session, loadBalancerName, instance string) error {
	svc := elb.New(sess)

	result, err := svc.DescribeInstanceHealth(&elb.DescribeInstanceHealthInput{
		LoadBalancerName: &loadBalancerName,
		Instances: []*elb.Instance{
			&elb.Instance{InstanceId: &instance},
		},
	})

	if err != nil {
		return err
	}

	for _, state := range result.InstanceStates {
		if *state.State != StateInService {
			return fmt.Errorf("Instance isn't ready: %s", *state.Description)
		}
	}

	return nil
}
