package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	v1 "github.com/srinivas-poturi-3/aws-controller/api/v1"
)

type AwsSession struct {
	sess *session.Session
}

func NewAwsSession(sess *session.Session) *AwsSession {
	return &AwsSession{
		sess: sess,
	}
}

func (c *AwsSession) CreateVM(vm *v1.Vm) error {
	svc := ec2.New(c.sess)

	// Specify instance details
	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(vm.Spec.ImageId),
		InstanceType: aws.String(vm.Spec.InstanceType),
		MinCount:     aws.Int64(int64(vm.Spec.MinCount)),
		MaxCount:     aws.Int64(int64(vm.Spec.MaxCount)),
		KeyName:      aws.String(vm.Spec.KeyName),
		SubnetId:     aws.String(vm.Spec.SubnetId),
		// Add other parameters (e.g., security groups, key pair, etc.)
	}

	runOutput, err := svc.RunInstances(runInput)
	if err != nil {
		fmt.Printf("Error creating EC2 instance: %v\n", err)
		return err
	}

	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runOutput.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(vm.Spec.Name),
			},
		},
	})
	if errtag != nil {
		fmt.Println("Could not create tags for instance", runOutput.Instances[0].InstanceId, errtag)
		return nil
	}

	// Store instance ID in VM status
	for i := range runOutput.Instances {
		vm.Status.Instance = append(vm.Status.Instance, *runOutput.Instances[i].InstanceId)
	}

	fmt.Printf("Created EC2 instance with ID: %s\n", vm.Status.Instance)
	return nil
}

func (c *AwsSession) GetExistingVM(vm *v1.Vm) (*ec2.DescribeInstancesOutput, error) {
	svc := ec2.New(c.sess)

	input := &ec2.DescribeInstancesInput{}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Printf("Error describing EC2 instance: %v\n", err)
		return nil, err
	}

	return result, nil
}

func (c *AwsSession) DeleteVM(vm *v1.Vm) error {
	svc := ec2.New(c.sess)

	// Specify instance ID for termination
	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice(vm.Status.Instance),
	}

	_, err := svc.TerminateInstances(terminateInput)
	if err != nil {
		return fmt.Errorf("error terminating EC2 instance: %v", err)

	}

	fmt.Printf("Terminated EC2 instance with ID: %s\n", vm.Status.Instance)

	return nil
}
