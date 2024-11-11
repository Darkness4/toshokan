#include "peek.h"

#include <archive.h>
#include <archive_entry.h>
#include <libgen.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct find_file_in_archive_ret find_file_in_archive(const char *archive_path,
                                                     const char *file_name) {
  struct archive *archive;
  struct archive_entry *entry;

  // Initialize the archive object
  archive = archive_read_new();
  archive_read_support_format_all(archive); // Enable support for all formats
  archive_read_support_filter_all(
      archive); // Enable support for all compression methods

  int result = archive_read_open_filename(archive, archive_path, 10240);
  if (result != ARCHIVE_OK) {
    fprintf(stderr, "Error opening archive: %s\n",
            archive_error_string(archive));
    archive_read_free(archive);
    return (struct find_file_in_archive_ret){.err = 1};
  }

  const char *file_path = NULL;

  // Loop through each entry in the archive
  while (archive_read_next_header(archive, &entry) == ARCHIVE_OK) {
    const char *entry_name = archive_entry_pathname(entry);
    char *base_name = basename(
        strdup(entry_name)); // Extract the base name (no directory path)

    // Compare base name of the current entry to the target filename
    if (strcmp(base_name, file_name) == 0) {
      // If the base name matches, store the full file path
      file_path = strdup(entry_name); // Duplicate the full path safely
      free(base_name);                // Clean up the strdup'd base name
      break;                          // Exit loop after finding the file
    }
  }

  // Close the archive
  archive_read_close(archive);
  archive_read_free(archive);

  int found = file_path != NULL;

  return (struct find_file_in_archive_ret){
      .file_path = file_path,
      .found = found,
      .err = 0,
  };
}
