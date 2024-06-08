# Design Document

Toshokan aims to provide a OPDS API for Tachiyomi and other manga readers. The API will follows the same as LANraragi's API, which allows for easy integration with Tachiyomi.

## Database

SQL is litteraly made for indexing and searching. The question is which SQL database to use. SQLite will be used for the first version knowing fully well its limitations. Other databases can be added in the future to allow for more scalability.

## Plugin System

We will use Lua for the plugin system. For now, we will simply implement `get_tags`.

To configure the plugin, for now, environment variables will be used: `TOSHOKAN_<PLUGIN_NAME>_<VARIABLE_NAME>`.

Utilities will be provided to make it easier to write plugins like extracting the archive, etc.

## OPDS API

The API will be based on the [OPDS 1.2](https://specs.opds.io/opds-1.2) specification.

## Project Tracking

[GitHub Project](https://github.com/users/Darkness4/projects/2/views/1)
