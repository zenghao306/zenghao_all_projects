package hdtcodec

/*
#include <stdlib.h>
#include "./hdtcodec.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func HdtEncodeV0(srcStr string) (string, error) {
	var cstr *C.char = C.CString(srcStr)
	var cdst *C.char
	defer C.free(unsafe.Pointer(cstr))

	// encode srcStr
	length := C.hdt_encode_v0(cstr, &cdst)
	if length < 0 {
		err := errors.New("ENCODE_ERR")
		return "", err
	}

	// return dstStr and release.
	encodedStr := C.GoString(cdst)
	//C.hdt_release(unsafe.Pointer(cstr) )
	C.hdt_release(unsafe.Pointer(cdst))
	return encodedStr, nil
}

func HdtDecodeV0(encStr string) (string, error) {
	var cstr *C.char = C.CString(encStr)
	var cdst *C.char
	defer C.free(unsafe.Pointer(cstr))

	// decode
	length := C.hdt_decode_v0(cstr, &cdst)
	if length < 0 {
		err := errors.New("DECODE_ERR")
		return "", err
	}

	// return dstStr and release.
	decodeStr := C.GoString(cdst)
	//C.hdt_release(unsafe.Pointer(cstr))
	C.hdt_release(unsafe.Pointer(cdst))
	return decodeStr, nil
}
