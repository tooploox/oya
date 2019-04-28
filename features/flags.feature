Feature: Running tasks

Background:
   Given I'm in project dir

Scenario: Pass flags and positional arguments to a task
  Given file ./Oyafile containing
    """
    Project: project
    task: |
      for i in $*; do
        echo $i
      done

      echo ${Oya[Args.0]}
      echo ${Oya[Args.1]}

      echo --switch = ${Oya[Flags.switch]}
      echo --value = ${Oya[Flags.value]}
      echo --other-switch = ${Oya[Flags.otherSwitch]}
    """
  When I run "oya run task --switch positional1 --value=5 positional2 --other-switch"
  Then the command succeeds
  And the command outputs
  """
  --switch
  positional1
  --value=5
  positional2
  --other-switch
  positional1
  positional2
  --switch = true
  --value = 5
  --other-switch = true

  """

Scenario: Have args survive 'set' command (regression)
  Given file ./Oyafile containing
    """
    Project: project
    task: |
      set -e
      echo $1

    """
  When I run "oya run task arg"
  Then the command succeeds
  And the command outputs
  """
  arg

  """
