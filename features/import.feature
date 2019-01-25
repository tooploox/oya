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
      other: github.com/bilus/oya/other
    """
  When I run "oya import github.com/bilus/oya/next"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Import:
      oya: github.com/bilus/oya/next
      other: github.com/bilus/oya/other
    """

Scenario: Import a pack to empty Oyafile
  Given file ./Oyafile containing
    """
    """
  When I run "oya import github.com/bilus/oya/next"
  Then the command succeeds
  And file ./Oyafile contains
    """

    Import:
      oya: github.com/bilus/oya/next
    """
