package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

const StateInService = "InService"

func main() {
	sess := session.Must(session.NewSession())

	err := checkInstance(sess, "public-ingress", "i-0c03fae448fa9e17b")

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Instance '%s' of LoadBalancer '%s' is InService.\n", "foo", "bar")
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
