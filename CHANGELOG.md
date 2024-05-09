# Changelog

## [1.0.0](https://github.com/lukecold/event-driver/compare/v0.3.1...v1.0.0) (2024-05-09)


### ⚠ BREAKING CHANGES

* Move event-driver under honestbank domain ([#55](https://github.com/lukecold/event-driver/issues/55))
* Move event-driver to honestbank domain
* **joiner+transformer:** Make joiner be able to join the map typed contents as is ([#54](https://github.com/lukecold/event-driver/issues/54))
* **joiner+transformer:** Make joiner be able to join the map typed contents as is
* **extensions/google-cloud:** GCSEventStore keep all storage by default ([#53](https://github.com/lukecold/event-driver/issues/53))
* **extensions/google-cloud:** GCSEventStore should keep all storage by default
* **golang+handler:** Use slog for formatted logging ([#47](https://github.com/lukecold/event-driver/issues/47))

### Features

* EventStore support looking up sources by key ([#49](https://github.com/lukecold/event-driver/issues/49)) ([282ba46](https://github.com/lukecold/event-driver/commit/282ba46169525739da1430e43a68532c8d34573b))
* Extensions use the local version of event-driver as dependency ([#51](https://github.com/lukecold/event-driver/issues/51)) ([9f2cf48](https://github.com/lukecold/event-driver/commit/9f2cf48e4ce502d1ab176d6f95612967c6f2bd29))
* **extensions/google-cloud:** GCSEventStore keep all storage by default ([#53](https://github.com/lukecold/event-driver/issues/53)) ([bdf52a6](https://github.com/lukecold/event-driver/commit/bdf52a67636ab31aaa4d08294096686cfd8e3015))
* **extensions/google-cloud:** GCSEventStore should keep all storage by default ([1f26c97](https://github.com/lukecold/event-driver/commit/1f26c97a1ad4c97f2fc912d4f94586b876f13822))
* **golang+handler:** Use slog for formatted logging ([#47](https://github.com/lukecold/event-driver/issues/47)) ([a93ad91](https://github.com/lukecold/event-driver/commit/a93ad915594cd5e1d7a2efaa3f80e3dfae4d9abd))
* **joiner+transformer:** Make joiner be able to join the map typed contents as is ([c41f38a](https://github.com/lukecold/event-driver/commit/c41f38a120d783af8c6372a935a16df9346da383))
* **joiner+transformer:** Make joiner be able to join the map typed contents as is ([#54](https://github.com/lukecold/event-driver/issues/54)) ([1f00c7d](https://github.com/lukecold/event-driver/commit/1f00c7d4ad455c68dbaa8d112c9a66928a146ab5))
* Move event-driver to honestbank domain ([d2f6051](https://github.com/lukecold/event-driver/commit/d2f605114af17ea0ea640fb6c3ee1c04715907ee))
* Move event-driver under honestbank domain ([#55](https://github.com/lukecold/event-driver/issues/55)) ([187f1de](https://github.com/lukecold/event-driver/commit/187f1dedbcd46f4ad7339bb40dbed2545442ee9c))
* Support compressing the data before writing into GCS ([#52](https://github.com/lukecold/event-driver/issues/52)) ([55b8cfc](https://github.com/lukecold/event-driver/commit/55b8cfc6c8cefa2e18cc84ff6588000e6454e99c))

## [0.3.1](https://github.com/lukecold/event-driver/compare/v0.3.0...v0.3.1) (2024-04-17)


### Bug Fixes

* **extensions/google-cloud:** GCSEventStore should return null message when object is not found ([bf35fad](https://github.com/lukecold/event-driver/commit/bf35fad9effc8e505ee6e3df3bef41086534d31a))
* **extensions/google-cloud:** GCSEventStore should return null message when object is not found ([#39](https://github.com/lukecold/event-driver/issues/39)) ([9493321](https://github.com/lukecold/event-driver/commit/9493321eca2194de75de8bc2685694aaaf2516f2))

## [0.3.0](https://github.com/lukecold/event-driver/compare/v0.2.0...v0.3.0) (2024-04-17)


### Features

* Add logging in the handlers ([a9beb8a](https://github.com/lukecold/event-driver/commit/a9beb8aa8db2f068f9f0d6f1b4f1d41451b8f6ee))
* Add logging in the handlers ([#37](https://github.com/lukecold/event-driver/issues/37)) ([8c0ebc4](https://github.com/lukecold/event-driver/commit/8c0ebc4441df2dbed1af20e381583622ab9234c5))

## [0.2.0](https://github.com/lukecold/event-driver/compare/v0.1.0...v0.2.0) (2024-04-03)


### Features

* GCS event store support grouping storage with folders ([8d53f1b](https://github.com/lukecold/event-driver/commit/8d53f1b29b01df3640ecafa942c6df203d7d4ac1))
* GCS event store support grouping storage with folders ([#35](https://github.com/lukecold/event-driver/issues/35)) ([73f0c1d](https://github.com/lukecold/event-driver/commit/73f0c1df0978aa07509c75dc236e77fe9d63ee29))
* Support customizing key extractor in cache ([#36](https://github.com/lukecold/event-driver/issues/36)) ([cac6ea4](https://github.com/lukecold/event-driver/commit/cac6ea4d8f95bf6c9b60a73210d2a358684c42d6))
* Sync the github.com/lukecold/event-driver version to extensions ([#33](https://github.com/lukecold/event-driver/issues/33)) ([88ed501](https://github.com/lukecold/event-driver/commit/88ed50113460b35babf33006eb178363e37c8335))

## [0.1.0](https://github.com/lukecold/event-driver/compare/0.0.3...v0.1.0) (2024-02-10)


### Features

* Add syntactic sugar of skip on conflict ([e3efe6e](https://github.com/lukecold/event-driver/commit/e3efe6ed63ccdb2a883989dd13fac9c96046f42a))
* Add syntactic sugar of skip on conflict ([#21](https://github.com/lukecold/event-driver/issues/21)) ([7c9f305](https://github.com/lukecold/event-driver/commit/7c9f30572408302bd4dc852733698587c3138f80))
* Create knative and google-cloud modules under extensions path ([#28](https://github.com/lukecold/event-driver/issues/28)) ([d657b3e](https://github.com/lukecold/event-driver/commit/d657b3e0e048b25fac1cc7ea0fb3cebe8f5499a2))
* Make pipeline respect the context timeouts ([#32](https://github.com/lukecold/event-driver/issues/32)) ([6981375](https://github.com/lukecold/event-driver/commit/698137570c890b2943f66b80b4cb43683eaa6a73))
