secretrabbit is an idiomatic, pure Go port of [libsamplerate](http://libsndfile.github.io/libsamplerate/).

It is built atop a [low level transpiled pure Go port](https://pkg.go.dev/modernc.org/libsamplerate) created by [modernc.org/ccgo](https://pkg.go.dev/modernc.org/ccgo).

So far, only the [simple](http://libsndfile.github.io/libsamplerate/api_simple.html) API has been implemented. For the [other APIs](http://libsndfile.github.io/libsamplerate/api.html), ask, or even better, send a PR. :)

### Status

Partial coverage. API is not stable yet. Otherwise ready to use.

### Acknowledgement

Some API design and tests have been inspired by and/or copied from [gosamplerate](https://github.com/dh1tw/gosamplerate), which is a CGo wrapper around libsamplerate.

### License

MIT, go nuts
