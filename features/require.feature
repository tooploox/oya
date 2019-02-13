Feature: Dependency management

Background:
   Given I'm in project dir

Scenario: Get a specific pack version
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya get github.com/tooploox/oya-fixtures@v1.0.0"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/VERSION contains
    """
    1.0.0

    """
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures: v1.0.0

    """

Scenario: Get the latest pack version
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya get github.com/tooploox/oya-fixtures"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.1.0/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.1.0/VERSION contains
    """
    1.1.0

    """
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures: v1.1.0

    """

Scenario: Get pack from a multi-pack repo
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya get github.com/tooploox/oya-fixtures/pack1@v1.1.1"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/VERSION contains
    """
    1.1.1

    """
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.1.1

    """

Scenario: Fetch only the package, not the entire repo
  Given file ./Oyafile containing
    """
    Project: project
    """
  When I run "oya get github.com/tooploox/oya-fixtures/pack1@v1.1.1"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.1.0/Oyafile does not exist
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.1.0/VERSION does not exist

Scenario: Require pack
  Given file ./Oyafile containing
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures: v1.0.0
    foo: echo "bar"
    """
  When I run "oya run foo"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures@v1.0.0/VERSION contains
    """
    1.0.0

    """

Scenario: Require pack from multi-pack repo
  Given file ./Oyafile containing
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.0.0
    foo: echo "bar"
    """
  When I run "oya run foo"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.0.0/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.0.0/VERSION contains
    """
    1.0.0

    """

Scenario: Require two packs from multi-pack repo
  Given file ./Oyafile containing
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.1.1
      github.com/tooploox/oya-fixtures/pack2: v1.1.2
    foo: echo "bar"
    """
  When I run "oya run foo"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/VERSION contains
    """
    1.1.1

    """
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.1.2/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.1.2/VERSION contains
    """
    1.1.2

    """

Scenario: Get command does not upgrade pack by default
  Given file ./Oyafile containing
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.0.0
      github.com/tooploox/oya-fixtures/pack2: v1.0.0
    foo: echo "bar"
    """
  When I run "oya run foo"
  And I run "oya get github.com/tooploox/oya-fixtures/pack1"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.0.0/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.0.0/VERSION contains
    """
    1.0.0

    """
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.0.0
      github.com/tooploox/oya-fixtures/pack2: v1.0.0
    foo: echo "bar"
    """

Scenario: Upgrade single pack using get command
  Given file ./Oyafile containing
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.0.0
      github.com/tooploox/oya-fixtures/pack2: v1.0.0
    foo: echo "bar"
    """
  When I run "oya run foo"
  And I run "oya get -u github.com/tooploox/oya-fixtures/pack1"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/VERSION contains
    """
    1.1.1

    """
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.0.0/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.0.0/VERSION contains
    """
    1.0.0

    """
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.1.1
      github.com/tooploox/oya-fixtures/pack2: v1.0.0
    foo: echo "bar"

    """

Scenario: Upgrade pack by editing the Require section
  Given file ./Oyafile containing
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.0.0
      github.com/tooploox/oya-fixtures/pack2: v1.0.0
    foo: echo "bar"
    """
  When I run "oya run foo"
  And I modify file ./Oyafile to contain
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.1.1
      github.com/tooploox/oya-fixtures/pack2: v1.0.0
    foo: echo "bar"
    """
  And I run "oya run foo"
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/VERSION contains
    """
    1.1.1

    """

Scenario: Generate requires from imports
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1
    foo: echo "bar"
    """
  And file ./subdir/Oyafile containing
    """
    Import:
      pack2: github.com/tooploox/oya-fixtures/pack2
    """
  When I run "oya run foo"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.1.1/VERSION contains
    """
    1.1.1

    """
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.1.2/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.1.2/VERSION contains
    """
    1.1.2

    """
  And file ./Oyafile contains
    """
    Project: project
    Require:
      github.com/tooploox/oya-fixtures/pack2: v1.1.2
      github.com/tooploox/oya-fixtures/pack1: v1.1.1
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1
    foo: echo "bar"

    """

Scenario: Preserve versions when generating requires from imports
  Given file ./Oyafile containing
    """
    Project: project
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1
    Require:
      github.com/tooploox/oya-fixtures/pack1: v1.0.0
    foo: echo "bar"
    """
  And file ./subdir/Oyafile containing
    """
    Import:
      pack2: github.com/tooploox/oya-fixtures/pack2
    """
  When I run "oya run foo"
  Then the command succeeds
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.0.0/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1@v1.0.0/VERSION contains
    """
    1.0.0

    """
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.1.2/Oyafile exists
  And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2@v1.1.2/VERSION contains
    """
    1.1.2

    """
  And file ./Oyafile contains
    """
    Project: project
    Import:
      pack1: github.com/tooploox/oya-fixtures/pack1
    Require:
      github.com/tooploox/oya-fixtures/pack2: v1.1.2
      github.com/tooploox/oya-fixtures/pack1: v1.0.0
    foo: echo "bar"

    """


# Not supported yet (?)
# Scenario: Require two packs from multi-pack repo by git sha
#   Given file ./Oyafile containing
#     """
#     Project: project
#     Require:
#       github.com/tooploox/oya-fixtures/pack1: aaaa
#       github.com/tooploox/oya-fixtures/pack2: bbbb
#     foo: echo "bar"
#     """
#   When I run "oya run foo"
#   Then the command succeeds
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1/Oyafile exists
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1/VERSION contains
#     """
#     1.0.0

#     """
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2/Oyafile exists
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack2/VERSION contains
#     """
#     1.1.0

#     """

# Scenario: Indirect requirements
#   Given file ./Oyafile containing
#     """
#     Project: project
#     Require:
#       github.com/tooploox/oya-fixtures/pack-requiring-pack1: v1.1.0
#     foo: echo "bar"
#     """
#   When I run "oya run foo"
#   Then the command succeeds
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1/Oyafile exists
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1/VERSION contains
#     """
#     v1.1.0
#     """
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack-requiring-pack1/Oyafile exists
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack-requiring-pack1/VERSION contains
#     """
#     v1.1.0
#     """
#   And file ./Oyafile contains
#     """
#     Project: project
#     Require:
#       github.com/tooploox/oya-fixtures/pack-requiring-pack1: v1.1.0
#       github.com/tooploox/oya-fixtures/pack1: v1.1.0
#     foo: echo "bar"
#     """
# Scenario: Indirectly required higher version
#   Given file ./Oyafile containing
#     """
#     Project: project
#     Require:
#       github.com/tooploox/oya-fixtures/pack-requiring-pack1: v1.1.0
#       github.com/tooploox/oya-fixtures/pack1: v1.0.0
#     foo: echo "bar"
#     """
#   When I run "oya run foo"
#   Then the command succeeds
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1/Oyafile exists
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack1/VERSION contains
#     """
#     v1.1.0
#     """
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack-requiring-pack1/Oyafile exists
#   And file ./.oya/packs/github.com/tooploox/oya-fixtures/pack-requiring-pack1/VERSION contains
#     """
#     v1.1.0
#     """
#   And file ./Oyafile contains
#     """
#     Project: project
#     Require:
#       github.com/tooploox/oya-fixtures/pack-requiring-pack1: v1.1.0
#       github.com/tooploox/oya-fixtures/pack1: v1.1.0  # indirect
#     foo: echo "bar"
#     """

#   # Two different major versions -- different paths
#   # Two different major versions -- same path (conflict)
