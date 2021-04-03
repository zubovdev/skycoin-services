# Tool to download repos from github organization

## Prerequisite

- Go v1.14+
- Git CLI is installed. This github download tool relies on the `git` CLI to download repos so that all git commits are
saved.
- If the organization is private, the git must be configured to be able to access the organization.

## Install

```sh
go install ./...
```

## Example

### Download a public organization

```
# download repos in https://github.com/skycoin org
# go to the dir where you'd like to store the download repos
cd $HOME/Download/skycoin_backup/

gh-downloader -org skycoin
```

### Download a private organization

To download repos of a private organization, the following requirements must be met:

- A github account that has access to the organization
- A github access token, which will be used to access repos in the private organization.

#### Generate access token

Go to the `Settings` of your github account, then `Developer settings` > `Personal access tokens`. Click the `Generate new token` button and select all the options in `repo` scope. Click `Generate token` and save the token, we will use it later.

#### Download the repos

For example, download repos from `github.com/skycoinpro`

```sh
gh-downloader -org skycoinpro -token $TOKEN
```

The value of `$TOKEN` is the access token we generated before.

## Get repo url list only

If we only need to get the repo urls of an organization, for example, get the repo list of github.com/skycoin:

```sh
gh-downloader -org skycoin -urls-only
```

If the organization is **private**, for example, get the repo list of github.com/skycoinpro:

```sh
gh-downloader -org skycoinorg -t $TOKEN -urls-only
```
