# k8s-ingress-elb-check

Container which checks, if an ingress replica is InService on the ELB.

## Configuration

AWS settings are configured via environemnt variables. These are the same like in the AWS CLI.

```
AWS_PROFILE=production AWS_REGION=eu-west-1 k8s-ingress-elb-check
```
