package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/joho/godotenv"
)

var (
	apiKey       string
	instanceType string
	sshKeyID     int64
	errorWait    int
	client       *hcloud.Client
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	apiKey = os.Getenv("HETZNER_API_KEY")
	if apiKey == "" {
		log.Fatalln("HETZNER_API_KEY is required")
	}

	instanceType = os.Getenv("HETZNER_INSTANCE_TYPE")
	if instanceType == "" {
		log.Fatalln("HETZNER_INSTANCE_TYPE is required")
	}

	sshKeyIDStr := os.Getenv("SSH_KEY_ID")
	if sshKeyIDStr == "" {
		log.Fatalln("SSH_KEY_ID is required")
	}
	keyId, err := strconv.ParseInt(sshKeyIDStr, 10, 64)
	if err != nil {
		log.Fatalln("SSH_KEY_ID couldn't be converted to int64")
	}
	sshKeyID = keyId

	errorWait, _ = strconv.Atoi(os.Getenv("ERROR_WAIT_SECONDS"))
	if errorWait == 0 {
		errorWait = 60
	}

	client = hcloud.NewClient(hcloud.WithToken(apiKey))
}

func launchInstanceLoop() {
	for {
		log.Println("========================================")

		// Get all images
		images, err := client.Image.All(context.Background())
		if err != nil {
			log.Println("Error getting images:", err)
			time.Sleep(time.Duration(errorWait) * time.Second)
			continue
		}

		// Find the latest Ubuntu image
		ubuntuImages := []*hcloud.Image{}
		for _, image := range images {
			if strings.Contains(image.Name, "ubuntu") {
				ubuntuImages = append(ubuntuImages, image)
			}
		}
		latestUbuntuImage := ubuntuImages[len(ubuntuImages)-1]

		log.Println("Latest Ubuntu image:", latestUbuntuImage.Name)

		serverCreateOpts := hcloud.ServerCreateOpts{
			Name:       "hetzner-bot-" + time.Now().Format("2006-01-02-15-04-05"),
			Image:      &hcloud.Image{Name: latestUbuntuImage.Name},
			ServerType: &hcloud.ServerType{Name: instanceType},
			SSHKeys:    []*hcloud.SSHKey{{ID: sshKeyID}},
		}
		serverCreateResult, _, err := client.Server.Create(context.Background(), serverCreateOpts)
		if err != nil {
			log.Println("Error creating server:", err)
			time.Sleep(time.Duration(errorWait) * time.Second)
			continue
		}
		log.Println("Server created:", serverCreateResult.Server.Name)
		break
	}
}

func main() {
	log.Println("Starting hetzner-bot...")
	launchInstanceLoop()
}
