package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	v1 "github.com/srinivas-poturi-3/aws-controller/api/v1"
)

func CreateVM(sess *session.Session, vm *v1.Vm) error {
	svc := ec2.New(sess)

	// Specify instance details
	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(vm.Spec.ImageId),
		InstanceType: aws.String(vm.Spec.InstanceType),
		MinCount:     aws.Int64(int64(vm.Spec.MinCount)),
		MaxCount:     aws.Int64(int64(vm.Spec.MaxCount)),
		KeyName:      aws.String(vm.Spec.KeyName),
		SubnetId:     aws.String(vm.Spec.Subnet),
		// Add other parameters (e.g., security groups, key pair, etc.)
	}

	runOutput, err := svc.RunInstances(runInput)
	if err != nil {
		fmt.Printf("Error creating EC2 instance: %v\n", err)
		return err
	}

	// Store instance ID in VM status
	vm.Status.InstanceId = *runOutput.Instances[0].InstanceId
	fmt.Printf("Created EC2 instance with ID: %s\n", vm.Status.InstanceId)
	return nil
}

func GetExistingVM(sess *session.Session, vm *v1.Vm) (*ec2.DescribeInstancesOutput, error) {
	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(vm.Status.InstanceId)},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Printf("Error describing EC2 instance: %v\n", err)
		return nil, err
	}

	return result, nil
}

func DeleteVM(sess *session.Session, vm *v1.Vm) error {
	svc := ec2.New(sess)

	// Specify instance ID for termination
	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(vm.Status.InstanceId)},
	}

	_, err := svc.TerminateInstances(terminateInput)
	if err != nil {
		return fmt.Errorf("error terminating EC2 instance: %v", err)

	}

	fmt.Printf("Terminated EC2 instance with ID: %s\n", vm.Status.InstanceId)

	return nil
}
