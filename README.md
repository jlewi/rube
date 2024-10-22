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

The command to fetch the logs for the `rube` app did not return any output. This indicates that there may be no logging entries available for the past hour, or the pod may not have generated logs due to its state.
