package config

import (
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"net"
	"os"
	"paymentService/model"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/secretmanager/apiv1"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newSecretClient(ctx context.Context) (*secretmanager.Client, error) {
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath != "" {
		return secretmanager.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	}

	return secretmanager.NewClient(ctx)
}

func AccessSecret(ctx context.Context, secretName string) (string, error) {
	client, err := newSecretClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to setup secretmanager client: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret version: %w", err)
	}

	secretData := string(result.Payload.Data)
	return secretData, nil
}

func buildSecretPath(projectID, secretName string) string {
	return fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName)
}

func ConnectToDB() (*gorm.DB, error) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("PROJECT_ID environment variable is not set")
	}

	secretdbUser := buildSecretPath(projectID, "DB_USER")
	secretdbPwd := buildSecretPath(projectID, "DB_PASSWORD")
	secretdbName := buildSecretPath(projectID, "DB_NAME")
	secretinstanceConnectionName := buildSecretPath(projectID, "INSTANCE_CONNECTION_NAME")

	dbUser, err := AccessSecret(ctx, secretdbUser)
	dbPwd, err := AccessSecret(ctx, secretdbPwd)
	dbName, err := AccessSecret(ctx, secretdbName)
	instanceConnectionName, err := AccessSecret(ctx, secretinstanceConnectionName)
	usePrivate := os.Getenv("PRIVATE_IP")

	d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithLazyRefresh())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}

	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}
	mysqlDriver.RegisterDialContext("cloudsqlconn", func(ctx context.Context, addr string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName, opts...)
	})

	dsn := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}

	db.AutoMigrate(&model.Payment{})
	return db, nil
}
