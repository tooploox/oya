Feature: Importing packs

Background:
   Given I'm in project dir

Scenario: Import a pack
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya import github.com/tooploox/oya-fixtures/pack1"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.1.1
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1

    """

Scenario: Import a pack to other already imported
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1

    task: |
      echo "check"
    """
  When I run "oya import github.com/tooploox/oya-fixtures/pack2"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack2: v1.1.2
      github.com/tooploox/oya-fixtures/pack1: v1.1.1
    Import:
      pack2: github.com/tooploox/oya-fixtures/pack2
      pack1: github.com/tooploox/oya-fixtures/pack1

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
  And I run "oya import github.com/tooploox/oya-fixtures/pack1"
  Then the command succeeds
  And file ./subdir/Oyafile contains
    """
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1

    """

Scenario: Import a pack to Oyafile with other things
  Given file ./Oyafile containing
    """
    Project: project

    task: |
      echo "check"
    """
  When I run "oya import github.com/tooploox/oya-fixtures/pack1"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.1.1
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1

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

Scenario: Import a pack with long name should have lower camelcase alias
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya import github.com/tooploox/oya-fixtures/pack3-and-a-half"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack3-and-a-half: v1.1.0
    Import:
      pack3AndAHalf: github.com/tooploox/oya-fixtures/pack3-and-a-half

    """


Scenario: Import a pack with alias from a parameter
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya import github.com/tooploox/oya-fixtures/pack3-and-a-half --alias pack3_5"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack3-and-a-half: v1.1.0
    Import:
      pack3_5: github.com/tooploox/oya-fixtures/pack3-and-a-half

    """

Scenario: Import a pack and expose it
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya import github.com/tooploox/oya-fixtures/pack3-and-a-half --expose"
  Then the command succeeds
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack3-and-a-half: v1.1.0
    Expose: pack3AndAHalf
    Import:
      pack3AndAHalf: github.com/tooploox/oya-fixtures/pack3-and-a-half

    """

Scenario: Try to expose two packs
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya import github.com/tooploox/oya-fixtures/pack1 --expose"
  And I run "oya import github.com/tooploox/oya-fixtures/pack2 --expose"
  Then the command fails with error matching
  """
  .*already exposed.*
  """
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.1.1
    Expose: pack1
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1

    """
