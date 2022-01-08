package deploy

import (
	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

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
		if service.Spec.Labels["com.docker.stack.image"] == image && service.Spec.Labels["swarm_deploy"] == "true" {
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
func UpdateAllServices(image string) {
	services := LocateServices(image)

	for _, service := range services {
		result := UpdateServiceImage(service, image)

		if result {
			log.WithFields(log.Fields{"service": service.Spec.Name, "image": image}).Info("Service image updated")
		}
	}
}
