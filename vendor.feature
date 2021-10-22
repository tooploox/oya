Feature: Vendoring

Background:
   Given I'm in project dir


Scenario: Ignore vendored Oyafiles
  Given file ./Oyafile containing
    """
    Project: project
    all: echo "main"
    """
  And file ./.oya/vendor/github.com/test/foo@v1.0.0/Oyafile containing
    """
    all: echo "vendored"
    """
  When I run "oya run all", awaiting its exit
  Then the command succeeds
  And the command outputs
  """
  main

  """

Scenario: Import tasks from vendored packs
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      foo: github.com/test/foo
    """
  And file ./.oya/vendor/github.com/test/foo/Oyafile containing
    """
    all: |
      foo=4
      if [ $$foo -ge 3 ]; then
        touch OK
      fi
      echo "Done"
    """
  When I run "oya run foo.all", awaiting its exit
  Then the command succeeds
  And the command outputs
  """
  Done

  """
  And file ./OK exists
