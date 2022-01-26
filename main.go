package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"swarm_deploy/lib/config"
	containers "swarm_deploy/lib/docker"
	githubModels "swarm_deploy/lib/github"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/go-github/v42/github"
)

const Version = "0.2.3"

var cnf = config.LoadConfig()

func WebhookTypeChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header
		eventType, ok := headers["X-Github-Event"]

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "UnprocessableEntity"})
		}

		for _, event := range cnf.GithubWebhookEvents() {
			if event == eventType[0] {
				return
			}
		}

		log.WithFields(log.Fields{"field": eventType[0]}).Error("Event not in the list of accepted events")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "UnprocessableEntity"})
	}
}

func WebhookSignatureValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Copy()

		payload, err := github.ValidatePayload(req.Request, []byte(cnf.Github.WebhookSecret))
		if err != nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(payload))
		c.Next()
	}
}

func main() {
	r := setupRouter()
	r.Run(fmt.Sprintf("%s:%d", cnf.WebServer.Host, cnf.WebServer.Port))
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.Use(WebhookTypeChecker())
	r.Use(WebhookSignatureValidator())

	r.POST("webhook/github", func(c *gin.Context) {
		log.Info("New incoming webhook from github")
		eventType := c.Request.Header["X-Github-Event"][0]

		if eventType == "ping" {
			var ping github.PingEvent
			if e := c.ShouldBindBodyWith(&ping, binding.JSON); e == nil {
				log.Debug("Received an ping event")
				c.JSON(http.StatusOK, gin.H{"status": "Hello github!"})
			} else {
				log.Error(e)
				c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request"})
			}
			return
		}

		var package_update githubModels.PackageEvent
		if e := c.ShouldBindBodyWith(&package_update, binding.JSON); e == nil {
			escapedImage := strings.Replace(package_update.Package.PackageVersion.PackageURL, "\n", "", -1)
			escapedImage = strings.Replace(escapedImage, "\r", "", -1)
			image, tag, err := containers.ParseImageName(escapedImage)

			if err != nil {
				log.WithFields(log.Fields{"image": escapedImage}).Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Image could not be parsed", "image": escapedImage})
			}

			// Spawns an async process to update the services.
			// Responses just acknowledges the webhook so it won't be retried
			go containers.UpdateAllServices(image, tag)

			c.JSON(http.StatusCreated, gin.H{"status": "Thanks for the heads up :)"})

		} else {
			log.Error(e)
			c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request"})
		}
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}
