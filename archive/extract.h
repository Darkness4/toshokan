#ifndef EXTRACT_H
#define EXTRACT_H

int extract_file_from_archive(const char *archive_filename,
                              const char *target_file_path,
                              const char *output_filename);

int extract_all_from_archive(const char *archive_path, const char *output_dir);

#endif // EXTRACT_H
