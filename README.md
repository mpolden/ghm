# ghm

[![Build Status](https://travis-ci.org/martinp/ghm.svg)](https://travis-ci.org/martinp/ghm)

A backup tool for GitHub repositories.

## Installation

Configure [GOPATH](https://golang.org/doc/code.html#GOPATH) and run:

    $ go get github.com/martinp/ghm

## Examples

Mirror all public repositories of user *martinp* into `~/git`:

    $ ghm martinp ~/git

Print git commands that would be run:

    $ ghm -n martinp ~/git
    /usr/bin/git -C /home/martin/git/dotfiles.git fetch --prune
    /usr/bin/git clone --mirror git@github.com:martinp/emacs.d.git /home/martin/emacs.d.git
    ...

*NOTE: If the repository doesn't already exist on disk, `git clone` will be run,
otherwise the repository is updated using `git fetch`.*

Speed up mirroring by running git commands concurrently:

    # Updating 52 repositories with concurrency level 1 (default)
    $ time ghm -q martinp ~/git
    ghm martinp ~/git  0.88s user 0.06s system 1% cpu 1:30.71 total

    # with concurrency level 20
    $ time ghm -c 20 -q martinp ~/git
    ghm -c 20 -q martinp ~/git  0.88s user 0.13s system 16% cpu 5.995 total

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
  -s, --skip-fork                   Skip forked repositories
  -c, --concurrency=COUNT           Mirror COUNT repositories concurrently (default: 1)

Help Options:
  -h, --help                        Show this help message

Arguments:
  github-user:                      GitHub username
  path:                             Path where repositories should be mirrored
```
