#!/usr/bin/env bash
# One-shot deploy of the Deals API to MicroK8s. Run from the task/ folder.
set -euo pipefail
IMAGE=localhost:32000/deals-api:1.0

echo "==> Enabling MicroK8s add-ons (dns, storage, registry)"
microk8s enable dns hostpath-storage registry

echo "==> Building image $IMAGE"
docker build -t "$IMAGE" .
echo "==> Pushing image to the MicroK8s registry"
docker push "$IMAGE"

echo "==> Applying Kubernetes manifests"
microk8s kubectl apply -f k8s/

echo "==> Waiting for everything to be ready"
microk8s kubectl -n deals rollout status statefulset/mongodb --timeout=180s
microk8s kubectl -n deals rollout status deployment/redis --timeout=120s
microk8s kubectl -n deals rollout status deployment/deals-api --timeout=180s

microk8s kubectl -n deals get pods,svc
echo
echo "==> Test it:"
echo "curl http://localhost:30080/api/healthchecker"
