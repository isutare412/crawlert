# Crawlert

Crawl any JSON APIs and receive telegram messages when specific condition are
met.  Knowledge of [jq](https://jqlang.github.io/jq/) is required to write
queries.

## Configuration

Reference comments in [sample config file](config.yaml).

## Run Locally

1. Modify [config](config.yaml) to suit your needs.
2. `make run`

## Deploy to Kubernetes

1. Prepare you own values file (e.g. `my_values.yaml`)
2. `helm upgrade --install crawlert ./charts/crawlert -f my_values.yaml`
