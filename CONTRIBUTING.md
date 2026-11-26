# Contributing to linstor-scheduler-extender

The [Piraeus contributing guide](https://github.com/piraeusdatastore/piraeus/blob/master/CONTRIBUTING.md) describes how to report issues, write commits and changelog entries, and open pull requests.
The sections that follow describe what is specific to this repository.

## Target branch

Open pull requests against the `master` branch.

## Build and test

Install [Go](https://go.dev/) in the version declared in `go.mod`.
Run the unit tests:

```
$ go test ./...
```

On every pull request, CI builds the container image.
CI does not run the unit tests.
Run them locally before you open a pull request.

## Deployment manifests

Two deployment methods exist:

* The plain manifests in `deploy/linstor-scheduler.yaml` and `deploy/linstor-scheduler-admission.yaml`.
* The `linstor-scheduler` chart in [helm-charts](https://github.com/piraeusdatastore/helm-charts).

If you add a command-line option that users must set, or if the scheduler extender requires new RBAC permissions, update the plain manifests and open a pull request for the chart.

## Testing on a cluster

To try a change on a cluster, build an image and push it to your registry:

```
$ make upload REGISTRY=registry.example.com/myname TAG=dev
```

By default, the image is built for `linux/amd64` and `linux/arm64`.
To build for a single platform, set `PLATFORMS`, for example `PLATFORMS=linux/amd64`.
The target also tags the image as `latest`.
Replace the image in your deployment with the new image.
