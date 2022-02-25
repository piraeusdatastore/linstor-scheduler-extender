# linstor-scheduler-extender

![Latest release](https://img.shields.io/github/v/release/piraeusdatastore/linstor-scheduler-extender)

LINSTOR scheduler extender plugin for Kubernetes which allows a storage driver to give the Kubernetes scheduler hints about where to place a new pod so that it is optimally located for storage performance.

## Get started

If you want to get started directly with an existing LINSTOR setup, check out the [single file deployment](./deploy/all.yaml)
The deployment will create:

* All needed RBAC resources
* A Deployment spawning 2 replicas of the Kube-scheduler with linstor-scheduler-existing, configured to connect to `http://piraeus-op-cs.default.svc`

Copy the file, make any desired changes (see the [options](#options) below) and apply:

```console
$ kubectl apply -f deploy/all.yaml
configmap/linstor-scheduler created
deployment.apps/linstor-scheduler created
serviceaccount/linstor-scheduler created
clusterrole.rbac.authorization.k8s.io/linstor-scheduler created
role.rbac.authorization.k8s.io/linstor-scheduler created
clusterrolebinding.rbac.authorization.k8s.io/linstor-scheduler created
rolebinding.rbac.authorization.k8s.io/linstor-scheduler created

$ kubectl -n kube-system get pods -l app.kubernetes.io/name=linstor-scheduler
NAME                                 READY   STATUS    RESTARTS   AGE
linstor-scheduler-6bb88fc66c-bq4fm   2/2     Running   0          52s
linstor-scheduler-6bb88fc66c-tllgm   2/2     Running   0          52s
```

## Usage

To make your applications using linstor-scheduler, use the `schedulerName: linstor` field.
For example, Pod Templates in a StatefulSet should look like:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: my-stateful-app
spec:
  serviceName: my-stateful-app
  selector:
    matchLabels:
      app.kubernetes.io/name: my-stateful-app
  template:
    metadata:
      labels:
        app.kubernetes.io/name: my-stateful-app
    spec:
      schedulerName: linstor
    ...
```

## Configuration

`linstor-scheduler-extender` uses the environment variables specified in the [`golinstor` library](https://pkg.go.dev/github.com/LINBIT/golinstor/client#NewClient)
for configuration.

| Variable              | Description                                                         |
|-----------------------|---------------------------------------------------------------------|
| `LS_CONTROLLERS`      | A comma-separated list of LINSTOR controller URLs to connect to.    |
| `LS_USERNAME`         | Username to use for HTTP basic auth.                                |
| `LS_PASSWORD`         | Password to use for HTTP basic auth.                                |
| `LS_ROOT_CA`          | CA certificate to use for authenticating the server.                |
| `LS_USER_KEY`         | TLS key to use for authenticating the client to the server.         |
| `LS_USER_CERTIFICATE` | TLS certificate to use for authenticating the client to the server. |


The linstor-scheduler-extender itself can be configured using the following flags:

```
--verbose      Enable verbose logging
```

