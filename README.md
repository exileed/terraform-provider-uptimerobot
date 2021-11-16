# Terraform Provider Uptimerobot

[![Terrafrom Basic Test Process](https://github.com/exileed/terraform-provider-uptimerobot/actions/workflows/test.yml/badge.svg)](https://github.com/exileed/terraform-provider-uptimerobot/actions/workflows/test.yml)

This repository provides both a Terraform provider for the [UptimeRobot](https://uptimerobot.com).

## Getting Started

In order to get started, use [the documentation included in this repository](docs/index.md).


## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.


## Using the provider

Fill this in for each provider

## Contributing

When contributing, please also add documentation to help other users.
Debugging the provider

Debugging is available for this provider through the Terraform Plugin SDK versions 2.0.0. Therefore the plugin can be started with the debugging flag --debug.

For example (using delve as Debugger):

dlv exec --headless ./terraform-provider-my-provider -- --debug

For more information about debugging a provider please see: [Debugger-Based Debugging](https://www.terraform.io/docs/extend/debugging.html#debugger-based-debugging)

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
