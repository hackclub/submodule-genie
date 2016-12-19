# Submodule Genie

This is a small little script which will make sure that a given repositories
submodules are always up to date!

## Usage

You can run `submodule-genie` in two ways, as a regular binary or in a Docker
container.

### Configuration

| Docker Env | CLI Flag | Purpose |
|------------|----------|---------|
| `SG_DIRECTORY` | `--directory`   | The directory which the git repository should be cloned into, or where it is already located |
| `SG_REMOTE` | `--remote`  | The remote of the local repository which you want to use for pushing and pulling |
| `SG_BRANCH` | `--branch`  | The remote of `$SG_REMOTE` which you want to use for pushing and pulling |
| `SG_FORK_OWNER` | `--fork-owner`  | The username of a GitHub user which has write access on `$SG_FORK_REPO`, and will be used to send Pull Requests` |
| `SG_FORK_REPO` | `--fork-repo`  | The repository name of the GitHub repository which will be used to send Pull Requests from |
| `SG_FORK_GIT_REPO` | Not needed  | The git URL of the repository which will be used to send Pull Requests from |
| `SG_OWNER` | `--owner` | The owner of the repository which Pull Requests will be sent to |
| `SG_REPO` | `--repo` | The name of the repository which will be sent to |
| `SG_UPSTREAM_REMOTE` | `--upstream` | The remote of the repository which you want to send Pull Requests to |
| `SG_UPSTREAM_BRANCH` | `--to-branch` | The branch of the repository which you want to send Pull Requests to |
| `SG_TOKEN` | `--token` | A GitHub authentication which has access to the `repo` scope and belongs (usually) to `$FORK_OWNER`. [You can create them here](https://github.com/settings/tokens) |
| `SG_SUBMODULES` | `--submodules` | The paths of the submodules (relative to `$SG_DIRECTORY`) which you want to update, if multple separate them with a space |

It's handy to note that remotes can be provided as URLs, without them being previously added to the repository. For example, `git@github.com:hackclub/lecture-hall` instead of `upstream`.

### Running the binary

If you choose to run `submodule-genie` through its binary form, you will need to
take care of managing ssh keys and scheduling yourself.

The following is a generic example of the usage of the command, from the
perspective of GitHub user @paked, who wants to keep the submodules of Hack
Club's `lecture-hall` repository up to date.

```
submodule-genie --directory /lecture-hall \
    --remote origin \
    --branch master \
    --fork-owner paked \
    --fork-repo lecture-hall \
    --owner hackclub \
    --repo lecture-hall \
    --upstream git@github.com:hackclub/lecture-hall \
    --to-branch master \
    --submodules vendor/hackclub \
    --token <github-auth-token-here>
```

### Running with Docker

The Dockerfile is set up with a cron job that will run the `submodule-genie`
command once a day. So that git will be able to push and pull from inside the
container, **you will need a set of ssh keys in the `.ssh/` folder** named `id_rsa`
and `id_rsa.pub` respectively. These need to be set up with the GitHub account
you will be sending Pull Requests from.

These are the commands which I have been using to build and run the container:

- `docker build -t genie .`
- `docker run --name submodule-genie --rm genie`
