```sh
kubectl -n rube port-forward service/rube 8080:80
```

# Look at logs

* Fetch the logs for trace ef2d93e03f47ca8e4f391ba923d4a680
* Use gcloud to fetch the logs from Cloud Logging

```bash {"id":"01JAV9WVEKJ7BH09D4DZ207JS4","interactive":"false"}
gcloud logging read "resource.type=\"k8s_container\" AND jsonPayload.traceId=\"ef2d93e03f47ca8e4f391ba923d4a680\"" --limit=100 --format="table(severity,timestamp,jsonPayload.message)"
```

```bash
gcloud logging read "resource.type=\"k8s_container\" AND jsonPayload.traceId=\"ef2d93e03f47ca8e4f391ba923d4a680\"" --limit=100 --format=yaml
```

##  List the GCP secrets in project foyle-dev

```bash {"id":"01JAVA11G7MTAZ9G6F5QPN3FYX","interactive":"false"}
gcloud secrets list --project=foyle-dev
```

redeploy rube

```bash {"id":"01JAVA22CPD5CB5H6YKNKBHTFY","interactive":"false"}
kubectl apply -f manifests
```

```bash {"interactive":"false"}
kubectl -n rube get pods
```

* To debug the `CrashLoopBackOff` status of the pod `rube-7dc988fc48-k6c7p`, we should check the pod logs to understand what is causing the crash.

```bash {"id":"01JAVA3QJQXS2PZNWD5WSRP4WW","interactive":"false"}
kubectl -n rube logs rube-7dc988fc48-k6c7p
```

* Check whether the service rube-demo in project foyle-dev can read the secret rube-demo-openai-apikey in project foyle-dev
* Use gcloud to troubleshoot the policy

```bash {"id":"01JAVA7B8TMQK6BPFWC68VFSTB","interactive":"false"}
gcloud policy-troubleshoot iam \
    --principal-email="rube-demo@foyle-dev.iam.gserviceaccount.com" \
    //secretmanager.googleapis.com/projects/513170322007/secrets/rube-demo-openai-apikey \
    --permission="secretmanager.versions.access"
```

* How do we give the service account rube-demo in project foyle-dev the permission to read the secret r

To grant the service account `rube-demo@foyle-dev.iam.gserviceaccount.com` in project `foyle-dev` the permission to read the secret, you can use the following command to add the IAM policy binding

* Use terraform to grant the permission
* Write the terraform code to the file iac/main.tf

```bash {"id":"01JAVAABCF98AQ2FYGRRNJ2F8N","interactive":"false"}
# 1. Append the Terraform policy to the iac/main.tf file to grant secret accessor permissions
echo 'resource "google_secret_manager_secret_iam_member" "accessor" {
  secret_id = "rube-demo-openai-api-key"
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:rube-demo@foyle-dev.iam.gserviceaccount.com"
}' >> iac/main.tf
```

```bash {"id":"01JAVAAWWH6RF8W4KXQHHTNYT4","interactive":"true"}
# 2. Apply the Terraform changes to ensure the permissions are granted
cd iac
terraform init
terraform apply -auto-approve
```

```bash {"id":"01JAVAD7N2DTDS23CT9VPXTT8R","interactive":"false"}
# 3. Verify that the IAM policy has been applied correctly
gcloud secrets get-iam-policy rube-demo-openai-apikey
```

The output confirms that the service account `rube-demo@foyle-dev.iam.gserviceaccount.com` has been successfully granted the `roles/secretmanager.secretAccessor` permission for the secret `rube-demo-openai-apikey`. The IAM policy binding is correctly set as shown in the `get-iam-policy` output.

```bash {"id":"01JAVADVWA922DQZQX7J31V8QM","interactive":"false"}
gcloud policy-troubleshoot iam \
    --principal-email="rube-demo@foyle-dev.iam.gserviceaccount.com" \
    //secretmanager.googleapis.com/projects/513170322007/secrets/rube-demo-openai-apikey \
    --permission="secretmanager.versions.access"
```

* Check the Rube pods

```bash {"id":"01JAVAGMWH977V4RQ59YZXE993","interactive":"false"}
kubectl -n rube get pods
```

* Delete the pod rube-7dc988fc48-k6c7p

```bash {"id":"01JAVAH6ZTDN9NK5GKQXS3FNQZ","interactive":"false"}
kubectl -n rube delete pod rube-7dc988fc48-k6c7p
```

```sh
kubectl -n rube get pods
```