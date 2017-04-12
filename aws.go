package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

var ELBWaitSleep = 3 * time.Second

const (
	StateInService    = "InService"
	StateOutOfService = "OutOfService"
)

type AWS struct {
	session *session.Session
}

func NewAWS() *AWS {
	return &AWS{
		session: session.Must(session.NewSession()),
	}
}

func (aws *AWS) InstanceID() (string, error) {
	meta := ec2metadata.New(aws.session)
	identity, err := meta.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}

	return identity.InstanceID, nil
}

func (aws *AWS) ELBRequireState(loadBalancerName, instance, desired string) error {
	svc := elb.New(aws.session)

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
		if *state.State != desired {
			return fmt.Errorf("Instance isn't '%s': %s", desired, *state.Description)
		}
	}

	return nil
}

func (aws *AWS) ELBWaitState(loadBalancerName, instance, desired string) {
	for {
		if aws.ELBRequireState(loadBalancerName, instance, desired) == nil {
			return
		}

		time.Sleep(ELBWaitSleep)
	}
}

func (aws *AWS) ELBRegisterInstance(loadBalancerName, instance string) error {
	svc := elb.New(aws.session)

	_, err := svc.RegisterInstancesWithLoadBalancer(
		&elb.RegisterInstancesWithLoadBalancerInput{
			LoadBalancerName: &loadBalancerName,
			Instances: []*elb.Instance{
				&elb.Instance{InstanceId: &instance},
			},
		})
	if err != nil {
		return err
	}

	aws.ELBWaitState(loadBalancerName, instance, StateInService)
	return nil
}

func (aws *AWS) ELBDeregisterInstance(loadBalancerName, instance string) error {
	svc := elb.New(aws.session)

	_, err := svc.DeregisterInstancesFromLoadBalancer(
		&elb.DeregisterInstancesFromLoadBalancerInput{
			LoadBalancerName: &loadBalancerName,
			Instances: []*elb.Instance{
				&elb.Instance{InstanceId: &instance},
			},
		})
	if err != nil {
		return err
	}

	aws.ELBWaitState(loadBalancerName, instance, StateOutOfService)
	return nil
}
