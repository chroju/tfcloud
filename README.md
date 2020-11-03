WIP: tfcloud
=======

`tfcloud` is a partial [Terraform Cloud](https://www.terraform.io/docs/cloud/index.html) CLI tool.

Notes
-----

This is a command line tool that created for the personal use. So, this tool will not be compatible with all Terraform Cloud API.

Set up
------

`tfcloud` requires a Terraform Cloud token. This corresponds to the description on the [CLI Configuration - Terraform by HashiCorp](https://www.terraform.io/docs/commands/cli-config.html#credentials-1) (like `$HOME/.terraformrc`, `TF_CLI_CONFIG_FILE` environment variable).

Commands
--------

# run

```
$ tfc run list <organization>
```

# workspace

```
$ tfc workspace create
$ tfc workspace list <organization>
$ tfc workspace update
```

# module

```
$ tfc module list
$ tfc module versions <module>
```

LICENSE
-------

[MIT](https://github.com/chroju/tfcloud/blob/main/LICENSE)

Author
------

chroju https://chroju.dev
