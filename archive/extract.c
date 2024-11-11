#include "extract.h"

#include <archive.h>
#include <archive_entry.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int extract_file_from_archive(const char *archive_filename,
                              const char *target_file_path,
                              const char *output_filename) {
  struct archive *archive;
  struct archive_entry *entry;
  int result;
  FILE *output_file = NULL;

  // Initialize the archive object
  archive = archive_read_new();
  archive_read_support_format_all(archive);
  archive_read_support_filter_all(archive);

  // Open the archive file
  result = archive_read_open_filename(archive, archive_filename, 10240);
  if (result != ARCHIVE_OK) {
    fprintf(stderr, "Error opening archive: %s\n",
            archive_error_string(archive));
    archive_read_free(archive);
    return 1;
  }

  // Loop through each entry in the archive
  while (archive_read_next_header(archive, &entry) == ARCHIVE_OK) {
    const char *file_path = archive_entry_pathname(entry);

    // Check if this is the target file
    if (strcmp(file_path, target_file_path) == 0) {
      printf("Extracting file: %s\n", file_path);

      // Open the output file for writing
      output_file = fopen(output_filename, "wb");
      if (!output_file) {
        perror("Error creating output file");
        archive_read_free(archive);
        return 1;
      }

      // Read the file data and write it to the output file
      const void *buffer;
      size_t size;
      int64_t offset;
      while ((result = archive_read_data_block(archive, &buffer, &size,
                                               &offset)) == ARCHIVE_OK) {
        fwrite(buffer, 1, size, output_file);
      }

      if (result != ARCHIVE_EOF) {
        fprintf(stderr, "Error reading data block: %s\n",
                archive_error_string(archive));
        fclose(output_file);
        archive_read_free(archive);
        return 1;
      }

      printf("File extracted successfully to %s\n", output_filename);
      fclose(output_file);
      break;
    }

    // Skip to the next entry if it's not the target
    archive_read_data_skip(archive);
  }

  if (output_file == NULL) {
    fprintf(stderr, "File not found in archive: %s\n", target_file_path);
  }

  // Close the archive
  archive_read_close(archive);
  archive_read_free(archive);
  return output_file == NULL ? 1 : 0;
}

int extract_all_from_archive(const char *archive_path, const char *output_dir) {
  struct archive *archive;
  struct archive_entry *entry;
  int r;

  // Open the archive for reading
  archive = archive_read_new();
  archive_read_support_format_all(archive);
  archive_read_support_filter_all(archive);

  if (archive_read_open_filename(archive, archive_path, 10240) != ARCHIVE_OK) {
    fprintf(stderr, "Could not open archive: %s\n",
            archive_error_string(archive));
    archive_read_free(archive);
    return 1;
  }

  // Iterate through each entry in the archive and extract it
  while ((r = archive_read_next_header(archive, &entry)) == ARCHIVE_OK) {
    const char *entry_name = archive_entry_pathname(entry);
    char output_path[1024];

    // Build the output path (combine output_dir and entry_name)
    snprintf(output_path, sizeof(output_path), "%s/%s", output_dir, entry_name);

    // Ensure the parent directories of the output path exist
    char *dir_name = strdup(output_path);
    char *last_slash = strrchr(dir_name, '/');
    if (last_slash != NULL) {
      *last_slash = '\0';
      mkdir(dir_name, 0755); // Create the directory (if it doesn't exist)
    }
    free(dir_name);

    // Set the entry's file path and extract the data
    archive_entry_set_pathname(entry, output_path);

    // Extract the file content
    r = archive_read_data_into_fd(archive, fileno(fopen(output_path, "wb")));
    if (r != ARCHIVE_OK) {
      fprintf(stderr, "Failed to extract file %s: %s\n", output_path,
              archive_error_string(archive));
      archive_read_free(archive);
      return 1;
    }

    printf("Extracted: %s\n", output_path);
  }

  // Clean up
  if (r != ARCHIVE_EOF) {
    fprintf(stderr, "Error occurred while reading archive: %s\n",
            archive_error_string(archive));
  }

  archive_read_close(archive);
  archive_read_free(archive);

  return 0;
}
