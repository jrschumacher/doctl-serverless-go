# Go deployment tool for DigitalOcean Serverless

When using DigitalOcean Serverless, you can deploy your Go functions using `doctl serverless deploy`, however they don't support **private Go modules** or **vendoring** or **building locally**. This tool is a workaround to this limitation.

## How it works

### Deploy

```bash
doctl-serverless-go deploy
```

This tool will run on any Go packages within the DigitalOcean serverless monorepo. Detection of go packages is determined by the presence of a `go.mod` file. If your project doesn't include a `go.mod` file, this tool will not run.

Using the `GOPRIVATE` environment variable in the `package.yaml` file, this tool will clone the repository within the `private` directory and modify the `go.mod` file to include replacement directive for the private module.

> [!NOTE]
> A **checksum** of the new `go.mod` file is stored as part of the filename of the backup `go.mod` file in the form of `go.mod.<checksum>.bak`. This is used to restore the `go.mod` file to its original state.

### Clean

```bash
doctl-serverless-go clean
```

This command will remove the `private` directory and any temporary files created during the deployment process. It will then restore the `go.mod` file to its original state based on the checksum of changed `go.mod` file.

> [!WARNING]
> If the checksum doesn't match the `go.mod` will not be restored and an error will be thrown.

## Installation

```bash
go install github.com/jrschumacher/doctl-serverless-go@latest
```

## Usage

```bash
usage: doctl-serverless-go <command> [<monorepo-path>]

        deploy  Deploy the go monorepo to digitalocean serverless
        clean   Clean the go monorepo after deployment
```
