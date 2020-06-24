# Oya changelog

## v0.0.20 (Unreleased)

### Fixed

- Ensure non-zero exit code from a command in Oya tasks, including sub-commands,
  propagates to the shell invoking `oya run`, even without `set -e`.
  
### Added 

- A simple REPL, started using `oya repl`, helping build scripts interactively
  with access to values in .oya files, example session:
  
  ```
  oya run repl
  $ echo ${Oya[someValue]}
  foobar
  ```
