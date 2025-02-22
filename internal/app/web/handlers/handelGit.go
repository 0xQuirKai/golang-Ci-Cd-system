package handlers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Struct for GitHub Push Event Payload
type PushEvent struct {
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commits"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Ref string `json:"ref"`
}

// ✅ Function to handle push events and print commit IDs
func HandlePushEvent(c *fiber.Ctx) error {
	var event PushEvent

	// 📖 Parse the JSON payload
	if err := json.Unmarshal(c.Body(), &event); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("❌ Invalid payload")
	}

	// 📢 Print repository and branch
	fmt.Printf("🚀 Push to repository: %s on branch: %s\n", event.Repository.FullName, event.Ref)

	// 📋 Loop through commits and print their IDs

	rb := RequestBody{
		Url:    os.Getenv("GITHUB_REPO_URL"),
		Branch: os.Getenv("GITHUB_BRANCH"),
	}
	primaryLang, err := detectFrameworkFromGitHub(event.Repository.FullName)

	err = PostCheckItWorks(rb, primaryLang)
	fmt.Println("Primary Language: ", primaryLang)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("❌ Pipeline failed: %v", err))
	}
	return c.SendString("✅ Push event processed")
}
