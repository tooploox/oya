# Oya Changelog

## v0.0.20 (Unreleased)

### Fixed

- Ensure non-zero exit code from a command in Oya tasks, including sub-commands,
  propagates to the shell invoking `oya run`, even without `set -e`.
  
### Added 

- REPL, helping build scripts interactively with access to values in .oya files
  and auto-completion, started using `oya repl`, an example session:
  
        oya run repl
        $ echo ${Oya[someValue]}
        foobar
  
