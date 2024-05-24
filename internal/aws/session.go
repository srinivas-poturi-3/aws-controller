package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsv1 "github.com/srinivas-poturi-3/aws-controller/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Creds struct {
	accessID    string
	accessKey   string
	accessToken string
	region      string
}

func GetSession(ctx context.Context, creds Creds) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(creds.region),
		Credentials: credentials.NewStaticCredentials(creds.accessID, creds.accessKey, creds.accessToken),
	})
	return sess, err
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
	if key, ok := secret.Data["accessID"]; !ok {
		creds.accessID = string(key)
		return Creds{}, fmt.Errorf("secret %s missing accessID data", secretKey)
	}
	if key, ok := secret.Data["accessKey"]; !ok {
		creds.accessKey = string(key)
		return Creds{}, fmt.Errorf("secret %s missing accessKey data", secretKey)
	}
	if key, ok := secret.Data["tokenKey"]; !ok {
		creds.accessToken = string(key)
		return Creds{}, fmt.Errorf("secret %s missing secretKey data", secretKey)
	}
	if key, ok := secret.Data["region"]; !ok {
		creds.region = string(key)
		return Creds{}, fmt.Errorf("secret %s missing region data", secretKey)
	}

	return creds, nil
}
