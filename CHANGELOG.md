# Oya Changelog


## v0.0.20 (Released)

### Added

- Added an `Import.Alias` built-in variable, usable inside packs, containing the
  alias under which the pack was imported. Consider the following import:

    ```yaml
    Import:
       bar: github.com/bilus/foo
    ```
    
    When you run the task below, defined inside the `foo` pack, it'll print "bar":
    
    ```yaml
    whoami: |
       echo ${Oya[Import.Alias]}
    ```
    
- Added build for Apple Silicon processors `arm64`.


### Changed

- Report an error if a `Require` or `Replace` directive appears in an `Oyafile`
  without a `Project` directive; you can manage an Oya project dependencies in
  the top-level `Oyafile`.

### Fixed

- Ensure non-zero exit code from a command in Oya tasks, including sub-commands,
  propagates to the shell invoking `oya run`, even without `set -e`.
  
  
## v0.0.19 (Released)

### Added

- Added a `oya secrets init` command.


