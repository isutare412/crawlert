# Crawlert

Crawl any JSON APIs and receive Telegram notifications when specified conditions
are met. Familiarity with [jq](https://jqlang.github.io/jq/) is needed for
writing queries.

## Configuration

Refer to comments in the [sample config file](config.yaml) for setup.

## Run Locally

1. Edit the [config file](config.yaml) to match your requirements.
2. Run `make run`

## Deploy to Kubernetes

1. Create a custom values file (e.g. `my_values.yaml`)
2. Run `helm upgrade --install crawlert ./charts/crawlert -f my_values.yaml`
