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
  Given file ./templates/file.txt containing
    """
    $foo
    """
  When I run "oya render -f ./Oyafile ./templates/file.txt"
  Then the command succeeds
  And file ./file.txt contains
  """
  xxx
  """

@render
Scenario: Render a template directory
  Given file ./Oyafile containing
    """
    Module: project
    Values:
      foo: xxx
      bar: yyy
    """
  Given file ./templates/file.txt containing
    """
    $foo
    """
  Given file ./templates/subdir/file.txt containing
    """
    $bar
    """
  When I run "oya render -f ./Oyafile ./templates/"
  Then the command succeeds
  And file ./file.txt contains
  """
  xxx
  """
  And file ./subdir/file.txt contains
  """
  yyy
  """
