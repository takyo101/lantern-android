# Lantern on Android

The `lantern-android` repository provides documentation and scripts for
building a basic [flashlight][1] client library for Android devices.

## Prerequisites

* An OSX or Linux box.
* [docker][2].
* [Android Studio][3].
* [Go][4].
* [GNUMake][6].

### Temporal hack

This is an experimental feature so we need to do some minor hacks in order to
test it. We're going to work with the `experimental/lantern-android` branch of
[flashlight-build][5]:

```
mkdir -p $GOPATH/src/github.com/getlantern
cd $GOPATH/src/github.com/getlantern
git clone https://github.com/getlantern/flashlight-build.git
cd flashlight-build
git checkout -b experimental/lantern-android remotes/origin/experimental/lantern-android
```

This is only a temporary hack while we wait for the required changes to hit
upstream.

## Building the Android library

Set the `GOPATH` environment variable to
`$GOPATH/src/github.com/getlantern/flashlight-build` for the current session,
the [flashlight-build][5] repository has everything we need to build the
[flashlight][1] lightweight web proxy:

```
export GOPATH=$GOPATH/src/github.com/getlantern/flashlight-build
```

Now, get the `libflashlight` package using `go get`:

```
go get github.com/getlantern/lantern-android/libflashlight
```

Finally, change directory into
`$GOPATH/src/github.com/getlantern/lantern-android/` and pass the build task to
the `make` command.

```
make
```

This will create a new `app` subdirectory with an example android project. You
may import the contents of the `app` subdirectory into Android Studio to see it
working.

## Testing the example project

(pending)

## Building a stand-alone client binary for Android devices

(pending)

[1]: https://github.com/getlantern/flashlight
[2]: https://www.docker.com/
[3]: http://developer.android.com/tools/studio/index.html
[4]: http://golang.org/
[5]: https://github.com/getlantern/flashlight-build
[6]: http://www.gnu.org/software/make/
