# git-mirror

[![Build Status](https://travis-ci.org/martinp/git-mirror.svg)](https://travis-ci.org/martinp/git-mirror)

A backup tool for GitHub repositories.

## Usage
```
$ git-mirror -h
Usage:
  git-mirror [OPTIONS] USER PATH

Application Options:
  -g, --git=    Path to git executable (default: git)
  -q, --quiet   Only print errors
  -n, --dryrun  Print what would be done and exit

Help Options:
  -h, --help    Show this help message

Arguments:
  USER:         GitHub username
  PATH:         Path where repositories should be mirrored
```
