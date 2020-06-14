# Oya changelog

## v0.0.20 (Upcoming)

### Fixed

- Ensure non-zero exit code from a command in Oya tasks, including sub-commands,
  propagates to the shell invoking `oya run`, even without `set -e`.
