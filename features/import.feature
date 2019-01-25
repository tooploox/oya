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

    task: |
      echo "check" 
    """
  When I run "oya import github.com/bilus/oya/next"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Import:
      next: github.com/bilus/oya/next
      other: github.com/bilus/oya/other
    
    task: |
      echo "check" 
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
      next: github.com/bilus/oya/next

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
