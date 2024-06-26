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

// NewAwsSession creates a new AWS session.
func NewAwsSession(sess *session.Session) *AwsSession {
	return &AwsSession{
		sess: sess,
	}
}

// CreateVM creates a new EC2 instance with given specs.
func (c *AwsSession) CreateVM(vm *v1.Vm) error {
	svc := ec2.New(c.sess)

	// Specifying instance details
	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String(vm.Spec.ImageId),
		InstanceType: aws.String(vm.Spec.InstanceType),
		MinCount:     aws.Int64(int64(vm.Spec.MinCount)),
		MaxCount:     aws.Int64(int64(vm.Spec.MaxCount)),
		KeyName:      aws.String(vm.Spec.KeyName),
		SubnetId:     aws.String(vm.Spec.SubnetId),
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
		vm.Status.InstanceStatus = append(vm.Status.InstanceStatus, v1.InstanceStatus{
			InstanceId: *runOutput.Instances[i].InstanceId,
			State:      *runOutput.Instances[i].State.Name,
		})
	}

	return nil
}

// GetExistingVM gets the existing EC2 instance details.
func (c *AwsSession) GetExistingVM(vm *v1.Vm) error {
	svc := ec2.New(c.sess)

	instancesIds := make([]string, len(vm.Status.InstanceStatus))
	for i, id := range vm.Status.InstanceStatus {
		instancesIds[i] = id.InstanceId
	}
	input := &ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice(instancesIds)}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		fmt.Printf("Error describing EC2 instance: %v\n", err)
		return err
	}
	vm.Status.InstanceStatus = []v1.InstanceStatus{}
	// Store details in VM status
	for i := range result.Reservations {
		for j := range result.Reservations[i].Instances {
			instance := v1.InstanceStatus{
				InstanceId: *result.Reservations[i].Instances[j].InstanceId,
				State:      *result.Reservations[i].Instances[j].State.Name,
			}
			if result.Reservations[i].Instances[j].PrivateIpAddress != nil {
				instance.PrivateIpAddresses = *result.Reservations[i].Instances[j].PrivateIpAddress
			}
			if result.Reservations[i].Instances[j].PublicIpAddress != nil {
				instance.PublicIpAddresses = *result.Reservations[i].Instances[j].PublicIpAddress
			}
			vm.Status.InstanceStatus = append(vm.Status.InstanceStatus, instance)
		}
	}
	return nil
}

// DeleteVM deletes the existing EC2 instance.
func (c *AwsSession) DeleteVM(vm *v1.Vm) error {
	svc := ec2.New(c.sess)

	instancesIds := make([]string, len(vm.Status.InstanceStatus))
	for i, id := range vm.Status.InstanceStatus {
		instancesIds[i] = id.InstanceId
	}
	// Specifying instance ID for termination
	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: aws.StringSlice(instancesIds),
	}

	_, err := svc.TerminateInstances(terminateInput)
	if err != nil {
		return fmt.Errorf("error terminating EC2 instance: %v", err)

	}

	fmt.Printf("Terminated EC2 instance with ID: %s\n", instancesIds)

	return nil
}
