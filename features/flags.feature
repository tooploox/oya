Feature: Running tasks

Background:
   Given I'm in project dir

Scenario: Pass flags and positional arguments to a task
  Given file ./Oyafile containing
    """
    Project: project
    task: |
      bashVariable=42
      $for i, arg in Args:
        echo Args[$i] = $Args[i]
      $end
      echo Flags.switch = $Flags.switch
      echo Flags.value = $Flags.value
      echo bashVariable = $$bashVariable
    """
  When I run "oya run task positional1 positional2 --switch --value=5"
  Then the command succeeds
  And the command outputs to stdout
  """
  Args[0] = positional1
  Args[1] = positional2
  Flags.switch = true
  Flags.value = 5
  bashVariable = 42

  """
