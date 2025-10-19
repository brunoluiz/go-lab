Can you change app/error.go so:

1. It is under a package "lib/errx"
2. The errors should be builders, instead of sentinel errors
3. The calls to the sentinel errors should be replaced with calls to the builders
