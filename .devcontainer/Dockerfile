# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.241.1/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.18, 1.17, 1-bullseye, 1.18-bullseye, 1.17-bullseye, 1-buster, 1.18-buster, 1.17-buster
ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        httpie \
        curl \
        jq

RUN curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.12.0/kind-linux-amd64 \
    && install -o root -g root -m 755 ./kind /usr/local/bin/kind \
    && rm -v ./kind

RUN curl -Lo ./kustomize.tar.gz https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv4.5.5/kustomize_v4.5.5_linux_amd64.tar.gz \
    && tar -xzf ./kustomize.tar.gz \
    && install -o root -g root -m 755 ./kustomize /usr/local/bin/kustomize \
    && rm -v ./kustomize.tar.gz

RUN curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 \
    && install -o root -g root -m 755 ./skaffold /usr/local/bin/skaffold \
    && rm -v skaffold

RUN go install github.com/yannh/kubeconform/cmd/kubeconform@v0.4.13
