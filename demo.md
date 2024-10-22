# Demo

## Setup

Port forward to the service to make it accessible from the local machine:

```bash {"id":"01JAV2SP5BFF7W5N4AH7K1PMWF","interactive":"false"}
kubectl -n rube port-forward svc/rube 8080:80
```

## Setup

* Remove the permission to read the secret openai-api-key by the rube-demo account
* Make sure it is configured to use openai-api-key and deployed.

## Observability

* Open the app in the browser and get a traceId

```bash {"id":"01JAV2V66C9HVH0CE381VS4MTH","interactive":"false"}
open http://localhost:8080
```

Fetch the rube logs for trace cdd7cd1432902c8ee88c6a9861cb4926

```bash {"id":"01JAV2QF6Y01XA9JQDHM1DTPKY","interactive":"false"}
gcloud logging read "resource.type=\"k8s_container\" AND labels.k8s-pod/app=\"rube\" AND jsonPayload.traceId=\"cdd7cd1432902c8ee88c6a9861cb4926\"" --limit=100 --format="table(severity,timestamp,jsonPayload.message,jsonPayload.traceId)"
```

```bash {"id":"01JAV2PN7NB7D9D7Y2KZBPV575","interactive":"false"}
gcloud logging read "resource.type=\"k8s_container\" AND labels.k8s-pod/app=\"rube\"" --limit=100 --freshness=1h --format="table(severity,timestamp,jsonPayload.message,jsonPayload.traceId)"
```

## Debug rube app

Why isn't the rube app running correctly?

To debug why the rube app isn't running correctly, we can take the following steps:

```bash {"id":"01JAV4KQT031ZMZ6TTAHHA6ZRC","interactive":"false"}
gcloud logging read "resource.type=\"k8s_container\" AND labels.k8s-pod/app=\"rube\" severity=\"ERROR\"" --limit=1 --freshness=1h --format="table(severity,timestamp,jsonPayload.error,jsonPayload.traceId)"
```

To resolve the error indicated in the logs, it appears that the service account being used does not have the necessary permissions to access the secret containing the API key.

* Use gcloud to get the IAM policy for the secret

```bash {"id":"01JAV5W898X74H5F59FF036JC2","interactive":"false"}
gcloud secrets get-iam-policy openai-api-key
```

* Use gcloud iam policy troubleshooter to check whether the service account rube-demo in project foyle-dev can
  read the secret 

```bash {"id":"01JAV62Z5A1JAXK63MBSY4ZHN9","interactive":"true"}
gcloud policy-troubleshoot iam \
    --principal-email="rube-demo@foyle-dev.iam.gserviceaccount.com" \
    //secretmanager.googleapis.com/projects/513170322007/secrets/openai-api-key \
    --permission="secretmanager.versions.access"
```

* Use gcloud to decribe the secretAccessor role

```bash {"id":"01JAV6H5ZA5KG9ZYFN71J2N1G1","interactive":"false"}
# 1. Describe the secretAccessor role to understand its permissions
gcloud iam roles describe roles/secretmanager.secretAccessor
```

* Describe the secret openai-api-key in project foyle-dev

```bash {"id":"01JAV6KJ5AC4GG5ZX51P0MYHJ1","interactive":"false"}
gcloud secrets describe openai-api-key --project="foyle-dev"
```

## Fix the permission

* Write a terraform policy to grant the GCP service account rube-dev secret accessor permissions on the secret openai-api-key
* Append the terraform policy to the file iac/main.tf

```bash {"id":"01JAV6XANFZVRD4FMDW2G8JPNN","interactive":"false"}
# 1. Append the Terraform policy to the iac file to grant secret accessor permissions
echo 'resource "google_secret_manager_secret_iam_member" "accessor" {
  secret_id = "openai-api-key"
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:rube-demo@foyle-dev.iam.gserviceaccount.com"
}' >> iac/main.tf
```

```bash {"id":"01JAV6XY73KNM1TJ0PFH2GMKFW","interactive":"false"}
# 2. Apply the Terraform configuration to ensure the permissions are updated
cd iac
terraform init
terraform apply -auto-approve
```

The output indicates that the `google_secret_manager_secret_iam_member.accessor` has been successfully created, providing the service account `rube-demo@foyle-dev.iam.gserviceaccount.com` with the role `roles/secretmanager.secretAccessor` for the secret `openai-api-key`.

* Lets verify the permissions using iam policy troubleshooter
* Use gcloud to check whether the rube-demo email can access the secret

```bash {"id":"01JAV6ZGCS9YKXEH6A2M7J6A30","interactive":"false"}
gcloud policy-troubleshoot iam \
    --principal-email="rube-demo@foyle-dev.iam.gserviceaccount.com" \
    //secretmanager.googleapis.com/projects/513170322007/secrets/openai-api-key \
    --permission="secretmanager.versions.access"
```

* Now check the rube pods to see if they are healthy

```bash {"id":"01JAV70VMD3NJNWQV39E47FP00","interactive":"false"}
kubectl -n rube get pods
```

To debug the `CrashLoopBackOff` status of the `rube-8689fcbdb-nsglc` pod, follow these steps:

* Fetch the logs for the rube app using gcloud 

```bash {"id":"01JAV724F7SEVY36SRG0Z53ZMN","interactive":"false"}
gcloud logging read "resource.type=\"k8s_container\" AND labels.k8s-pod/app=\"rube\"" --limit=100 --freshness=1h --format="table(severity,timestamp,jsonPayload.message,jsonPayload.error,jsonPayload.traceId)"
```