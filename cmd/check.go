package cmd

import (
	"fmt"

	"github.com/rebuy-de/k8s-ingress-elb-check/pkg/aws"
	"github.com/spf13/cobra"
)

func NewCheckCommand(loadBalancerName *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "check if the instance is registered to the container",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		api := aws.New()

		instanceID, err := api.InstanceID()
		must(err)

		err = api.ELBRequireState(*loadBalancerName, instanceID, aws.StateInService)
		if err != nil {
			fmt.Printf("Instance '%s' of LoadBalancer '%s' is not InService: %v\n",
				instanceID, *loadBalancerName, err)
		}

		fmt.Printf("Instance '%s' of LoadBalancer '%s' is InService.\n",
			instanceID, *loadBalancerName)
	}

	return cmd
}
