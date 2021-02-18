# ghm

![Build Status](https://github.com/mpolden/ghm/workflows/ci/badge.svg)

A backup tool for GitHub repositories.

## Installation

Configure [GOPATH](https://golang.org/doc/code.html#GOPATH) and run:

    $ go get github.com/mpolden/ghm

## Examples

Mirror all public repositories of user *mpolden* into `~/git`:

    $ ghm mpolden ~/git

Print git commands that would be run:

    $ ghm -n mpolden ~/git
    /usr/bin/git -C /home/martin/git/dotfiles.git fetch --prune
    /usr/bin/git clone --mirror git@github.com:mpolden/emacs.d.git /home/martin/emacs.d.git
    ...

*NOTE: If the repository doesn't already exist on disk, `git clone` will be run,
otherwise the repository is updated using `git fetch`.*

Speed up mirroring by running git commands concurrently:

    # Updating 52 repositories with concurrency level 1 (default)
    $ time ghm -q mpolden ~/git
    ghm mpolden ~/git  0.88s user 0.06s system 1% cpu 1:30.71 total

    # with concurrency level 20
    $ time ghm -c 20 -q mpolden ~/git
    ghm -c 20 -q mpolden ~/git  0.88s user 0.13s system 16% cpu 5.995 total

## Usage
```
$ ghm -h
Usage of ghm:
  -a	Skip archived repositories
  -c int
    	Number of repositories to mirror concurrently (default 1)
  -n	Print commands that would be run and exit
  -p string
    	Protocol to use for mirroring [git|https|ssh] (default "ssh")
  -q	Only print errors
  -s	Skip forked repositories

Arguments:
  <github-user>	GitHub username
  <path>	Path where repositories should be mirrored
```
