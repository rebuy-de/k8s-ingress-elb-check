package cmd

import (
	"fmt"

	"github.com/rebuy-de/k8s-ingress-elb-check/pkg/aws"
	"github.com/spf13/cobra"
)

func NewDeregisterCommand(loadBalancerName *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister",
		Short: "deregister an instance from an ELB",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		api := aws.New()

		instanceID, err := api.InstanceID()
		must(err)

		err = api.ELBDeregisterInstance(*loadBalancerName, instanceID)
		must(err)

		fmt.Printf("Deregistered '%s' from LoadBalancer '%s'.\n", instanceID, *loadBalancerName)
	}

	return cmd
}
