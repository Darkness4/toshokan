#include "engine.h"
#include "msgpack.lua.h"
#include <lua5.4/lauxlib.h>
#include <lua5.4/lua.h>
#include <lua5.4/lualib.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int msgpack_openf(lua_State *L) {
  /* Get the module name as passed to luaL_requiref. */
  char const *const modname = lua_tostring(L, 1);

  int res;

  /*
  Check if we know the module. We can use this function to load many
  different Lua modules uniquely identified by modname.
  */
  if (strcmp(modname, "msgpack") == 0) {
    /*
    Parses the Lua source code and leaves the compiled function on the top
    of the stack if there are no errors.
    */
    res = luaL_loadbufferx(L, msgpack_lua_data, sizeof(msgpack_lua_data) - 1,
                           "msgpack", "t");
  } else {
    /* Unknown module. */
    return luaL_error(L, "unknown module \"%s\"", modname);
  }

  /* Check if the call to luaL_loadbufferx was successful. */
  if (res != LUA_OK) {
    return lua_error(L);
  }

  /*
  Runs the Lua code and returns whatever it returns as the result of openf,
  which will be used as the value of the module.
  */
  lua_call(L, 0, 1);
  return 1;
}

int do_execute_lua_script_from_file(lua_State *L, const char *filename) {
  if (filename == NULL) {
    fprintf(stderr, "Error: Invalid filename\n");
    return -1;
  }

  int ret = luaL_dofile(L, filename);
  if (ret != LUA_OK) {
    fprintf(stderr, "Error: Could not execute Lua script, exited with: %d\n",
            ret);
    const char *error = lua_tostring(L, -1);
    fprintf(stderr, "Error: %s\n", error);
    return -1;
  }

  return 0;
}

struct execute_lua_script_from_file_ret
execute_lua_script_from_file(const char *filename, const char *archive_path) {
  lua_State *L = luaL_newstate();
  struct execute_lua_script_from_file_ret ret;
  if (L == NULL) {
    fprintf(stderr, "Error: Could not create Lua state\n");
    return ret;
  }

  // Load the Lua standard libraries
  luaL_openlibs(L);

  luaL_requiref(L, "msgpack", msgpack_openf, 1);

  lua_pushstring(L, archive_path);
  lua_setglobal(L, "archive_path");

  ret.ret = do_execute_lua_script_from_file(L, filename);

  // Get return value from Lua script
  if (!lua_isstring(L, -1)) {
    fprintf(stderr, "Expected string return value from Lua script\n");
    ret.ret = -1;
    goto end;
  }

  // Retrieve the return value from the stack
  const char *result = lua_tostring(L, -1);

  // Duplicate the string so it can be returned
  ret.message = strdup(result);

// Clean up
end:
  lua_close(L);
  return ret;
}

int is_lua_script(const char *filename) {
  lua_State *L = luaL_newstate(); // Create a new Lua state
  if (L == NULL) {
    fprintf(stderr, "Failed to create Lua state.\n");
    return -1; // Indicates an error in setting up Lua
  }

  int status = luaL_loadfile(L, filename);
  lua_close(L);

  return status == LUA_OK;
}
