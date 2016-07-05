ghg
=======

[![Build Status](https://travis-ci.org/Songmu/ghg.png?branch=master)][travis]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/ghg?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/ghg
[coveralls]: https://coveralls.io/r/Songmu/ghg?branch=master
[license]: https://github.com/Songmu/ghg/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/ghg

## Description

Get the executable from github releases

## Installation

    % go get github.com/Songmu/ghg/cmd/ghg

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
