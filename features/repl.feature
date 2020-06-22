Feature: REPL

Background:
   Given I'm in project dir

@current
Scenario: Successfully run a command
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya repl"
  And I send "touch ./OK" to repl
  Then file ./OK exists
