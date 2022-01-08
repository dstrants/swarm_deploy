package main

import (
	deploy "swarm_deploy/lib/docker"
)

func main() {
	deploy.UpdateAllServices("nginx:latest")

}
