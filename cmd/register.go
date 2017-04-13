package cmd

import (
	"fmt"

	"github.com/rebuy-de/k8s-ingress-elb-check/pkg/aws"
	"github.com/spf13/cobra"
)

func NewRegisterCommand(loadBalancerName *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register an instance to an ELB",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		api := aws.New()

		instanceID, err := api.InstanceID()
		must(err)

		err = api.ELBRegisterInstance(*loadBalancerName, instanceID)
		must(err)

		fmt.Printf("Registered '%s' to LoadBalancer '%s'.\n", instanceID, *loadBalancerName)
	}

	return cmd
}
