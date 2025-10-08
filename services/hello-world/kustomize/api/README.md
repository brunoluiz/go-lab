# Kustomize for hello-world service

This folder contains the kustomize configuration to deploy the `hello-world` service.

## How to use

To generate the kubernetes manifests, run the following command:

```
kubectl kustomize .
```

This will output the generated manifests to stdout.

To apply the manifests to a kubernetes cluster, run:

```
kubectl apply -k .
```
