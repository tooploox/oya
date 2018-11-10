Feature: Running hooks

Background:
   Given I'm in project dir

Scenario: Import hooks from vendored idr
  Given file ./Oyafile containing
    """
    Import:
      foo: github.com/test/foo
    """
  And file ./oya/vendor/github.com/test/foo/Oyafile containing
    """
    all: |
      foo=4
      if [ $$foo -ge 3 ]; then
        touch OK
      fi
      echo "Done"
    """
  When I run "oya run foo/all"
  Then the command succeeds
  And the command outputs to stdout
  """
  Done

  """
  And file ./OK exists

@current
Scenario: Import hooks using pack values
  Given file ./Oyafile containing
    """
    Import:
      foo: github.com/test/foo
    """
  And file ./oya/vendor/github.com/test/foo/Oyafile containing
    """
    Values:
      foo: xxx
    all: |
      bar=$foo
      echo $$bar
    """
  When I run "oya run foo/all"
  Then the command succeeds
  And the command outputs to stdout
  """
  xxx

  """
