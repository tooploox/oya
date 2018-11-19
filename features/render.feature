Feature: Rendering templates

Background:
   Given I'm in project dir

Scenario: Render a template
  Given file ./Oyafile containing
    """
    Project: project
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

Scenario: Render a template directory
  Given file ./Oyafile containing
    """
    Project: project
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

Scenario: Render templated paths
  Given file ./Oyafile containing
    """
    Project: project
    Values:
      foo: xxx
      bar: yyy
    """
  Given file ./templates/${foo}.txt containing
    """
    $foo
    """
  Given file ./templates/$bar/${foo}.txt containing
    """
    $bar
    """
  When I run "oya render -f ./Oyafile ./templates/"
  Then the command succeeds
  And file ./xxx.txt contains
  """
  xxx
  """
  And file ./yyy/xxx.txt contains
  """
  yyy
  """
