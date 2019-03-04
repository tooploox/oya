Feature: Importing packs

Background:
   Given I'm in project dir

Scenario: Import a pack
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya import github.com/bilus/oya"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Import:
      oya: github.com/bilus/oya

    """

Scenario: Import a pack to other already imported
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      other: github.com/tooploox/oya/other

    task: |
      echo "check"
    """
  When I run "oya import github.com/tooploox/oya/next"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Import:
      next: github.com/tooploox/oya/next
      other: github.com/tooploox/oya/other

    task: |
      echo "check"

    """

Scenario: Import a pack to empty Oyafile
  Given file ./Oyafile containing
    """
    Project: project
    """
  And file ./subdir/Oyafile containing
    """
    """
  When I'm in the ./subdir dir
  And I run "oya import github.com/tooploox/oya/next"
  Then the command succeeds
  And file ./subdir/Oyafile contains
    """
    Import:
      next: github.com/tooploox/oya/next

    """

Scenario: Import a pack to Oyafile with other things
  Given file ./Oyafile containing
    """
    Project: project

    task: |
      echo "check"
    """
  When I run "oya import github.com/bilus/oya"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Import:
      oya: github.com/bilus/oya

    task: |
      echo "check"

    """

Scenario: Import a pack which is already imported
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      oya: github.com/bilus/oya

    task: |
      echo "check"
    """
  When I run "oya import github.com/bilus/oya"
  Then the command fails with error matching
  """
  .*Pack already imported: github.com/bilus/oya.*
  """
