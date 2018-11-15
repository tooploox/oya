Feature: Rendering templates

Background:
   Given I'm in project dir

@render
Scenario: Render a template
  Given file ./Oyafile containing
    """
    Module: project
    Values:
      foo: xxx
    """
  Given file ./templates/test.txt containing
    """
    $foo
    """
  When I run "oya render -f ./Oyafile ./templates/test.txt"
  Then the command succeeds
  And file ./test.txt contains
  """
  xxx
  """
