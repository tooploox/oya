Feature: Error reporting

Background:
  Given I'm in project dir


Scenario: Script error in task
  Given file ./Oyafile containing
    """
    Project: project
    fail: |
      echo "Stdout"
      echo "Stderr" >&2
      exit 1
    """
  When I run "oya run fail"
  Then the command fails
  And the command outputs to stdout
  """
  Task "fail" in ./Oyafile failed with exit code 1.
  --------------------------------------------------------------------------------
  > Stdout
  2> Stderr
  --------------------------------------------------------------------------------

  """

Scenario: Script error in nested task
Scenario: Script exit code is propagated
Scenario: Variable missing in task
Scenario: Variable missing in nested task
Scenario: Variable missing in template
Scenario: Missing pack on run
Scenario: Missing pack on get
Scenario: Invalid argument
