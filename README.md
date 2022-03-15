# `k8s-context-i`

This is a tiny little command-line tool, based on [`git branch-i`](
https://github.com/JoelOtter/git-branch-i
), for switching Kubernetes contexts. You'll need `kubectl` installed.

## Installation

Install it using Go, like so:

```sh
go install github.com/JoelOtter/k8s-context-i@latest
```

Ensure your Go directory is on your system path. You might want to alias this
tool to something easier to type - I use `kx`.

## Usage

* Contexts can be navigated using the arrow keys, j and k, Pg Up/Down, or
Ctrl+N and Ctrl+P.
* Switch to a context with the return key.
* Exit with Escape or Ctrl+C.
