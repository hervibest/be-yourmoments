package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

var (
	vaultConfig struct {
		Host   string
		Port   string
		Auth   string
		Token  string
		Engine string
		Path   string
	}
	vaultClient *vault.Client
	vaultOnce   sync.Once
)

func init() {
	vaultConfig.Host = os.Getenv("VAULT_HOST")
	vaultConfig.Port = os.Getenv("VAULT_PORT")
	vaultConfig.Auth = os.Getenv("VAULT_AUTH")
	vaultConfig.Token = os.Getenv("VAULT_TOKEN")
	vaultConfig.Engine = os.Getenv("VAULT_ENGINE")
	vaultConfig.Path = os.Getenv("VAULT_PATH")
}

func logFailure(key string, sources []string) {
	logger := NewLogger()
	logger.Errorw("Failed to get key-value pair", "failed sources", sources, "key", key)
}

func GetEnv(key string) string {
	var failedSources []string

	if value := getOSEnv(key, &failedSources); value != "" {
		return value
	}

	if value := getDotEnv(key, &failedSources); value != "" {
		return value
	}

	// if value := getVaultEnv(key, &failedSources); value != "" {
	// 	return value
	// }

	logFailure(key, failedSources)
	return ""
}

func getOSEnv(key string, failedSources *[]string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	*failedSources = append(*failedSources, "os")
	return ""
}

func getDotEnv(key string, failedSources *[]string) string {
	if err := godotenv.Load("../../.env"); err != nil {
		*failedSources = append(*failedSources, ".env")
		return ""
	}
	return getOSEnv(key, failedSources)
}

func getVaultEnv(key string, failedSources *[]string) string {
	client, err := getVaultClient()
	if err != nil {
		*failedSources = append(*failedSources, "vault")
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	secret, err := client.KVv2(vaultConfig.Engine).Get(ctx, vaultConfig.Path)
	if err != nil {
		*failedSources = append(*failedSources, "vault")
		return ""
	}

	if value, ok := secret.Data[key].(string); ok {
		return value
	}
	*failedSources = append(*failedSources, "vault")
	return ""
}

func getVaultClient() (*vault.Client, error) {
	var err error
	vaultOnce.Do(func() {
		if !isVaultConfigValid() {
			err = fmt.Errorf("invalid vault configuration")
			return
		}

		vaultURL := fmt.Sprintf("http://%s:%s", vaultConfig.Host, vaultConfig.Port)
		if !isVaultReachable(vaultURL) {
			err = fmt.Errorf("vault is not reachable")
			return
		}

		config := vault.DefaultConfig()
		config.Address = vaultURL

		vaultClient, err = vault.NewClient(config)
		if err != nil {
			return
		}
		vaultClient.SetToken(vaultConfig.Token)
	})

	return vaultClient, err
}

func isVaultConfigValid() bool {
	return vaultConfig.Host != "" && vaultConfig.Port != "" && vaultConfig.Auth != "" &&
		vaultConfig.Token != "" && vaultConfig.Engine != "" && vaultConfig.Path != ""
}

func isVaultReachable(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
