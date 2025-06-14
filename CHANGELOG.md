# Changelog

## [v0.3.2](https://github.com/Songmu/ghg/compare/v0.3.1...v0.3.2) - 2025-06-08
- don't update CREDITS in release process by @Songmu in https://github.com/Songmu/ghg/pull/31

## [v0.3.1](https://github.com/Songmu/ghg/compare/v0.3.0...v0.3.1) - 2025-06-08
- introduce Songmu/gitconfig by @Songmu in https://github.com/Songmu/ghg/pull/26
- replace ioutil by @yulog in https://github.com/Songmu/ghg/pull/28
- update deps by @Songmu in https://github.com/Songmu/ghg/pull/29

## [v0.3.0](https://github.com/Songmu/ghg/compare/v0.2.0...v0.3.0) - 2023-01-02
- update release related stuffs by @Songmu in https://github.com/Songmu/ghg/pull/18
- update deps by @Songmu in https://github.com/Songmu/ghg/pull/20
- use os.UserHomeDir instead of go-homedir by @Songmu in https://github.com/Songmu/ghg/pull/21
- migrate to go-github from octokit by @Songmu in https://github.com/Songmu/ghg/pull/22
- quess download url before requesting assets list by @Songmu in https://github.com/Songmu/ghg/pull/23
- detect latest tag without using API by @Songmu in https://github.com/Songmu/ghg/pull/24

## [v0.2.0](https://github.com/Songmu/ghg/compare/v0.1.4...v0.2.0) (2019-01-29)

* apply go modules and update mholt/archiver unarchive method [#17](https://github.com/Songmu/ghg/pull/17) ([zrma](https://github.com/zrma))

## [v0.1.4](https://github.com/Songmu/ghg/compare/v0.1.3...v0.1.4) (2018-07-08)

* update deps [#15](https://github.com/Songmu/ghg/pull/15) ([Songmu](https://github.com/Songmu))
* Detect executables in windows validly [#14](https://github.com/Songmu/ghg/pull/14) ([delphinus](https://github.com/delphinus))

## [v0.1.3](https://github.com/Songmu/ghg/compare/v0.1.2...v0.1.3) (2018-06-24)

* Fix rename error that occurs if the target is a naked binary [#13](https://github.com/Songmu/ghg/pull/13) ([yuuki](https://github.com/yuuki))

## [v0.1.2](https://github.com/Songmu/ghg/compare/v0.1.1...v0.1.2) (2018-01-01)

* introduce goxz [#12](https://github.com/Songmu/ghg/pull/12) ([Songmu](https://github.com/Songmu))
* Update deps [#11](https://github.com/Songmu/ghg/pull/11) ([Songmu](https://github.com/Songmu))

## [v0.1.1](https://github.com/Songmu/ghg/compare/v0.1.0...v0.1.1) (2017-10-10)

* Adjust releasing flow and introduce dep [#10](https://github.com/Songmu/ghg/pull/10) ([Songmu](https://github.com/Songmu))
* Change working directory location [#9](https://github.com/Songmu/ghg/pull/9) ([Songmu](https://github.com/Songmu))

## [v0.1.0](https://github.com/Songmu/ghg/compare/v0.0.3...v0.1.0) (2017-07-18)

* Support  naked binary without being compressed [#8](https://github.com/Songmu/ghg/pull/8) ([Songmu](https://github.com/Songmu))
* Add newline for bin and version command [#6](https://github.com/Songmu/ghg/pull/6) ([syohex](https://github.com/syohex))
* Fix extraction code failure by updating it for archiver 2.0 API [#4](https://github.com/Songmu/ghg/pull/4) ([BlueSkyDetector](https://github.com/BlueSkyDetector))

## [v0.0.3](https://github.com/Songmu/ghg/compare/v0.0.2...v0.0.3) (2016-08-17)

* Copy & set permission if failed to rename [#3](https://github.com/Songmu/ghg/pull/3) ([delphinus35](https://github.com/delphinus35))

## [v0.0.2](https://github.com/Songmu/ghg/compare/v0.0.1...v0.0.2) (2016-07-08)

* follow newer github.com/mholt/archiver [#2](https://github.com/Songmu/ghg/pull/2) ([Songmu](https://github.com/Songmu))

## [v0.0.1](https://github.com/Songmu/ghg/releases/tag/v0.0.1) (2016-07-03)

- initial implement
