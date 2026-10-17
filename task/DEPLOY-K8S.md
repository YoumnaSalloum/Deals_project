# Deploying Deals API to Kubernetes (MicroK8s)

This guide takes the Deals currency-exchange API (Go + Gin + MongoDB + Redis) from
source code to a running app on a **local** Kubernetes cluster using **MicroK8s**.
No prior Kubernetes knowledge needed.

## 0. The 60-second mental model

| Word | What it means here |
| --- | --- |
| **Image** | A packaged, runnable copy of the app (built from the `Dockerfile`). |
| **Container** | A running image. |
| **Pod** | Kubernetes' smallest unit: one (or a few) containers running together. |
| **Deployment** | "Keep N copies of this pod running; replace them if they die." Used for the API and Redis. |
| **StatefulSet** | Like a Deployment, but for databases: stable name + its own disk. Used for MongoDB. |
| **Service** | A stable name + address for a group of pods. The API reaches Mongo at `mongodb:27017` thanks to this. |
| **Secret** | Where passwords/connection strings live. |
| **Namespace** | A folder grouping all of this app's resources (`deals`). |
| **kubectl** | The command-line tool you use to talk to the cluster. |

Flow: **code → Docker image → registry → Kubernetes pulls it → pods run → Service exposes it.**

```
            your machine :30080
                   │  (NodePort Service)
          ┌────────▼────────┐
          │  deals-api x2   │  Deployment
          └───┬─────────┬───┘
     mongodb:27017   redis:6379
     ┌───────▼──┐   ┌───▼─────┐
     │ MongoDB  │   │  Redis  │
     │ + 1Gi PV │   │         │
     └──────────┘   └─────────┘
```

## 1. What was changed in the code to make it cloud-ready

1. `databases/database.go` hard-coded `mongodb://127.0.0.1:27017`. Inside Kubernetes Mongo
   is a different pod, so it now reads `MONGODB_LOCAL_URI` from the environment, and it
   reuses **one** client instead of opening a new connection on every request.
2. `config/default.go` crashed if `app.env` didn't exist. Now `app.env` is optional and real
   environment variables (set by Kubernetes) are used. `PORT` defaults to `8000`.
3. `routes/create.go` used `log.Fatal` on a DB error, which killed the whole server. Now it
   logs the error and the handler returns HTTP 500.
4. `main.go`: the handler no longer writes two JSON responses on success; a small
   variable-shadowing bug in `init()` was fixed.

Local `docker-compose up -d && go run main.go` still works exactly as before.

## 2. New files

```
task/
├── Dockerfile          # multi-stage build → small Alpine image, non-root user
├── .dockerignore
├── deploy.sh           # does steps 5–7 in one command
└── k8s/
    ├── 00-namespace.yaml
    ├── 01-secret.yaml   # Mongo user/password + connection string
    ├── 02-mongodb.yaml  # Service + StatefulSet with 1Gi persistent volume
    ├── 03-redis.yaml    # Service + Deployment
    └── 04-api.yaml      # API Deployment (2 replicas, health probes) + NodePort 30080
```

## 3. Install MicroK8s

**Ubuntu / Linux**

```bash
sudo snap install microk8s --classic
sudo usermod -a -G microk8s $USER
mkdir -p ~/.kube && sudo chown -R $USER ~/.kube
newgrp microk8s                     # or log out and back in
microk8s status --wait-ready
```

**Windows / macOS:** download the installer from <https://microk8s.io> (it runs MicroK8s in a
small Ubuntu VM via Multipass). Then open the `microk8s` shell and follow the Linux steps
from step 4. Ports like `30080` are then on the VM's IP (`multipass list`), not `localhost`.

Optional shortcut so you can type `kubectl` instead of `microk8s kubectl`:

```bash
sudo snap alias microk8s.kubectl kubectl
```

## 4. Turn on the add-ons we need

```bash
microk8s enable dns hostpath-storage registry
```

- **dns** – lets pods find each other by name (`mongodb`, `redis`).
- **hostpath-storage** – gives MongoDB a disk on your machine.
- **registry** – a private image registry inside the cluster at `localhost:32000`.

## 5. Install Docker and build the image

Install Docker if you don't have it (<https://docs.docker.com/engine/install/>). Then:

```bash
git clone https://github.com/YoumnaSalloum/Deals_project.git
cd Deals_project/task
docker build -t localhost:32000/deals-api:1.0 .
docker push localhost:32000/deals-api:1.0
```

> **No Docker?** Build with any tool that makes a tarball, then:
> `docker save localhost:32000/deals-api:1.0 | microk8s ctr image import -`
> and set `imagePullPolicy: IfNotPresent` in `k8s/04-api.yaml`.

## 6. Deploy

```bash
microk8s kubectl apply -f k8s/
```

Watch it come up (Ctrl+C to stop watching):

```bash
microk8s kubectl -n deals get pods -w
```

You want to see something like:

```
NAME                         READY   STATUS    RESTARTS
mongodb-0                    1/1     Running   0
redis-6c9f...                1/1     Running   0
deals-api-7d4b...            1/1     Running   0
deals-api-7d4b...            1/1     Running   0
```

The API pods show `Init:0/1` for a few seconds first — that's the init container waiting
for Mongo and Redis to be up.

**Or do steps 4–6 in one go:** `./deploy.sh`

## 7. Test it

```bash
curl http://localhost:30080/api/healthchecker
# {"message":"Welcome to Golang with Redis and MongoDB","status":"success"}

curl -X POST http://localhost:30080/api/exchange \
  -H 'Content-Type: application/json' \
  -d '{"rates":[{"cfrom":"USD","cto":"GBP","conv":1.5}],"amount":5000,"codeId":1,"created_at":"2022-10-14T17:31:03.735Z"}'
# 201 → "NewAmount": 7500
```

Check the record landed in MongoDB:

```bash
microk8s kubectl -n deals exec -it mongodb-0 -- \
  mongosh -u root -p password123 --quiet --eval 'db.getSiblingDB("myTaskDB").Currs.find()'
```

## 8. Everyday commands

| Task | Command |
| --- | --- |
| See everything | `microk8s kubectl -n deals get all` |
| API logs | `microk8s kubectl -n deals logs -l app=deals-api -f` |
| Why is a pod failing? | `microk8s kubectl -n deals describe pod <name>` |
| Scale to 4 API copies | `microk8s kubectl -n deals scale deploy/deals-api --replicas=4` |
| Ship a new version | `docker build -t localhost:32000/deals-api:1.1 . && docker push localhost:32000/deals-api:1.1 && microk8s kubectl -n deals set image deploy/deals-api api=localhost:32000/deals-api:1.1` |
| Roll back | `microk8s kubectl -n deals rollout undo deploy/deals-api` |
| Port-forward instead of NodePort | `microk8s kubectl -n deals port-forward svc/deals-api 8000:80` |
| Dashboard UI | `microk8s enable dashboard && microk8s dashboard-proxy` |
| Delete everything | `microk8s kubectl delete namespace deals` |

## 9. Troubleshooting

| Symptom | Likely cause / fix |
| --- | --- |
| `ImagePullBackOff` on deals-api | Image wasn't pushed. Re-run step 5 and check `curl localhost:32000/v2/_catalog`. |
| `docker push` says "http: server gave HTTP response to HTTPS client" | Use exactly `localhost:32000` (not your IP), or add it to `insecure-registries` in `/etc/docker/daemon.json`. |
| API stuck in `Init:0/1` | Mongo/Redis not ready. `microk8s kubectl -n deals get pods` and `describe` the db pod. |
| `mongodb-0` stuck in `Pending` | No storage. Run `microk8s enable hostpath-storage`. |
| `CrashLoopBackOff` on deals-api | `microk8s kubectl -n deals logs <pod> --previous` — usually a wrong `MONGODB_LOCAL_URI` in the Secret. |
| `curl localhost:30080` refuses (Windows/macOS) | Use the VM IP from `multipass list` instead of `localhost`. |

## 10. Where to go next

- Put the API behind an **Ingress** (`microk8s enable ingress`) to get `http://deals.local` instead of a port.
- Package the manifests as a **Helm chart** or **Kustomize** overlays (dev/prod).
- Add a **GitHub Actions** workflow that builds the image and pushes to GHCR on every push.
- Move the same manifests to a cloud cluster (GKE / EKS / AKS) — only the image registry and Service type change.
