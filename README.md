# rube

Hello world for Foyle Demo

## Build the image

```bash {"id":"01JATZXXEKVEKSE44Z3PY1S45F","interactive":"true"}
KO_DOCKER_REPO=us-west1-docker.pkg.dev/foyle-public/images/rube ko build .
```

## Deploy it

```bash {"id":"01JAV0AB0THJFS20GKNWMXMN05","interactive":"true"}
kubectl apply -f ./manifests
```

```bash {"id":"01JAV0B6T625Z5G4KM3WA4SPCN","interactive":"false"}
kubectl -n rube get pods -w
```

To debug the `CrashLoopBackOff` status of the pod, you should inspect the logs of the pod to identify any errors or issues that might be causing the crash. You can do this using the following command:

## Monitoring

Fetch the logs for the rube app using gcloud

```bash {"id":"01JAV0FH9HF58VMYEWAJ7CKD5Z","interactive":"true"}
gcloud logging read "resource.type=\"k8s_container\" AND labels.k8s-pod/app=\"rube\"" --limit=100 --freshness=1h --format="table(severity,timestamp,jsonPayload.message,jsonPayload.traceId)"
```

```bash
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
