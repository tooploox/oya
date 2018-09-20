# Oya

1. Packs for CI/CD tools and workflows ala draft.sh

## Usage

1. Install oya and its dependencies:

        curl https://raw.githubusercontent/bilus/oya/master/scripts/setup.sh | sh

1. Initialize project to use a certain CI/CD tool and workflow. Example:

        oya init jenkins-monorepo

   It boostraps configuration for Jenkins pipelines supporting the 1.a workflow (see Workflows below), an Oyafile and supporting scripts.

1. Run a job:

        oya run "build"

   Right now it won't do anything as there are no buildable directories yet. Let's create one.

1. Create a buildable directory:

        mkdir app1
        cat > Toopfile
        build:

## Workflows

### Repo structure

1. Mono-repo:
   - Each app has its own directory
   - There is a directory/file containing deployment configuration

2. Multi-repo:
   - Each app has its own repo
   - Deployment configurations in its own repo

3. Mix:
   - Some/all apps share repos, some may have their own
   - Deployment configurations in its own repo

> Also submodules tried for NT/Switchboard and eventually ditched.

### Change control

a. cEach environment has its own directory

b. Each environment has its own branch

### Evaluation

| Workflow | Projects       | Pros                                       | Cons                                                    |
|----------|----------------|--------------------------------------------|---------------------------------------------------------|
| 1.a      | Elevar         | "Can share code"                           | Merge order dependent [1]                               |
|----------|----------------|--------------------------------------------|---------------------------------------------------------|
| 1.b      | DTS/CONRAD     | Single checkout                            | Complex automation [2]                                  |
|----------|----------------|--------------------------------------------|---------------------------------------------------------|
| 2.a      |                | Same as 2.b                                | Same as 2.b plus need to detect which directory changed |
|----------|----------------|--------------------------------------------|---------------------------------------------------------|
| 2.b      | NT/Switchboard | Better isolation [3] Simple automation [4] | More process overhead [5]                               |
|----------|----------------|--------------------------------------------|---------------------------------------------------------|
| 3.a      | Pledge?        | Can divide up a project however you like   | Complex automation [2]                                  |
|----------|----------------|--------------------------------------------|---------------------------------------------------------|
| 3.b      | Pledge?        | Simple deployment automation               | Same as 3.a                                             |

* [1] Code gets merged from branch to branch; works for small team.
* [2] Need to detect what changed between commits. Many CI/CD tools allow only one configuration per repo and require coding around the limitations, example: https://discuss.circleci.com/t/does-circleci-2-0-work-with-monorepos/10378/13
* [3] No way to just share code, need to package into libraries. Great for microservices and must have for large teams.
* [4] Just put a CI/CD config into the root.
* [5] No way to just share code, need to package into libraries. Bad for small teams wanting to quickly prototype.

## Model

1. A directory is buildable if it has a Toopfile in toml/zeus format.

1. Toopfile contains list of hooks -> script:
   * `changeset` -- CI/CD tool-specific script outputting list of modified files in buildable directories given the current job name.
     - each path must be normalized and prefixed with `+`
     - cannot be overriden, only valid for top-level Toopfile
     - in the future, you'll be able to override for a buildable directory and use `-` to exclude directories, `+` to include additional ones,
       and use wildcards, this will allow e.g. forcing running tests for all apps when you change a shared directory
     - git diff --name-only origin/master...$branch
     - https://dzone.com/articles/build-test-and-deploy-apps-independently-from-a-mo
     - https://stackoverflow.com/questions/6260383/how-to-get-list-of-changed-files-since-last-build-in-jenkins-hudson/9473207#9473207
   * `<job>` -- script to run for the job, optional dependencies

1. Top-level directory contains CI/CD configuration file, depending on what tool we ue

1. Hooks are inherited from parent.

1. .toopignore lists files whose changes do not trigger build for the containing buildable directory

1. It contains generators based on packs. https://github.com/Flaque/thaum
