# rube

Hello world for Foyle Demo

## Build the image

```bash {"id":"01JATZXXEKVEKSE44Z3PY1S45F","interactive":"true"}
KO_DOCKER_REPO=us-west1-docker.pkg.dev/foyle-public/images/rube ko build .
```

## Deploy it

```bash {"id":"01JAV0AB0THJFS20GKNWMXMN05","interactive":"true"}
kustomize build ./manifests/dev | kubectl apply -f -
kustomize build ./manifests/prod | kubectl apply -f -
```

```bash {"id":"01JAVQ4MPFX0G0E1A9RA1PEZWM","interactive":"true"}
kustomize build ./manifests/prod | kubectl apply -f -
```

```bash {"id":"01JAVQHJNJJE5XTY2B7HJJTC8P","interactive":"false"}
kubectl -n rube get pods -w
```

```bash {"id":"01JAV0B6T625Z5G4KM3WA4SPCN","interactive":"false"}
kubectl -n rube get pods -w
```

To debug the `CrashLoopBackOff` status of the pod, you should inspect the logs of the pod to identify any errors or issues that might be causing the crash. You can do this using the following command:

## Monitoring

### Logs

Fetch the logs for the rube app using gcloud

```bash {"id":"01JAV0FH9HF58VMYEWAJ7CKD5Z","interactive":"true"}
gcloud logging read "resource.type=\"k8s_container\" AND labels.k8s-pod/app=\"rube\" labels.k8s-pod/env=\"prod\"" --limit=100 --freshness=1h --format="table(severity,timestamp,jsonPayload.traceId,textPayload,jsonPayload.message)"
```

To proceed with debugging the `CrashLoopBackOff` status of the `rube-8689fcbdb-nsglc` pod, the next step is to fetch the logs specifically for that pod to examine any error messages or issues that may be causing the crashes. You can use the following command:

```bash
kubectl -n rube port-forward svc/rube 8080:80
```

### Monitor Latency

* Show the 95th percentile of latency

```bash
QUERY='{
  "time_range": 86400,
  "granularity": 0,
  "breakdowns": [],
  "calculations": [
    {
      "op": "P95",
      "column": "duration_ms"
    }
  ],
  "filters": [
    {
      "column": "name",
      "op": "=",
      "value": "/"
    }
  ],
  "filter_combination": "AND",
  "orders": [],
  "havings": [],
  "trace_joins": [],
  "limit": 1000
}'

hccli querytourl --dataset rube --query "$QUERY" --base-url https://ui.honeycomb.io/autobuilder/environments/prod --open
```

```bash {"id":"01JAWSAP5JEFGXDX3YN22Q89JW","interactive":"false"}
# Fetch detailed logs for the rube pod to check for errors
kubectl -n rube logs rube-8689fcbdb-nsglc

# Alternatively, fetch the logs for the rube app using gcloud
gcloud logging read "resource.type=\"k8s_container\" AND labels.k8s-pod/app=\"rube\" AND labels.k8s-pod/env=\"prod\"" --limit=100 --freshness=1h --format="table(severity,timestamp,jsonPayload.traceId,textPayload,jsonPayload.message)"

# If further inspection is needed, forward the port to access the application
kubectl -n rube port-forward svc/rube 8080:80
```

## Setup the gateway

* dev.foyle.io should be pointing at the IP address

```sh
# create the namespace
kubectl create namespace gateway
```

```bash {"id":"01JAVBT0TSW5X05QV13WTJDN6X","interactive":"true"}
kubectl apply -f manifests/gateway.yaml
kubectl apply -f manifests/httproute.yaml
```

```bash
kubectl -n gateway  describe gateway
```

```bash {"id":"01JAVCQ2EP3X6X8Z53R5HFJZSH","interactive":"false"}
# Now that the namespace and required resources are created, let's check the deployment status of the rube app.
kubectl -n rube get deployments
kubectl -n rube describe deployment rube
```

* List the static ip addresses

```bash {"id":"01JAVCDY7XDY2X58WPV8T6111R","interactive":"false"}
gcloud compute addresses list --global
```
