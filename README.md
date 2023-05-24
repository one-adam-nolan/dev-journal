# Dev Journal (DJ)

## Installing

1. `go build -o dj`
1. `sudo mv dj /usr/local/bin`

> If you use `go install` it will put the binary in the $GOBIN, which may make it so the autocompletion does not work properly. It is also possible, it will attempt to install as `dev-journal`, which is a lot of typing!