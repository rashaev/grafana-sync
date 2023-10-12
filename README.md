# grafana-sync

The service synchronizes users and groups from Keycloak to Grafana.


## Build

From the root of the repository run the command below

```sh
docker build -t grafana-sync:<tag> .
```


## Run

Before start the container edit the file with env vars according to your environment

```sh
docker run -d --env-file ./example/env.file grafana-sync:<tag>
```