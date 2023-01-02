ghg
=======

[![Test Status](https://github.com/Songmu/ghg/workflows/test/badge.svg?branch=main)][actions]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![Coverage Status](https://codecov.io/gh/Songmu/ghg/branch/main/graph/badge.svg)][codecov]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Songmu/ghg)][PkgGoDev]

[actions]: https://github.com/Songmu/ghg/actions?workflow=test
[codecov]: https://codecov.io/gh/Songmu/ghg
[license]: https://github.com/Songmu/ghg/blob/main/LICENSE
[PkgGoDev]: https://pkg.go.dev/github.com/Songmu/ghg

## Description

Get the executable from github releases

## Installation

    % go install github.com/Songmu/ghg/cmd/ghg@latest

Built binaries are available on gihub releases.
<https://github.com/Songmu/ghg/releases>

## Synopsis

    % ghg get tcnksm/ghr
    % $(ghg bin)/ghr             # you can run the executable
    % ghg get -u Songmu/retry    # upgrade and overwrite
    % ghg get motemen/ghq@v0.7.1 # specify the release version

## Commands

```
Available commands:
  bin      display bin dir
  get      get stuffs
  version  display version
```

## Author

[Songmu](https://github.com/Songmu)
