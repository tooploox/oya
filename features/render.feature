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
  Given file ./test.tpl containing
    """
    $foo
    """
  When I run "oya render -f ./Oyafile test.tpl"
  Then the command succeeds
  And the command outputs to stdout
  """
  xxx
  """
