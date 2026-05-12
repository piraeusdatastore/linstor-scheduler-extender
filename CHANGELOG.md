# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.4] - 2026-05-12

### Fixed

- Skip the `UpToDate` disk-state check for non-DRBD resources. `InspectVolume`
  previously rejected every node hosting a resource whose
  `Volumes[0].State.DiskState` was not `"UpToDate"`. `DiskState` is only
  populated when the resource is built on top of a DRBD layer; storage-only
  resources (e.g. single-replica ZFS without DRBD) report an empty string, so
  pods using such resources were filtered out from every candidate node and
  stayed `Pending` forever. The driver now walks the resource's layer tree and
  only enforces the `UpToDate` check when a DRBD layer is actually present.

## [0.3.3] - 2026-01-12

### Fixed

- Update `KubeSchedulerConfiguration` to the stable `v1` API. The `v1beta3` API
  was removed in Kubernetes 1.27, causing the scheduler to fail with
  `no kind KubeSchedulerConfiguration is registered for version kubescheduler.config.k8s.io/v1beta3`.

### Changed

- Bump `docker/build-push-action` from 5 to 6.
- Bump `google.golang.org/grpc` from 1.47.0 to 1.56.3.
- Update GitHub Actions runner image.
- Add Dependabot configuration.

## [0.3.2] - 2023-04-11

### Changed

- Fixed an issue when inspecting Pods without explicit namespace.

## [0.3.1] - 2023-02-20

### Changed

- Update dependencies

## [0.3.0] - 2023-01-18

### Added

- Introduce linstor-scheduler-admission controller

## [0.2.1] - 2022-05-13

### Added

- Support unbound PVCs if using WaitForFirstConsumer binding mode

## [0.2.0] - 2022-02-25

### Changed

- Moved to Piraeus Organization

## [0.1.0] - 2022-02-07

### Added

- Initial command, import stork scheduler extender and linstor driver code

[Unreleased]: https://github.com/piraeusdatastore/linstor-scheduler-extender/compare/v0.3.4...HEAD
[0.3.4]: https://github.com/piraeusdatastore/linstor-scheduler-extender/compare/v0.3.3...v0.3.4
[0.3.3]: https://github.com/piraeusdatastore/linstor-scheduler-extender/compare/v0.3.2...v0.3.3
[0.3.2]: https://github.com/piraeusdatastore/linstor-scheduler-extender/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/piraeusdatastore/linstor-scheduler-extender/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/piraeusdatastore/linstor-scheduler-extender/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/piraeusdatastore/linstor-scheduler-extender/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/piraeusdatastore/linstor-scheduler-extender/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/piraeusdatastore/linstor-scheduler-extender/releases/tag/v0.1.0
