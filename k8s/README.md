# Playlister Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the playlister downloader and UI.

## Prerequisites

1. A Kubernetes cluster (single-node is fine with local-path storage)
2. A container registry to host your images
3. A Git repository to store your data files
4. kubectl and (optionally) kustomize

## Architecture

- **Namespace**: `playlister` - isolated environment
- **PVC**: `playlister-data` - ReadWriteOnce storage using local-path
- **CronJob**: Runs daily at 6 AM, pulls git repo, runs downloader, commits and pushes changes
- **UI Deployment**: Single replica (required for ReadWriteOnce storage)
- **Secrets**: API key and git credentials managed via Kustomize or manual YAML

## Storage Model

This setup uses ReadWriteOnce storage with your cluster's default `local-path` StorageClass:
- The PVC can be mounted by pods on one node at a time
- UI is limited to 1 replica to respect ReadWriteOnce constraint
- All pods (UI + CronJob) are scheduled on the same node

**For multiple UI replicas**: Switch to ReadWriteMany storage (NFS, EFS, etc.) and update `k8s/pvc.yaml` to use `accessModes: [ReadWriteMany]`.

## Setup Instructions

### 1. Configure Secrets

**Option A: Kustomize (Recommended)**

Run the setup script to create local secret files:
```bash
./setup-secrets.sh
```

This creates:
- `k8s/secrets/api-key.txt` - Your YouTube API key
- `k8s/secrets/id_rsa` - SSH private key for git access
- Updates git configuration in `k8s/kustomization.yaml`

**Option B: Manual Secret YAML**

1. Copy the example: `cp k8s/secrets.yaml.example k8s/secrets.yaml`
2. Replace placeholder values with your actual secrets
3. Base64 encode binary values:
   ```bash
   echo -n "your_api_key" | base64
   echo -n "$(cat ~/.ssh/id_rsa)" | base64
   ```

### 2. Build and Push Images

Build the downloader image:
```bash
docker build -f Dockerfile.downloader -t your-registry/playlister-downloader:latest .
docker push your-registry/playlister-downloader:latest
```

Update image references in the manifests or use kustomize:
```bash
cd k8s && kustomize edit set image your-registry/playlister-downloader:latest
```

### 3. Deploy to Kubernetes

**With Kustomize (if using Option A):**
```bash
kubectl apply -k k8s/
```

**Manual deployment (if using Option B):**
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/pvc.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/git-init-job.yaml  # One-time setup
kubectl apply -f k8s/downloader-cronjob.yaml
kubectl apply -f k8s/ui-deployment.yaml  # Merge with your existing UI
```

### 4. Verify Deployment

Check the git initialization job:
```bash
kubectl logs -n playlister job/git-init
```

Check CronJob status:
```bash
kubectl get cronjobs -n playlister
kubectl get jobs -n playlister
```

Manually trigger a job for testing:
```bash
kubectl create job --from=cronjob/playlister-downloader -n playlister manual-test-$(date +%s)
kubectl logs -n playlister job/manual-test-$(date +%s) -f
```

## File Locking & Race Conditions

The system uses file locking to prevent race conditions:
- Downloader acquires `.downloader.lock` before running
- CronJob runs with `concurrencyPolicy: Forbid` to prevent overlapping executions
- UI can implement similar locking when modifying files:
  ```bash
  exec 200>/app/data/.ui.lock
  flock -w 10 200  # Wait up to 10s for lock
  # ... perform file operations
  ```

## Git Workflow

The CronJob follows this sequence:
1. **Pull**: `git pull origin main` to sync with remote
2. **Run**: Execute `./downloader` to update db.json and config.json
3. **Commit**: `git add . && git commit -m "Data update: $(date)"`
4. **Push**: `git push origin main` with one retry on conflict

## Scaling & Storage Upgrades

**To enable multiple UI replicas:**
1. Install NFS CSI driver or use cloud ReadWriteMany storage
2. Update `k8s/pvc.yaml` to use ReadWriteMany access mode
3. Increase UI replicas in `k8s/ui-deployment.yaml`
4. Remove node affinity constraints

**NFS Example:**
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/csi-driver-nfs/master/deploy/install-driver.yaml
# Then update pvc.yaml to use nfs-csi storageClassName
```

## Troubleshooting

**PVC won't mount**: Check if your cluster supports the storage class:
```bash
kubectl get storageclass
kubectl describe pvc playlister-data -n playlister
```

**Git authentication fails**: Verify SSH key format and GitHub access:
```bash
kubectl exec -n playlister deployment/playlister-ui -- ssh -T git@github.com
```

**CronJob not running**: Check schedule and job history:
```bash
kubectl describe cronjob playlister-downloader -n playlister
kubectl get jobs -n playlister
```