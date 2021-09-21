# Migration

## Unreleased

* Flag `--name` and `-n` has been replaced by `--wallet` and `-w` respectively.
* The service configuration `wallet-service/config.toml` no longer exists.
* The network configurations are located in the `wallet-service/networks` config folder.
* A new `--network` (shorthand `-n`) has been introduced on `command` and `service run` subcommands.
