#ifndef PEEK_H
#define PEEK_H

struct find_file_in_archive_ret {
  /// File path in the archive. Must be freed by the caller.
  const char *file_path;
  /// Boolean indicating if the file was found.
  int found;
  /// Error code
  int err;
};

// Find a file in the archive.
struct find_file_in_archive_ret find_file_in_archive(const char *archive_path,
                                                     const char *file_name);

#endif // PEEK_H
