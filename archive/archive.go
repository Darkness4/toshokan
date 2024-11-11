package archive

/*
#cgo LDFLAGS: -larchive
#include "extract.h"
#include "peek.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

func FindFile(archivePath string, fileName string) (path string, found bool, err error) {
	carchivePath := C.CString(archivePath)
	defer C.free(unsafe.Pointer(carchivePath))

	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.find_file_in_archive(carchivePath, cfileName)
	if ret.err > 0 {
		return "", false, errors.New("failed to find file in archive")
	}
	filePath := C.GoString(ret.file_path)
	C.free(unsafe.Pointer(ret.file_path))

	return filePath, ret.found == 1, nil
}

func ExtractFile(archivePath string, filePath string, outputPath string) error {
	carchivePath := C.CString(archivePath)
	defer C.free(unsafe.Pointer(carchivePath))

	cfilePath := C.CString(filePath)
	defer C.free(unsafe.Pointer(cfilePath))

	coutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(coutputPath))
	ret := C.extract_file_from_archive(carchivePath, cfilePath, coutputPath)
	if ret > 0 {
		return errors.New("failed to extract file from archive")
	}

	return nil
}

func Extract(archivePath string, outputPath string) error {
	carchivePath := C.CString(archivePath)
	defer C.free(unsafe.Pointer(carchivePath))

	coutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(coutputPath))
	ret := C.extract_all_from_archive(carchivePath, coutputPath)
	if ret > 0 {
		return errors.New("failed to extract archive")
	}

	return nil
}

func IsSupported(archivePath string) bool {
	carchivePath := C.CString(archivePath)
	defer C.free(unsafe.Pointer(carchivePath))
	ret := C.is_supported_archive(carchivePath)
	return ret == 1
}
