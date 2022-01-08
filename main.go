package main

import (
	"encoding/json"
	"net/http"
	github "swarm_deploy/lib/github"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	r := setupRouter()
	r.Run(":8080")
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.POST("webhook/github", func(c *gin.Context) {
		var package_update github.GithubPackageWebhook
		log.Info("New incoming webhook from github")
		if c.ShouldBindBodyWith(&package_update, binding.JSON) == nil {
			data, _ := json.Marshal(package_update)
			result := github.CalculateHMAC(data, "supersecret")
			header := c.Request.Header["X-Hub-Signature-256"][0]
			if "sha256="+result != header {
				log.WithFields(log.Fields{"incoming": header, "computed": "sha256=" + result}).Error("Authentication failed.")
				c.JSON(http.StatusForbidden, gin.H{"status": "You are not authenticated"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"status": "Thanks for the heads up :)"})

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request"})
		}
		return
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}
