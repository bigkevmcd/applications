# applications

Simple Applications CRD for Kubernetes.

## Installation

```console
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml
$ kubectl create -f deploy/crds/app_v1alpha1_application_crd.yaml
$ kubectl create -f deploy/operator.yaml
```

## Creating Applications

```console
$ kubectl create -f deploy/crds/app_v1alpha1_application_cr.yaml
```

## Building from Source

This uses the [`operator-sdk`](https://github.com/operator-framework/operator-sdk) to build, see the [installation instructions](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md) for details on how to install the tooling.

```console
$ operator-sdk build quay.io/example/app-operator
```

And push to your Docker image hosting provider of choice:

```console
$ docker push quay.io/example/app-operator
```

## Removing the operator

```console
$ kubectl create -f deploy/crds/app_v1alpha1_application_cr.yaml
$ kubectl create -f deploy/operator.yaml
$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/service_account.yaml
$ kubectl create -f deploy/role_binding.yaml
$ kubectl create -f deploy/crds/app_v1alpha1_application_crd.yaml
```
