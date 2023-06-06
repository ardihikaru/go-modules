<p align="center">
<a href="https://github.com/satumedishub/octopus-api-service">
<img src="https://raw.githubusercontent.com/gilbarbara/logos/master/logos/docker-icon.svg" width="30%" />
</a>
<br>
</p>
<p align="center">
<a href="#">
<img src="https://img.shields.io/badge/%20Platforms-Windows%20/%20Linux-blue.svg?style=flat-square" alt="Platforms" />
</a>
<a href="https://github.com/satumedishub/octopus-api-service/blob/master/LICENSE">
<img src="https://img.shields.io/badge/%20Licence-MIT-green.svg?style=flat-square" alt="license" />
</a>
</p>
<p align="center">
<a href="https://github.com/satumedishub/octopus-api-service/blob/master/CODE_OF_CONDUCT.md">
<img src="https://img.shields.io/badge/Community-Code%20of%20Conduct-orange.svg?style=flat-squre" alt="Code of Conduct" />
</a>
<a href="https://github.com/satumedishub/octopus-api-service/blob/master/SUPPORT.md">
<img src="https://img.shields.io/badge/Community-Support-red.svg?style=flat-square" alt="Support" />
</a>
<a href="https://github.com/satumedishub/octopus-api-service/blob/master/CONTRIBUTING.md">
<img src="https://img.shields.io/badge/%20Community-Contribution-yellow.svg?style=flat-square" alt="Contribution" />
</a>
</p>
<hr>

# Project Template

Please add project description here.

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

![Github Actions](https://github.com/satumedishub/octopus-api-service/workflows/Github%20Action/badge.svg)

## Table of Contents

* [Dependencies](#dependencies)
* [Prerequisites](#prerequisites)
* [Installation](#installation)
* [Development](#development)
* [Usage](#usage)
* [Configurations](#configurations)
* [Contributing](#contributing)
* [License](#license)
* [Misc](#misc)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

* `Git` installation (**Dude, please...**)
    * Run this command:
      ```shell
      apt-get -y install git
      ```
      Or
      ```shell
      yum -y install git
      ```

* `golang-ci` is one of the **IMPORTANT** packages. Any developer who will maintain this project should install it. The
  installation command is as follows:
    ```shell
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
    ```
* Everytime the code is updated, please run following command:
  ```shell
  golangci-lint run ./...
  ```
    * Whenever you get **RED colored logs**, means that **YOU ARE ABSOLUTELY** need to fix that before finalize your
      work (e.g. *Make a PR*)
        * An example of an error log:
          ```shell
          pkg/router/handlers/user.go:19:6: `xxx` is unused (deadcode)
          type xxx struct {
               ^
          ```
    * The only log that can be ignored is the `WARN`, such as:
      ```shell
      WARN [linters context] structcheck is disabled because of go1.18. You can track the evolution of the go1.18 support by following the https://github.com/golangci/golangci-lint/issues/2649.
      ```

### Installation

* TO BE ADDED

## Development

* *Release 1.0* : **Not Available Yet**.

## Usage

* TO BE ADDED

## Configurations

* TO BE ADDED

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

Looking to contribute to our code but need some help? There's a few ways to get information:

* Connect with me on [Twitter](https://twitter.com/satumedis)
* Connect with me on [Facebook](https://www.facebook.com/profile.php?id=100073637499682)
* Connect with me on [LinkedIn](https://www.linkedin.com/company/76233667)
* Log an issue here on GitHub

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see
the [tags on this repository](https://github.com/satumedishub/octopus-api-service/tags).

## Authors

* **[Muhammad Febrian Ardiansyah](https://github.com/ardihikaru)**
* ADD MORE HERE ...

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Misc

* [Golang CI Lint](https://golangci-lint.run/usage/install/)
    * Installation command:
      ```
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
      ```
    * Usage: `golangci-lint run <folderName>`
        * Example 1: `golangci-lint run cmd/...`
        * Example 2: `golangci-lint run pkg/...`

* using `go-modules` WITH tag
  * Once you finished added updates in this `go-modules`, please finalize the work by adding a new tag
  * create new tag:
      ```
      git tag -a <YOUR-TARGET-VERSION>
      ```
      * for example:
        ```shell
        git tag -a v0.0.3 -m "add `licence_revoked` field in `identities` storage module"
        ```
  * add tag to the remote:
      ```
      git push --tags


* Preparing the `go.mod` and `go.sum` ([ref](https://golangbyexample.com/go-mod-sum-module/))
    * Init: `go mod init github.com/satumedishub/go-modules`
    * Tidy up: `go mod tidy`

* In case we had created a new tag, but we want to re-install the corresponding package with the latest version. We need to delete the old one from our local repository. Here is the command:
    ```
    sudo rm -fr $(go env GOPATH)/pkg/mod/github.com/satumedishub/'go-modules@v0.0.2'/
    ```
    * where `v0.0.2` is the version we want to remove

<p> Copyright &copy; 2022 Private use in within the SatuMedis Development Team. All Rights Reserved.</p>
