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

- Added `Expose` statement, pulling imported tasks into global scope, for example:

        Project: someproject

        Import:
          bar: github.com/foo/bar

        Expose: bar

  With this command the imported tasks are available both under the `bar` alias
  (e.g. `oya run bar.doSomething`) as well as without it (e.g. `oya run
  doSomething`).

  The `oya import` command now has a flag, to expose the import, for example:

        oya import github.com/foo/bar --expose
