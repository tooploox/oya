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
  When I run "oya run task positional1 positional2 --switch --value=5 --other-switch"
  Then the command succeeds
  And the command outputs to stdout
  """
  positional1
  positional2
  positional1
  positional2
  --switch = true
  --value = 5
  --other-switch = true

  """
