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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/srinivas-poturi-3/aws-controller/api/v1"
	"github.com/srinivas-poturi-3/aws-controller/internal/aws"
)

// VmReconciler reconciles a Vm object
type VmReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=aws.my.domain,resources=vms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws.my.domain,resources=vms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws.my.domain,resources=vms/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Vm object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *VmReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var vm v1.Vm
	err := r.Get(ctx, req.NamespacedName, &vm)
	if err != nil {
		log.Error(err, "unable to get CRD object")
		return ctrl.Result{}, err
	}

	// Extract secret reference
	var secretRef *v1.CredentialsSecret
	if vm.CredentialsSecretRef.Name != "" {
		secretRef = &vm.CredentialsSecretRef
		secretRef.Namespace = vm.Namespace
	}

	// Check if credentials secret is specified
	if secretRef == nil {
		log.Info("Credentials secret not specified in CRD. Skipping AWS actions.")
		// Update CRD status to reflect (optional)
		return ctrl.Result{}, nil
	}

	// Retrieve AWS credentials (implement getAWSCredentials function)
	secret, err := aws.GetAWSCredentials(ctx, r.Client, secretRef)
	if err != nil {
		log.Error(err, "failed to retrieve AWS credentials")
		return ctrl.Result{}, err
	}

	// Create AWS session
	awsSession, err := aws.GetSession(ctx, secret)
	if err != nil {
		log.Error(err, "unable to create AWS session")
		return ctrl.Result{}, err
	}

	// Use AWS SDK for VM management
	switch {
	case vm.ObjectMeta.DeletionTimestamp != nil:
		// Handle VM deletion
		err := aws.DeleteVM(awsSession, &vm)
		if err != nil {
			log.Error(err, "failed to delete VM")
			return ctrl.Result{}, err
		}
		// Update CRD status to reflect deletion
		vm.Status.Status = "Deleted"
		err = r.Status().Update(ctx, &vm)
		return ctrl.Result{}, err
	default:
		// Handle VM creation or update
		existingInstance, err := aws.GetExistingVM(awsSession, &vm)
		if err != nil {
			log.Error(err, "failed to check existing VM")
			return ctrl.Result{}, err
		}

		if existingInstance == nil {
			// Create VM
			err := aws.CreateVM(awsSession, &vm)
			if err != nil {
				log.Error(err, "failed to create VM")
				return ctrl.Result{}, err
			}
			// Update CRD status with VM information (e.g., instance ID)
			vm.Status.Status = "Running"
			err = r.Status().Update(ctx, &vm)
			if err != nil {
				log.Error(err, "failed to update CRD status")
				return ctrl.Result{}, err
			}
		}

	}
	log.Info("Reconciling Vm")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VmReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Vm{}).
		Complete(r)
}
