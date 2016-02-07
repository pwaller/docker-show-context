#!/bin/bash

# A short utility to show packages used recursively.
# Useful for spotting missed vendorings.

go list -f '{{range .Deps}}{{.}}{{"\n"}}{{end}}' . |
  xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' |
  sort -u
