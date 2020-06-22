Feature: REPL

Background:
   Given I'm in project dir

Scenario: Successfully run a command
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya repl" interactively
  And I send "touch ./OK" to repl
  And I send "exit" to repl
  Then file ./OK exists

@current
Scenario: Access a value
  Given file ./Oyafile containing
    """
    Project: project

    Values:
      foo: bar
    """
  When I run "oya repl" interactively
  And I send "echo foo: ${Oya[foo]}" to repl
  And I send "exit" to repl
  Then the command outputs text matching
    """
    foo: bar
    """
