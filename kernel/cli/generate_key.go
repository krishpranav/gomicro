package cli

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/krishpranav/gomicro/register"
	"github.com/krishpranav/gomicro/tool"
)

type GenerateKey struct {
	register.Command
}

func (c *GenerateKey) Register() {
	c.Signature = "generate:key"
	c.Description = "Generate application key"
}

func (c *GenerateKey) Run() {
	fmt.Println("Generating new application KEY...")
	path := tool.GetDynamicPath("config/server.go")
	read, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	appKey, err := generateNewToken()
	if err != nil {
		log.Fatal(err)
	}

	newContent := strings.Replace(string(read), "REPLACE_WITH_CUSTOM_APP_KEY", appKey, -1)

	if err = os.WriteFile(path, []byte(newContent), 0); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Complete!")
}

func generateNewToken() (string, error) {
	data := make([]byte, 10)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	hashStr := fmt.Sprintf("%x", hash[:])

	return hashStr, nil
}
