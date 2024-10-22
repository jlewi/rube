# Configure the Google Cloud provider
provider "google" {
  project = "foyle-dev"
  region  = "us-west1"
}

# Configure the backend to use a GCS bucket
terraform {
  backend "gcs" {
    bucket = "foyle-dev-tfstate"
    prefix = "rube"
  }
}

resource "google_service_account" "rube_service_account" {
  account_id   = "rube-demo"
  display_name = "Rube Service Account"
}



# Grant the service account permission to read the secret for the honeycommb API key
resource "google_secret_manager_secret_iam_member" "secret_reader" {
  project   = "foyle-dev"
  secret_id = "honeycomb-api-key"
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:rube-demo@foyle-dev.iam.gserviceaccount.com"
}

# Enable workload identity
resource "google_service_account_iam_binding" "allow_impersonation" {  
  service_account_id = "projects/foyle-dev/serviceAccounts/rube-demo@foyle-dev.iam.gserviceaccount.com"
  role = "roles/iam.workloadIdentityUser"
  members = [
    "serviceAccount:foyle-dev.svc.id.goog[rube/rube]"
  ]
}
resource "google_secret_manager_secret_iam_member" "accessor" {
  secret_id = "openai-api-key"
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:rube-demo@foyle-dev.iam.gserviceaccount.com"
}
