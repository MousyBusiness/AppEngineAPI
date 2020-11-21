package secrets

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/cloudresourcemanager/v1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
	"os"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
)

func getProjectNumber() (int64, error) {
	crmService, err := cloudresourcemanager.NewService(context.Background())
	if err != nil {
		return 0, errors.Wrap(err, "failed to get cloud resource manager service")
	}

	project, err := crmService.Projects.Get(projectID).Do()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get project number from ID")
	}

	log.Println("project number for project:", project.ProjectNumber)

	return project.ProjectNumber, nil
}

func GetSecret(secretId string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to setup secret manager client")
	}

	projectNumber, err := getProjectNumber()
	if err != nil {
		return "", err
	}

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%d/secrets/%s/versions/%d", projectNumber, secretId, 1),
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return "", errors.Wrap(err, "couldnt get secret")
	}

	return string(result.Payload.Data), nil
}
