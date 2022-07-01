# kubernetes-sandbox

This repository is a sandbox for development, testing and building of a demo
Kubernetes application.

There are prequisites for using this sandbox:
- kind
- kustomize
- skaffold
- kubeconform
- Go

All of these are present in the [VSCode Remote Containers][remote-containers]
configuration in the `.devcontainer/` folder.

## Kubernetes Cluster

The sandbox workflow here uses a Kind cluster to test local changes. This can
be configured locally with the Makefile command

```sh
make start-kube
```

The ArgoCD components should be installed into a remote Kubernetes cluster, and
boostrapped with the Application custom resources in `kubernetes/argocd`.

```sh
kubectl apply -f kubectl/argocd/echo-server-application.yaml
```

This will configure ArgoCD to watch the sandbox repository and sync the
application resources to the cluster. The default polling period is 3 minutes.

The configuration watches the `main` branch of the sandbox repository. Changes
will be sync'd automatically.

## HTTP Echo Server

There is a dummy HTTP echo server written in Go in `main.go`. The application
can be tested by running

```sh
make test
```



## Container Images

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

Install the full-fat ArgoCD components. Includes Redis, Dex auth server and UI.

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

## Github Actions

This repository has a Github Actions workflow that will test the application
when a pull request is opened. There are also validations of the Kubernetes
manifests.

## Merging to Main

Once pull requests are merged to main the container image will be built and
published to [ghcr.io/hhtpcd/echo-server][ghcr]. The image is tagged with the 
short git SHA of the commit.

A further update to the repo is required to update the application manifests.

The container image tag in `kubernetes/base/deployment.yaml` should be bumped to
the tag that you would like to release.

This change in the manifest will cause ArgoCD to resync the Kubernetes
resources.

### Further Work

More work could be done to

- Run a Kind cluster in pull requests to validate the Kubernetes manifests can
  be applied to the cluster.

[remote-containers]: https://code.visualstudio.com/docs/remote/containers
[ghcr]: https://ghcr.io/hhtpcd/echo-server
