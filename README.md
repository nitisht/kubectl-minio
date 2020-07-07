# MinIO Kubectl Plugin

This is a `kubectl` plugin to interact with MinIO Operator on Kubernetes.

## Prerequisites

- Kubernetes cluster with storage nodes earmarked.
- Required number of drives attached to each storage node. All drives mounted and formatted.
- kubectl installed on your local machine, configured to talk to the Kubernetes cluster.

Common Flags

- `--namespace=minio-operator`
- `--kubeconfig=`

## Operator Deployment

Command:

`kubectl minio init [options] [flags]`

Options:

- `--image=minio/k8s-operator:2.0.6`
- `--namespace-to-watch=default`
- `--cluster-domain=cluster.local`
- `--service-account=minio-operator`

Flags:

- `--o`
- `--dry-run`

## Tenant

### MinIO Tenant Creation

Command:

`kubectl minio tenant create NAME <access_key> <secret_key> --zones="rack1:4,rack2:8" --volumesPerNode=4 --capacityPerVolume=1Ti [options] [flags]`

Options:

- `--namespace=minio`
- `--storageClass=local`
- `--kms-secret=secret-name`
- `--console-secret=secret-name`
- `--external-cert-secret=secret-name`

Flags:

- `--auto-tls`
- `-o`

### Scale Tenant Zones

Command:

`kubectl minio tenant scale NAME --zones="rack3,4","rack4,8"`

### Update Images

Command:

`kubectl minio tenant update NAME --image=minio/minio:RELEASE.2020-06-12T00-06-19Z`

### Remove Tenant

Command:

`kubectl minio tenant delete NAME`
