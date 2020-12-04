tfcloud
=======

![Test](https://github.com/chroju/tfcloud/workflows/Test/badge.svg)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/chroju/tfcloud?style=flat)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/chroju/tfcloud/main)


`tfcloud` is a partial [Terraform Cloud](https://www.terraform.io/docs/cloud/index.html) (and Terraform Enterprise) CLI tool.

Notes
-----

* `tfcloud` is created for a limited purpose and does not intend to implement all Terraform Cloud / Enterprise APIs.
* Since Terraform Cloud and Terraform Enterprise have the same API, it will probably work with Terraform Enterprise, but it has not been confirmed.

Install
-------

## Homebrew

```
$ brew install chroju/tap/tfcloud
```

## Download binary

Download the latest binary from here and put it in your `$PATH` directory.

https://github.com/chroju/tfcloud/releases


Usage
-----

### Authentication

`tfcloud` requires a Terraform Cloud token. This corresponds to the description on the [CLI Configuration - Terraform by HashiCorp](https://www.terraform.io/docs/commands/cli-config.html#credentials-1) (like `$HOME/.terraformrc`, `TF_CLI_CONFIG_FILE` environment variable).

### Commands

#### run

```bash
# Lists up the current all Terraform runs
$ tfcloud run list <organization>

# Approves the specified Terraform run
$ tfcloud run apply <run ID>
```

#### workspace

```bash
# Lists up the all workspace in the organization
$ tfcloud workspace list <organization>

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
# Lists up the all Terraform registry modules
$ tfcloud module list

# Lists up the all registry module versions
$ tfcloud module versions <organization> <provider> <module name>
```

LICENSE
-------

[MIT](https://github.com/chroju/tfcloud/blob/main/LICENSE)
