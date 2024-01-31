package secretrabbit

import (
	"fmt"
	"math"
	"runtime"
	"unsafe"

	"modernc.org/libc"
	"modernc.org/libsamplerate"
)

type Converter int

const (
	SincBestQuality   = Converter(libsamplerate.SRC_SINC_BEST_QUALITY)
	SincMediumQuality = Converter(libsamplerate.SRC_SINC_MEDIUM_QUALITY)
	SincFastest       = Converter(libsamplerate.SRC_SINC_FASTEST)
	ZeroOrderHold     = Converter(libsamplerate.SRC_ZERO_ORDER_HOLD)
	Linear            = Converter(libsamplerate.SRC_LINEAR)
)

func (c Converter) String() string {
	tls := libc.NewTLS()
	defer tls.Close()
	r := libsamplerate.Xsrc_get_name(tls, int32(c))
	if r == 0 {
		return "unknown samplerate converter"
	}
	return libc.GoString(r)
}

func Simple(in []float32, ratio float64, channels int, converter Converter) ([]float32, error) {
	if len(in) == 0 {
		return nil, nil
	}
	// TODO: TLS pool?
	tls := libc.NewTLS()
	defer tls.Close()
	if ratio <= 0 {
		return nil, asErrorf(tls, int32(ErrBadSrcRatio), "ratio=%v", ratio) // 6 = _SRC_ERR_BAD_SRC_RATIO
	}
	out := make([]float32, len(in)*int(math.Ceil(ratio)))

	pin := new(runtime.Pinner)
	defer pin.Unpin()

	srcDataPtr := libc.Xcalloc(tls, 1, libsamplerate.Tsize_t(unsafe.Sizeof(libsamplerate.TSRC_DATA{})))
	defer libc.Xfree(tls, srcDataPtr)
	// tools complain about this use of unsafe.Pointer,
	// but libc.Xmalloc guarantees that it is fine.
	srcData := (*libsamplerate.TSRC_DATA)(unsafe.Pointer(srcDataPtr))
	// For pinning slices, see https://github.com/golang/go/issues/65286
	pin.Pin(&in[0])
	srcData.Fdata_in = uintptr(unsafe.Pointer(&in[0]))
	pin.Pin(&out[0])
	srcData.Fdata_out = uintptr(unsafe.Pointer(&out[0]))
	srcData.Finput_frames = int64(len(in) / channels)
	srcData.Foutput_frames = int64(cap(out) / channels) // TODO: why cap?
	srcData.Fsrc_ratio = ratio

	r := libsamplerate.Xsrc_simple(tls, srcDataPtr, int32(converter), int32(channels))
	if r != libsamplerate.SRC_ERR_NO_ERROR {
		return nil, asError(tls, r)
	}
	return out[:srcData.Foutput_frames_gen*int64(channels)], nil
}

type ErrorCode int

const (
	ErrMallocFailed          = ErrorCode(libsamplerate.SRC_ERR_MALLOC_FAILED)
	ErrBadState              = ErrorCode(libsamplerate.SRC_ERR_BAD_STATE)
	ErrBadData               = ErrorCode(libsamplerate.SRC_ERR_BAD_DATA)
	ErrBadDataPtr            = ErrorCode(libsamplerate.SRC_ERR_BAD_DATA_PTR)
	ErrNoPrivate             = ErrorCode(libsamplerate.SRC_ERR_NO_PRIVATE)
	ErrBadSrcRatio           = ErrorCode(libsamplerate.SRC_ERR_BAD_SRC_RATIO)
	ErrBadProcPtr            = ErrorCode(libsamplerate.SRC_ERR_BAD_PROC_PTR)
	ErrShiftBits             = ErrorCode(libsamplerate.SRC_ERR_SHIFT_BITS)
	ErrFilterLen             = ErrorCode(libsamplerate.SRC_ERR_FILTER_LEN)
	ErrBadConverter          = ErrorCode(libsamplerate.SRC_ERR_BAD_CONVERTER)
	ErrBadChannelCount       = ErrorCode(libsamplerate.SRC_ERR_BAD_CHANNEL_COUNT)
	ErrSincBadBufferLen      = ErrorCode(libsamplerate.SRC_ERR_SINC_BAD_BUFFER_LEN)
	ErrSizeIncompatibility   = ErrorCode(libsamplerate.SRC_ERR_SIZE_INCOMPATIBILITY)
	ErrBadPrivPtr            = ErrorCode(libsamplerate.SRC_ERR_BAD_PRIV_PTR)
	ErrBadSincState          = ErrorCode(libsamplerate.SRC_ERR_BAD_SINC_STATE)
	ErrDataOverlap           = ErrorCode(libsamplerate.SRC_ERR_DATA_OVERLAP)
	ErrBadCallback           = ErrorCode(libsamplerate.SRC_ERR_BAD_CALLBACK)
	ErrBadMode               = ErrorCode(libsamplerate.SRC_ERR_BAD_MODE)
	ErrNullCallback          = ErrorCode(libsamplerate.SRC_ERR_NULL_CALLBACK)
	ErrNoVariableRatio       = ErrorCode(libsamplerate.SRC_ERR_NO_VARIABLE_RATIO)
	ErrSincPrepareDataBadLen = ErrorCode(libsamplerate.SRC_ERR_SINC_PREPARE_DATA_BAD_LEN)
	ErrBadInternalState      = ErrorCode(libsamplerate.SRC_ERR_BAD_INTERNAL_STATE)
)

type Error struct {
	ErrorCode ErrorCode
	// TODO: consider re-exporting all these constants? :up_arrow:
	s string // eagerly initialized because requires a TLS
	x string // extra context
}

func asErrorf(tls *libc.TLS, r int32, msg string, args ...any) Error {
	e := asError(tls, r)
	e.x = fmt.Sprintf(msg, args...)
	return e
}

func asError(tls *libc.TLS, r int32) Error {
	s := libc.GoString(libsamplerate.Xsrc_strerror(tls, int32(r)))
	return Error{ErrorCode: ErrorCode(r), s: s}
}

func (e Error) Error() string {
	if e.x != "" {
		return e.s + " (" + e.x + ")"
	}
	return e.s
}
