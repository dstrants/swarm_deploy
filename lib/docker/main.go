package containers

import (
	"errors"
	"fmt"
	"strings"

	"swarm_deploy/lib/slack"

	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

//Parse image name and tag
func ParseImageName(full_name string) (string, string, error) {
	result := strings.Split(full_name, ":")

	if len(result) != 2 {
		err := errors.New("Parsing failed for image: " + full_name)
		return "", "", err
	}

	return result[0], result[1], nil
}

// Selects the swarm services that need to be updated based on image
// and swarm label
func LocateServices(image string) []swarm.Service {
	client, err := docker.NewClientFromEnv()

	if err != nil {
		panic(err)
	}
	allServices, err := client.ListServices(docker.ListServicesOptions{})

	if err != nil {
		log.Printf("There was an error while trying to list the services. %v", err)
	}

	var toBeUpdated []swarm.Service

	for _, service := range allServices {
		currentImage, _, imageError := ParseImageName(service.Spec.Labels["com.docker.stack.image"])

		if imageError != nil {
			log.WithFields(
				log.Fields{
					"incoming_image": image,
					"current_image":  service.Spec.Labels["com.docker.stack.image"],
					"service":        service.PreviousSpec.Name,
				}).Error("Could not parse current image to compare with incoming. Service will be skipped.")
			continue
		}

		if currentImage == image && service.Spec.Labels["swarm_deploy"] == "true" {
			toBeUpdated = append(toBeUpdated, service)
		}
	}
	return toBeUpdated
}

// Updates a service to given image
func UpdateServiceImage(service swarm.Service, image string) bool {
	client, err := docker.NewClientFromEnv()

	if err != nil {
		panic(err)
	}

	service.Spec.TaskTemplate.ContainerSpec.Image = image

	var updateConfig docker.UpdateServiceOptions
	updateConfig.ServiceSpec = service.Spec
	updateConfig.Version = service.Version.Index

	err = client.UpdateService(service.ID, updateConfig)

	if err != nil {
		log.WithFields(log.Fields{"service": service, "image": image}).Error("Service update failed,")
		return false
	}
	return true
}

// Parent function that triggers locating and updating services
func UpdateAllServices(image, tag string) {
	services := LocateServices(image)

	if len(services) == 0 {
		log.WithFields(log.Fields{"image": image}).Warning("No services found for the given image.")
		return
	}

	full_image := fmt.Sprintf("%s:%s", image, tag)

	for _, service := range services {
		result := UpdateServiceImage(service, full_image)

		if result {
			slack.SendSimpleMessage(fmt.Sprintf("Service `%s` has been deployed with image `%s`", service.Spec.Name, full_image))
			log.WithFields(log.Fields{"service": service.Spec.Name, "image": full_image}).Info("Service image updated")
		}
	}
}
