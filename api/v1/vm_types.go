/*
Copyright 2024 Srinivas.poturi.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VmSpec defines the desired state of Vm
type VmSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Name               string   `json:"name,omitempty"`
	MaxCount           int      `json:"maxCount,omitempty"`
	MinCount           int      `json:"minCount,omitempty"`
	ImageId            string   `json:"imageId,omitempty"`
	InstanceType       string   `json:"instanceType,omitempty"`
	KeyName            string   `json:"keyName,omitempty"`
	SecurityGroupIds   []string `json:"securityGroupIds,omitempty"`
	SubnetId           string   `json:"subnetId,omitempty"`
	UserData           string   `json:"userData,omitempty"`
	DryRun             bool     `json:"dryRun,omitempty"`
	IamInstanceProfile string   `json:"iamInstanceProfile,omitempty"`
	// NetworkInterface              []ec2.InstanceNetworkInterfaceSpecification `json:"NetworkInterface,omitempty"`
	// BlockDeviceMapping            ec2.BlockDeviceMapping                      `json:"BlockDeviceMapping,omitempty"`
	// MetadataOptions               ec2.InstanceMetadataOptionsRequest          `json:"MetadataOptions,omitempty"`
	// PrivateDnsNameOptionsOnLaunch ec2.PrivateDnsNameOptionsOnLaunch           `json:"PrivateDnsNameOptionsOnLaunch,omitempty"`
}

// VmStatus defines the observed state of Vm
type VmStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status         string           `json:"status,omitempty"`
	Error          string           `json:"error,omitempty"`
	InstanceStatus []InstanceStatus `json:"instanceStatus,omitempty"`
}

type InstanceStatus struct {
	InstanceId         string `json:"instanceId,omitempty"`
	State              string `json:"state,omitempty"`
	PrivateIpAddresses string `json:"privateIpAddresses,omitempty"`
	PublicIpAddresses  string `json:"publicIpAddresses,omitempty"`
}

// CredentialsSecret defines the reference to the secret containing AWS credentials
type CredentialsSecret struct {
	// Name of the secret containing credentials
	Name string `json:"name"`

	// Namespace where the secret resides
	Namespace string `json:"namespace,omitempty"`

	// Region of the AWS account
	Region string `json:"region,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Vm is the Schema for the vms API
type Vm struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VmSpec   `json:"spec,omitempty"`
	Status VmStatus `json:"status,omitempty"`
	// CredentialsSecretRef specifies the reference to the secret containing AWS credentials (optional)
	CredentialsSecretRef CredentialsSecret `json:"credentialsSecretRef,omitempty"`
}

//+kubebuilder:object:root=true

// VmList contains a list of Vm
type VmList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Vm `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Vm{}, &VmList{})
}
