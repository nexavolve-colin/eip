# eip [![Go Report Card](https://goreportcard.com/badge/github.com/kaperys/eip)](https://goreportcard.com/report/github.com/kaperys/eip)

`eip` is a small Go command line tool which allows you to retrieve the public and private IP addresses of your AWS EC2 instances.

## Installation

If you have a [working Go environment](https://golang.org/doc/install#testing), the easiest way to install `eip` is using `go get`.

```bash
go get -u github.com/kaperys/eip
```

If you're not a Go user, you can [download the Linux binary](https://github.com/kaperys/eip/releases) and add it to your path.

```bash
tar -zxvf eip.tar.gz
sudo mv eip /usr/local/bin/
```

## Usage

`eip` supports both public and private IP addresses.

```bash
$ eip --public
216.58.212.110

$ eip --private
192.10.22.33
```

You can filter results using the `--filter` flag. The filter flag accepts any filter [supported by the AWS Go SDK](https://github.com/datacratic/aws-sdk-go/blob/master/service/ec2/api.go#L9532-L9754). Multiple filters are supported and you can provide one or more comma separated values to each filter.

```bash
$ eip --filter tag:SystemGroup=api,app --filter tag:Name=my-ec2-instance --private
10.11.12.13
```

To retrieve the private IP addresses associated with all network interfaces use both the `--private` and `--all` flags. The primary private IP address is always printed first.

```bash
eip --private --all
192.10.22.33
192.11.22.33
```
