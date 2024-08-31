#ifndef ENGINE_H
#define ENGINE_H

struct execute_lua_script_from_file_ret {
  int ret;
  const char *message;
};

struct execute_lua_script_from_file_ret
execute_lua_script_from_file(const char *filename, const char *archive_path);

int is_lua_script(const char *filename);

#endif // ENGINE_H
