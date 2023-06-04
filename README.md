tfcloud
=======

[![Test](https://github.com/chroju/tfcloud/workflows/Test/badge.svg)](https://github.com/chroju/tfcloud/actions/workflows/test.yaml)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/chroju/tfcloud?style=flat)](https://github.com/chroju/tfcloud/releases/latest)
[![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/chroju/tfcloud/main)](https://github.com/chroju/tfcloud/blob/main/go.mod)


`tfcloud` is a [Terraform Cloud](https://www.terraform.io/docs/cloud/index.html) (and Terraform Enterprise) CLI tool.

Notes
-----

* `tfcloud` is a Terraform Cloud (and Terraform Enterprise) CLI tool that is created with a specific purpose in mind and does not intend to implement all Terraform Cloud / Enterprise APIs.
* While `tfcloud` should work with Terraform Enterprise due to the similarity of their APIs, it hasn't been tested in actual Terraform Enterprise environments.

Install
-------

### Homebrew

```
brew install chroju/tap/tfcloud
```

### Download binary

Download the latest binary from the following link and place it in your `$PATH` directory.

https://github.com/chroju/tfcloud/releases


Usage
-----

### Authentication

`tfcloud` requires a Terraform Cloud / Enterprise token for authentication. The method for this is the same as described in [CLI Configuration - Terraform by HashiCorp](https://www.terraform.io/docs/commands/cli-config.html#credentials-1) such as

* `terraform login` command
* `$HOME/.terraformrc`
* `TF_CLI_CONFIG_FILE` environment variable
* `TF_TOKEN_hostname` environment variable

### Output format

`tfcloud` supports table and JSON output formats. The default is table format, but you can change it with the `--format` option. The table format is user-friendly, but does not contain all the information. If you want to get all the information, use the JSON format.

### Commands

```bash

#### run

```bash
# Lists all current Terraform runs
$ tfcloud run list <organization>

# Approves the specified Terraform run
$ tfcloud run apply <run ID>
```

#### workspace

```bash
# Lists all workspaces in the organization
$ tfcloud workspace list <organization>

# View Terraform Cloud workspace details
$ tfcloud workspace view # Read the remote config in the current directory
$ tfcloud workspace view --org <organization> --workspace <workspace name> # You can also specify the organization and workspace name

# Upgrades Terraform cloud workspace terraform version
$ tfcloud workspace upgrade [OPTION]

Notes:
  This command works by reading the remote config in the current directory.
  You must run this command in the directory where the target terraform file resides.
  Or you can specify the target directory with the --root-path option.

Options:
  --upgrade-version, -u    Terraform version to upgrade.
                           It must be in the correct semantic version format like 0.12.1, v0.12.2 .
                           Or you can specify "latest" to automatically upgrade to the latest version.
                           (default: latest)
  --root-path              Terraform config root path. (default: current directory)
  --auto-approve           Skip interactive approval of upgrade.
```

#### module

```bash
# Lists all Terraform registry modules
$ tfcloud module list

# Lists all registry module versions
$ tfcloud module versions <organization> <provider> <module name>
```

LICENSE
-------

[MIT](https://github.com/chroju/tfcloud/blob/main/LICENSE)
