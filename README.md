# kubernetes-sandbox

This repository is a sandbox for development, testing and building of a demo
Kubernetes application.

- [Prerequisites](#prerequisites)
- [Kubernetes Clusters](#kubernetes-clusters)
  - [Local](#local)
  - [Remote](#remote)
- [HTTP Echo Server](#http-echo-server)
- [Local Development Workflow](#local-development-workflow)
- [Install ArgoCD](#install-argocd)
- [Deploy Echo Server](#deploy-echo-server)
- [Github Actions](#github-actions)
- [Vulnerability Scanning](#vulnerability-scanning)
- [Releasing Changes](#releasing-changes)
- [Further Work](#further-work)

## Prerequisites

There are prequisites for using this sandbox:
- kind
- kustomize
- skaffold
- kubeconform
- Go

All of these are present in the [VSCode Remote Containers][remote-containers]
configuration in the `.devcontainer/` folder.

## Kubernetes Clusters

### Local

The sandbox workflow here uses a Kind cluster to test local changes. This can
be configured locally with the Makefile command

```sh
make start-kube
```

### Remote

The ArgoCD components should be installed into a remote Kubernetes cluster, and
boostrapped with the Application custom resources in `kubernetes/argocd`.

For this sandbox we can start a multi-node cluster with Kind using the
configuration at `kubernetes/kind`.

```sh
kind create cluster \
  --name production-cluster \
  --image kindest/node:v1.22.9 \
  --config kubernetes/kind/multi_node.yaml
```

Once the cluster is up we should have a control-plane and 3 worker nodes.

```sh
NAME                               STATUS   ROLES                  AGE    VERSION
production-5ksym8t-control-plane   Ready    control-plane,master   149m   v1.22.9
production-5ksym8t-worker          Ready    <none>                 148m   v1.22.9
production-5ksym8t-worker2         Ready    <none>                 148m   v1.22.9
production-5ksym8t-worker3         Ready    <none>                 148m   v1.22.9
```

## HTTP Echo Server

There is a dummy HTTP echo server written in Go in `main.go`. The application
can be tested by running

```sh
make test
```

## Local Development Workflow

> ⚠️ Skaffold will use your active Kubernetes context. If this is set as a
> production cluster you may cause issues or conflicts!

Use `skaffold` to run a local, fast feedback loop for Kubernetes development.

```sh
skaffold dev
```

This will setup a monitoring loop on your application files, Dockerfile and
Kubernetes manifests. When changes are detected, the application will be
rebuilt, the container image will be published into the Kind cluster node, and
the Kubernetes manifests re-applied to the cluster.

You can also trigger individual parts of the skaffold development loop.

Run the pipeline through from end-to-end.

```sh
skaffold run
```

Build the artifacts and exit

```sh
skaffold build
```

## Install ArgoCD

> ⚠️ Make sure your Kubernetes context is pointing at the Production/remote
> cluster to avoid conflicts with your local environment.

Install the full-fat ArgoCD components. Includes Redis, Dex auth server and UI.
This requires broad permissions in the target cluster for creating these
resources.

```sh
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

Retrieve the generated password for the ArgoCD admin user.

```sh
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo
```

This can be used to access the ArgoCD UI and observe the state of managed
resources.

```sh
kubectl port-forward --namespace argocd service/argocd-server :443
```

This is an interim way to access the ArgoCD UI before creating an Ingress
resource with a loadbalancer, valid TLS and user authentication.

## Deploy Echo Server

We can deploy the ArgoCD Application custom resource into the remote cluster.

```sh
kubectl apply -f kubernetes/argocd/echo-server-application.yaml
```

This will configure ArgoCD to watch the sandbox repository and sync the
application resources to the cluster. The default polling period is 3 minutes.

The configuration watches the `main` branch of the sandbox repository. Changes
will be sync'd automatically.

## Github Actions

This repository has a Github Actions workflow that will test the application
when a pull request is opened. There are also validations of the Kubernetes
manifests.

## Vulnerability Scanning

In the VSCode Remote Containers configuration `anchore/grype` is installed
for vulnerability scanning.

For scanning the application packages

```sh
grype dir:.
```

And for scanning the container images

```sh
IMAGE=$(skaffold build -q | jq -r .builds[].tag)
grype docker:$IMAGE
```

TODO: Add this to Github Actions.

## Releasing Changes

Once pull requests are merged to main the container image will be built and
published to [ghcr.io/hhtpcd/echo-server][ghcr]. The image is tagged with the 
short git SHA of the commit.

A further update to the repo is required to update the application manifests.

The container image tag in `kubernetes/base/deployment.yaml` should be bumped to
the tag that you would like to release.

This change in the manifest will cause ArgoCD to resync the Kubernetes
resources.

## Further Work

More work could be done to

- Run a Kind cluster in pull requests to validate the Kubernetes manifests can
  be applied to the cluster.

[remote-containers]: https://code.visualstudio.com/docs/remote/containers
[ghcr]: https://ghcr.io/hhtpcd/echo-server
