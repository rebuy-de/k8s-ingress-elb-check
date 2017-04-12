package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("USAGE: %s <check|register|deregister> <load-balancer-name>\n", os.Args[0])
	os.Exit(2)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustPretty(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	command := os.Args[1]
	loadBalancerName := os.Args[2]

	switch command {
	case "check":
		checkCommand(loadBalancerName)
	case "register":
		registerCommand(loadBalancerName)
	case "deregister":
		deregisterCommand(loadBalancerName)
	default:
		usage()
	}
}

func checkCommand(loadBalancerName string) {
	aws := NewAWS()

	instanceID, err := aws.InstanceID()
	must(err)

	err = aws.ELBRequireState(loadBalancerName, instanceID, StateInService)
	mustPretty(err)

	fmt.Printf("Instance '%s' of LoadBalancer '%s' is InService.\n", instanceID, loadBalancerName)
}

func registerCommand(loadBalancerName string) {
	aws := NewAWS()

	instanceID, err := aws.InstanceID()
	must(err)

	err = aws.ELBRegisterInstance(loadBalancerName, instanceID)
	must(err)

	fmt.Printf("Registered '%s' to LoadBalancer '%s'.\n", instanceID, loadBalancerName)
}

func deregisterCommand(loadBalancerName string) {
	aws := NewAWS()

	instanceID, err := aws.InstanceID()
	must(err)

	err = aws.ELBDeregisterInstance(loadBalancerName, instanceID)
	must(err)

	fmt.Printf("Deregistered '%s' from LoadBalancer '%s'.\n", instanceID, loadBalancerName)
}
