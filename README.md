# GitHub Daily Report

A simple script that generates daily report using github commit and comment api.
You have to set following two envvars to fetch data.

- `ORG_NAME` : Organization name of the repogitory
- `REPO_NME` : Repogitory name which you want to use.

The script also uses GitHub API and a personal access token is required.
Please set in form of envvar named `GITHUB_ACCESS_TOKEN`.

In short, you can set these vars to your shell configuration files suc as `.bashrc` or `.zshrc`.

```shell
export ORG_NAME="Organization name of the repogitory"
export REPO_NME="Repogitory name which you want to use"
export GITHUB_ACCESS_TOKEN="GitHub personal access token"
```

