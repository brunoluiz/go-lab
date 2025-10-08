<h1 align="center">
  Go Lab
</h1>

<p align="center">
  🧪 Laboratory for my Go and infrastructure experiments
</p>

# Kustomize

This repository contains a Kustomize architecture for deploying applications. You can find shared bases for deployments, cronjobs, and jobs in `infra/kustomize/shared`.

An example implementation for the `hello-world` service can be found in `services/hello-world/kustomize`.

# Commands

- `make run service=... cmd=...`: runs a binary with hot-reloading
