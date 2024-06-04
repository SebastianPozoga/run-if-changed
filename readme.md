# Command run-if-changed

## Overview

The run-if-changed is a utility designed to monitor specified directories for changes and execute a given command when changes are detected. It provides options to include and exclude directories, as well as specify file extensions to be monitored.

## Features

- Monitor specific directories for changes.
- Exclude specific directories from monitoring.
- Specify file extensions to be monitored.
- Execute a command when changes are detected.

## Usage

### Command Line Arguments

The application accepts the following command line arguments:

- `--include <directory>`: Specifies the directories to be monitored. This argument is required.
- `--exclude <directory>`: Specifies the directories to be excluded from monitoring. This argument is optional.
- `--exts <extension>`: Specifies the file extensions to be monitored. If not provided, all file extensions will be monitored. This argument is optional.
- `-- <command>`: Specifies the command to be executed when changes are detected. This argument is required and should be placed after all other arguments.
- `--cache <directory>`: Specifies the directories to persist changes hash.

### Syntax

```sh
run-if-changed --include <directory> [--exclude <directory>] [--exts <extension>] -- <command>
```
like 
```sh
run-if-changed --include myproject --exclude myproject/cache --exts .js -- npm run build
```