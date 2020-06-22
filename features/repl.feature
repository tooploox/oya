Feature: REPL

Background:
   Given I'm in project dir

@current
Scenario: Successfully run a command
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya repl" interactively
  And I send "touch ./OK" to repl
  And I send "exit" to repl
  Then file ./OK exists
