# Toshokan \[WIP\]

Toshokan is a OPDS API with indexing and searching capabilities for manga readers like Tachiyomi.

It also integrates a system for tagging similar to LANraragi.

## Features

- OPDS API v1.2
- Plugins to add tags
- PostgreSQL for database.
- Build with high availability.
- zig, rar, tar.gz, tar.xz, tar.bz2, 7z, zip, cbz, cbr... supported formats (whatever supports libarchive)
- LANraragi tags style: `namespace:tag`
- Sort by title, author, date
- Simple grouping of series

## Diagram

```mermaid
flowchart TB

archive --> plugins
plugins --> scan
db --> indexer
scan --> indexer

db --> opds
archive --> opds
```

## Follow the progress of the project

[GitHub Project](https://github.com/users/Darkness4/projects/2)
