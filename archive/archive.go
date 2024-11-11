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
	ret := C.find_file_in_archive(C.CString(archivePath), C.CString(fileName))
	if ret.err > 0 {
		return "", false, errors.New("failed to find file in archive")
	}
	filePath := C.GoString(ret.file_path)
	C.free(unsafe.Pointer(ret.file_path))

	return filePath, ret.found == 1, nil
}

func ExtractFile(archivePath string, filePath string, outputPath string) error {
	ret := C.extract_file_from_archive(C.CString(archivePath), C.CString(filePath), C.CString(outputPath))
	if ret > 0 {
		return errors.New("failed to extract file from archive")
	}

	return nil
}

func Extract(archivePath string, outputPath string) error {
	ret := C.extract_all_from_archive(C.CString(archivePath), C.CString(outputPath))
	if ret > 0 {
		return errors.New("failed to extract archive")
	}

	return nil
}
