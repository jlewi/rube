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

## Monitoring

Fetch the logs for the rube app using gcloud

```bash {"id":"01JAV0FH9HF58VMYEWAJ7CKD5Z","interactive":"true"}
gcloud logging read "resource.type=\"k8s_container\" AND labels.k8s-pod/app=\"rube\"" --limit=100 --freshness=1h --format="table(severity,timestamp,jsonPayload.message)"
```

```bash {"id":"01JAV0S2H6VJVWS2HS0JEEHEH0","interactive":"false"}
kubectl -n rube get secrets
```

To debug the `CrashLoopBackOff` status of the pod in the `rube` namespace, we will fetch the logs for the specific pod that's encountering issues.

```bash
kubectl -n rube logs  rube-8654b7854d-8p8zh
```
