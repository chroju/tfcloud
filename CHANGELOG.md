CHANGELOG
====

## v0.2.3 (2022-06-12)

### Others

* Update README.
* Update dependencies.
  * Update mitchellh/cli to 1.1.4 ([#20](https://github.com/chroju/tfcloud/pull/20))
  * Update hashicorp/go-tfe to 1.3.0 ([#21](https://github.com/chroju/tfcloud/pull/21))

## v0.2.2 (2022-05-20)

### Bugfix

* Fix a bug with `workspace view` command that cause it to crash when run in a workspace with no VCSRepo set.

### Others

* Update dependencies.

## v0.2.1 (2022-04-25)

### Bugfix

* Fix Terraform backend config parse error.

## v0.2.0 (2022-04-24)

### ENHANCEMENT

* Add `workspace view` subcommand to view Terraform Cloud workspace details.

### Bugfix

* Fix recursive directory search with `workspace upgrade` command when checking remote backend.
* Fix errors in the URLs displayed by each command.

### Others

* Performance improved.

## v0.1.1 (2020-12-04)

### Bugfix

* Fix some wrong outputs.

## v0.1.0 (2020-12-04)

* Initial release.

## v0.0.1 (2020-11-03)

* Beta release.
