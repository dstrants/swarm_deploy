# Swarm Deploy

A Simple CD tool for [docker swarm](https://docs.docker.com/engine/swarm/). The view of the tool is to accept webhooks from container registries and deploy that image to the services that have the related label and use the given image.
Current implementation is very simple and just deploys the image as is.


## Deploy in swarm as a service

```yaml
# swarm_deploy.yml
version: '3.8'

services:
    server:
        environment:
            # Use this for production mode
            GIN_MODE: release
            # You have to set this on the github side as well
            GITHUB_WEBHOOK_SECRET: supersecretpass
            SLACK_TOKEN: extrasupersecret

        deploy:
            replicas: 1
        image: ghcr.io/dstrants/swarm_deploy:latest
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock
```

Then run the command deploy to deploy the service to you cluster
```sh
docker stack deploy -c swarm_deploy.yml swarm_deploy
```

### Full Configuration

Name | Default | Description
-----|---------|------------
`WEBSERVER_HOST` | "" | The hostname that webserver will be listening to.
`WEBSERVER_PORT` | 8080 | The port that webserver will be listening to.
`GITHUB_WEBHOOK_SECRET` | **required** | The secret that github webhooks will use to produce the HMAC sha. Should be the same in all webhooks or call will be rejected.
`GITHUB_WEBHOOK_EVENTS` | `package,ping` | The webhooks events that the app is accepting (should be comma separated values)
`SLACK_TOKEN` | **required** | The slack app token to deliver notifications.
`SLACK_CHANNEL` | `infra` | The slack channel that notifications will be delivered to.


## TODO
- [x] ~Add code related ci tools (linting, etc)~
- [ ] Add regex filters as deploy labels to filter what out images tags
- [ ] Adds webhook endpoint for dockerhub
- [ ] Add protected endpoint for manual triggering.
- [ ] Slack notifications
- [ ] Improve documentation
