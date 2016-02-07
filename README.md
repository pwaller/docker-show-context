# `docker-show-context`

Ever wonder why docker pauses when you do `docker build`, and what you can do
about it? You know, when it says `Sending build context to Docker daemon`?

This program shows where time and bytes are spent when building a docker context.

It is based directly on the same logic that Docker itself uses to build the
context.

## Getting started (binaries):

Binaries are available on the
[releases page](https://github.com/pwaller/docker-show-context/releases/tag/v1.0).
Just grab the binary and put it in your path, then invoke it as
`docker-show-context`. Use at your own risk.

## Getting started (building from source):

```
$ go get -v -u github.com/pwaller/docker-show-context
$ # Note: above command might take a minute or two
$ # because docker/docker is vendored and big.
```

# What the output looks like

The output looks something like this. It's easy to see now that I accidentally
included some large binary content (`*.deb` and `*.pdf` files in particular),
so I can now go and add those to my `.dockerignore` or delete them.

```
$ cd ~/path/to/project/using/docker
$ docker-show-context
2016/02/07 13:30:06 Building context...

(note: totals do not take into account dockerignore)
Creating tar: 4417 / 13949 (165 / 2188 MiB) (5.4s elapsed) .. completed

Top 10 directories by time spent:
  139 ms: vendor/github.com/google/go-github/github
   97 ms: vendor/golang.org/x/net/ipv4
   94 ms: vendor/github.com/foo/bar/baz
   85 ms: http
   82 ms: vendor/github.com/stripe/stripe-go
   80 ms: vendor/golang.org/x/net/ipv6
   61 ms: vendor/golang.org/x/net/html/testdata/webkit
   51 ms: vendor/github.com/prometheus/procfs
   51 ms: vendor/github.com/prometheus/client_golang/prometheus
   50 ms: vendor/github.com/aws/aws-sdk-go/internal/features/smoke/emr

Top 10 directories by storage:
 100.01 MiB: .
  13.72 MiB: pkg/format/xlsx/mapping
   4.75 MiB: vendor/golang.org/x/text/collate
   2.82 MiB: vendor/golang.org/x/text/display
   2.33 MiB: http/images
   1.01 MiB: tests/integration-data
   0.92 MiB: vendor/github.com/aws/aws-sdk-go/service/ec2
   0.85 MiB: vendor/golang.org/x/text/encoding/simplifiedchinese
   0.80 MiB: vendor/golang.org/x/text/encoding/japanese
   0.79 MiB: vendor/golang.org/x/text/encoding/traditionalchinese

Top 10 directories by file count:
  102: vendor/github.com/google/go-github/github
   92: vendor/github.com/foo/bar/baz
   80: vendor/golang.org/x/net/ipv4
   78: vendor/golang.org/x/net/ipv6
   45: vendor/golang.org/x/net/html/testdata/webkit
   45: http
   44: vendor/github.com/tealeg/xlsx
   39: vendor/github.com/stripe/stripe-go
   37: vendor/golang.org/x/crypto/ssh
   36: vendor/github.com/prometheus/client_golang/prometheus

Top 10 file extensions by storage:
  80.99 MiB: .deb
  29.38 MiB: .go
  15.43 MiB: .com
  12.95 MiB: .ttx
   7.22 MiB: .json
   5.31 MiB: .pdf
   3.67 MiB: .py
   3.21 MiB: .png
   1.05 MiB: .ttf
   0.68 MiB: .js

2016/02/07 13:30:12 Total files: 4437 total content: 165.48 MiB (+ 3.07 MiB tar overhead)
2016/02/07 13:30:12 Took 5.38 seconds to build tar
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

The MIT License (MIT)

Copyright (c) 2016 Peter Waller <p@pwaller.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
