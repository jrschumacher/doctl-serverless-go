name: goreleaser

on:
  push:
    branches:
      - main

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.17.x]
    steps:
      -
        name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      -
        name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      -
        name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          # either 'goreleaser' (default) or 'goreleaser-pro'
          distribution: goreleaser
          # 'latest', 'nightly', or a semver
          version: '~> v2'
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          # Your GoReleaser Pro key, if you are using the 'goreleaser-pro' distribution
          # GORELEASER_KEY: ${{ secrets.GORELEASER_KEY }}