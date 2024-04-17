# Changelog

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
