# ghm

[![Build Status](https://travis-ci.org/martinp/ghm.svg)](https://travis-ci.org/martinp/ghm)

A backup tool for GitHub repositories.

## Usage
```
$ ghm -h
Usage:
  ghm [OPTIONS] github-user path

Application Options:
  -g, --git=PATH                    Path to git executable (default: git)
  -q, --quiet                       Only print errors
  -n, --dryrun                      Print commands that would be run and exit
  -p, --protocol=[ssh|https|git]    Use the given protocol when mirroring (default: ssh)

Help Options:
  -h, --help                        Show this help message

Arguments:
  github-user:                      GitHub username
  path:                             Path where repositories should be mirrored
```
