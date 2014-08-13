# check-repo

A tool that allows you to run a command when a file changes in the given
path(s).

## Usage

```
$ check-repo -h
Check if local copies of GitHub repos exists in given path.

    Usage:
      check-repo [-s] <username> <path>
      check-repo -h | --help

    Options:
      -h --help             Show help
      -s --skip-fork        Skip forked repositories
```

## Example

`$ check-repo martinp /home/martin/git`
