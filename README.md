# My personal stateful monolith Bonzai™ commander (z)

[![GoDoc](https://godoc.org/github.com/murtaza-u/z?status.svg)](https://godoc.org/github.com/murtaza-u/z)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

I use a personal fork of
[rwxrob/bonzai](https://github.com/rwxrob/bonzai) which adds flag-like
key-value pair arguments to our beloved commander. Go workspaces are
used to point the dependencies to my local forked version. The `go.work`
file is hence intentionally added to Git.

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C z z
```

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.
