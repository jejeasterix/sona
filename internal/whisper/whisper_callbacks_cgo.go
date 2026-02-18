//go:build linux || darwin || windows

package whisper

/*
#include <whisper.h>
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

//export sonaGoProgressCB
func sonaGoProgressCB(handle uintptr, progress int32) {
	h := cgo.Handle(handle)
	cb := h.Value().(*StreamCallbacks)
	if cb.OnProgress != nil {
		cb.OnProgress(int(progress))
	}
}

//export sonaGoSegmentCB
func sonaGoSegmentCB(handle uintptr, ctxPtr unsafe.Pointer, nNew int32) {
	h := cgo.Handle(handle)
	cb := h.Value().(*StreamCallbacks)
	if cb.OnSegment != nil {
		ctx := (*C.struct_whisper_context)(ctxPtr)
		nSegments := int(C.whisper_full_n_segments(ctx))
		for i := nSegments - int(nNew); i < nSegments; i++ {
			seg := Segment{
				Start: int64(C.whisper_full_get_segment_t0(ctx, C.int(i))),
				End:   int64(C.whisper_full_get_segment_t1(ctx, C.int(i))),
				Text:  C.GoString(C.whisper_full_get_segment_text(ctx, C.int(i))),
			}
			cb.OnSegment(seg)
		}
	}
}

//export sonaGoAbortCB
func sonaGoAbortCB(handle uintptr) int32 {
	h := cgo.Handle(handle)
	cb := h.Value().(*StreamCallbacks)
	if cb.ShouldAbort != nil && cb.ShouldAbort() {
		return 1
	}
	return 0
}
