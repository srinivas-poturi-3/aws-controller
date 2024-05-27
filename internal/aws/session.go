package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsv1 "github.com/srinivas-poturi-3/aws-controller/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Creds struct {
	accessID  string
	accessKey string
	region    string
}

const (
	defaultRegion      = "us-east-1"
	accessID           = "access_id"
	accessKey          = "access_key"
	awsAccessKeyID     = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
)

func GetSession(ctx context.Context, creds Creds) (*AwsSession, error) {
	os.Setenv(awsAccessKeyID, creds.accessID)
	os.Setenv(awsSecretAccessKey, creds.accessKey)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(creds.region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &AwsSession{sess: sess}, err
}

func GetAWSCredentials(ctx context.Context, k8sClient client.Client, secretRef *awsv1.CredentialsSecret) (Creds, error) {
	// Get the secret object based on the reference
	secret := &corev1.Secret{}
	secretKey := types.NamespacedName{Name: secretRef.Name, Namespace: secretRef.Namespace}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: secretKey.Name, Namespace: secretKey.Namespace}, secret)
	if err != nil {
		return Creds{}, fmt.Errorf("failed to get secret %s: %w", secretKey, err)
	}

	creds := Creds{}
	// Check if required data keys exist
	if _, ok := secret.Data[accessID]; !ok {
		return Creds{}, fmt.Errorf("secret %s missing accessID data", secretKey)
	}
	if _, ok := secret.Data[accessKey]; !ok {
		return Creds{}, fmt.Errorf("secret %s missing accessKey data", secretKey)
	}

	if secretRef.Region == "" {
		creds.region = defaultRegion
	} else {
		creds.region = secretRef.Region
	}

	creds.accessID = string(secret.Data[accessID])
	creds.accessKey = string(secret.Data[accessKey])

	return creds, nil
}
