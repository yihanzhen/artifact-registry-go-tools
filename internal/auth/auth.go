package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2/google"
)

// applicationDefault returns a token of Application Default Credentials.
func applicationDefault(ctx context.Context) (string, error) {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", fmt.Errorf("ApplicationDefault: %v", err)
	}
	tk, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("ApplicationDefault: %v", err)
	}
	return tk.AccessToken, nil
}

// gcloud returns a token by running `gcloud auth print-access-token` is a separate process.
func gcloud(ctx context.Context) (string, error) {
	gcloud := "gcloud"
	if runtime.GOOS == "windows" {
		gcloud = "gcloud.cmd"
	}
	cmd := exec.CommandContext(ctx, gcloud, "auth", "print-access-token")
	token, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Gcloud: %v", err)
	}
	return string(token), nil
}

// Token returns oauth2 access token from the environment. It looks for Application Default Credentials
// first and if not found, the credentials of the user logged into gcloud.
func Token(ctx context.Context) (string, error) {
	token, adcErr := applicationDefault(ctx)
	if adcErr != nil {
		var gcloudErr error
		token, gcloudErr = gcloud(ctx)
		if gcloudErr != nil {
			return "", fmt.Errorf("Token: can't find either Application Default Credentials: %v or gcloud credentials %v", adcErr, gcloudErr)
		}
	}
	return token, nil
}

// EncodeJsonKey base64 encodes a service account JSON key file.
func EncodeJsonKey(keyPath string) (string, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("EncodeJsonKey: %v", err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}