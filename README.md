# `docker-show-context`

Ever wonder why docker pauses when you do `docker build`, and what you can do
about it? You know, when it says `Sending build context to Docker daemon`?

This program shows where time and bytes are spent when building a docker context.

It is based directly on the same logic that Docker itself uses to build the
context.

## Getting started (binaries):

Binaries are available on the
[releases page](https://github.com/pwaller/docker-show-context/releases).
Just grab the binary and put it in your path, then invoke it as
`docker-show-context`. Use at your own risk.

## Getting started (building from source):

This requires vgo. At some point vgo will become a part of Go, but until then you need to obtain the vgo tool:

```
go get golang.org/x/vgo
```

Then run:

```
git clone https://github.com/pwaller/docker-show-context
cd docker-show-context
vgo install -v
```

# What the output looks like

The output looks something like this. It's easy to see now that I accidentally
included some large binary content (`*.deb` and `*.pdf` files in particular),
so I can now go and add those to my `.dockerignore` or delete them.

```
$ cd ~/path/to/project/using/docker
$ docker-show-context
Scanning local directory (in tar / on disk):
  24 / 1057 (62 / 216 MiB) (0.0s elapsed) .. completed

Excluded by .dockerignore: 1033 files totalling 153.98 MiB

Final .tar:
  24 files totalling 61.83 MiB (+ 0.02 MiB tar overhead)
  Took 0.04 seconds to build

Top 10 directories by time spent:
   40 ms: .
    1 ms: example

Top 10 directories by storage:
  61.83 MiB: .
   0.00 MiB: example

Top 10 directories by file count:
   23: .
    1: example

Top 10 file extensions by storage:
  57.10 MiB: 
   4.71 MiB: .exe
   0.01 MiB: .pprof
   0.01 MiB: .md
   0.01 MiB: .go
   0.00 MiB: .sum
   0.00 MiB: .mod
   0.00 MiB: .sh
   0.00 MiB: .gitignore
   0.00 MiB: .dockerignore
```

# Notes about the current behaviour

This documents the current behaviour, which may not be ideal, but it is what it
is for now. Pull requests welcome.

* The amounts shown don't show recursive usage, they just show a single level
  of the directory. (Otherwise, the root would always be the biggest thing).

* Time records the amount of time between `tarFile.Next()` calls. I assume that
  this approximates the amount of time `docker/pkg/archive` spent constructing
  one tar entry. It might not be precise.

* "Total content" shows the uncompressed bytes inside files inside the tar.
  The total amount sent to the docker daemon is this amount plus the tar
  overhead.

* At this moment, only running with the build context root as the current
  working directory is supported, with a dockerfile named `Dockerfile`.
  Pull requests welcome to add parameters, so long as the existing default
  behaviour is preserved.

# How can I use this to make building faster?

Frequently, I find that `docker build` suddenly takes longer than I expect. It
is often the case that I have accidentally included some binaries or something
which I did not intend to include. This the purpose of this tool is to give
visibility into this, taking into account your existing `.dockerignore`,
so that you can improve your `.dockerfile` or delete assets you don't need.

It scratches an itch.

# License

> The MIT License (MIT)
> 
> Copyright (c) 2016-2018 Peter Waller <p@pwaller.net>
> 
> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:
> 
> The above copyright notice and this permission notice shall be included in all
> copies or substantial portions of the Software.
> 
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.
