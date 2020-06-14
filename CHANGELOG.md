# Oya Changelog


## v0.0.20 (Upcoming)

### Changed

- Report an error if a `Require` or `Replace` directive appears in an `Oyafile`
  without a `Project` directive; Oya dependencies are managed in the top-level
  `Oyafile`.

### Fixed

- Ensure non-zero exit code from a command in Oya tasks, including sub-commands,
  propagates to the shell invoking `oya run`, even without `set -e`.
  
  
## v0.0.19 (Released)

### Added

- Added a `oya secrets init` command.


