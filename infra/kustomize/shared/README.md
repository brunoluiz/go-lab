# Shared Kustomize Bases

This folder contains shared kustomize bases for deploying applications.

## Available bases

- `deployment`: A base for deploying applications as a kubernetes deployment.
- `cronjob`: A base for deploying applications as a kubernetes cronjob.
- `job`: A base for deploying applications as a kubernetes job.

## How to use

To use a shared base, add it to the `bases` section of your `kustomization.yaml` file.

For example:

```
bases:
- ../../../infra/kustomize/shared/deployment
```
