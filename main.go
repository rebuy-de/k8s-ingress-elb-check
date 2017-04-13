package main

import (
	"fmt"
	"os"

	"github.com/rebuy-de/k8s-ingress-elb-check/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
