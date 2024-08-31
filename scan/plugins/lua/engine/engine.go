//go:generate ./embed.sh
package engine

/*
#cgo LDFLAGS: -llua5.4 -lm
#include "engine.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func ExecuteFromPath(path string, archivePath string) (string, error) {
	// Convert Go string to C string
	cscriptPath := C.CString(path)
	defer C.free(unsafe.Pointer(cscriptPath))

	carchivePath := C.CString(archivePath)
	defer C.free(unsafe.Pointer(carchivePath))

	// Execute the Lua script from the file
	result := C.execute_lua_script_from_file(cscriptPath, carchivePath)
	if result.ret != 0 {
		return "", fmt.Errorf("failed to execute Lua script from file")
	}
	var resultMessage string
	if result.message != nil {
		resultMessage = C.GoString(result.message)
		defer C.free(unsafe.Pointer(result.message))
	}
	return resultMessage, nil
}

func IsLuaScript(path string) bool {
	// Convert Go string to C string
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ret := C.is_lua_script(cpath)
	if ret < 0 {
		fmt.Println("failed to check if the file is a Lua script")
		return false
	}

	// Check if the file is a Lua script
	return ret == 1
}
