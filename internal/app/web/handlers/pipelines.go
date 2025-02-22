package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func SetupPipelines(app *fiber.App) {
	pipelinesGroup := app.Group("/pipelines")
	pipelinesGroup.Post("/github-webhook", HandlePushEvent)

}

type RequestBody struct {
	Url    string `json:"url" xml:"url" form:"url"`
	Branch string `json:"branch" xml:"branch" form:"branch"`
}
