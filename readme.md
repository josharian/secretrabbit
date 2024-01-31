secretrabbit is an idiomatic, pure Go port of [libsamplerate](http://libsndfile.github.io/libsamplerate/).

It is built atop a [low level transpiled pure Go port](https://pkg.go.dev/modernc.org/libsamplerate) created by [modernc.org/ccgo](https://pkg.go.dev/modernc.org/ccgo).

So far, the [Simple](http://libsndfile.github.io/libsamplerate/api_simple.html) and [Full](http://www.mega-nerd.com/SRC/api_full.html) APIs have been implemented. For the [callback API](http://www.mega-nerd.com/SRC/api_callback.html), ask, or even better, send a PR. :)

### Status

Partial coverage. API is not stable yet. Otherwise ready to use.

### Acknowledgement

Some API design and tests have been inspired by and/or copied from [gosamplerate](https://github.com/dh1tw/gosamplerate), which is a CGo wrapper around libsamplerate.

### License

MIT, go nuts
