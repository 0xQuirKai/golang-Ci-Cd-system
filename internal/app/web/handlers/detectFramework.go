package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type RequestBodyNew struct {
	Repo        string `json:"repo"`        // e.g., "username/repository"
	Branch      string `json:"branch"`      // e.g., "main"
	AccessToken string `json:"accessToken"` // GitHub App installation token or PAT
}

func postDetectFramework(c *fiber.Ctx) error {
	body := &RequestBodyNew{}

	if err := c.BodyParser(body); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}

	// ✅ Detect framework directly using GitHub API
	framework, err := detectFrameworkFromGitHub(body.Repo)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error detecting framework: %v", err))
	}

	// ✅ Return the detected framework
	return c.JSON(fiber.Map{
		"repository": body.Repo,
		"branch":     body.Branch,
		"framework":  framework,
	})
}

func detectFrameworkFromGitHub(repo string) (string, error) {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	// 📡 GitHub /languages API
	url := fmt.Sprintf("https://api.github.com/repos/%s/languages", repo)
	fmt.Println("Fetching languages from:", url)

	// 🌐 Create HTTP Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "unknown", err
	}

	// 🛡️ Set Headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// 💥 Send Request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "unknown", err
	}
	defer resp.Body.Close()

	// ⚠️ Handle Non-200 Responses
	if resp.StatusCode != http.StatusOK {
		return "unknown", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	// 📖 Read and Parse Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "unknown", err
	}

	// 🖨️ Print Raw JSON Response
	fmt.Println("GitHub API /languages Response:", string(body))

	// 🔄 Parse JSON into Map
	var languages map[string]int
	if err := json.Unmarshal(body, &languages); err != nil {
		return "unknown", err
	}
	// 🎯 Identify Primary Language
	var primaryLanguage string
	var maxBytes int
	for lang, bytes := range languages {
		if bytes > maxBytes {
			primaryLanguage = lang
			maxBytes = bytes
		}
	}

	if primaryLanguage != "" {
		return primaryLanguage, nil
	}

	return "unknown", nil
}
