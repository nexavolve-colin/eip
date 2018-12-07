# eip

`eip` is a small command line tool which allows you to retrieve public and private IP addresses of your AWS EC2 instances.

## Installation

If you have Go installed and configured, the easiest way to install this is using `go get`.

```bash
$ go get -u github.com/kaperys/eip
```

If you're not a Go user, you can [download the Linux binary](https://github.com/kaperys/eip/releases) and add it to your path.

```bash
$ tar -zxvf eip.tar.gz
$ sudo mv eip /usr/local/bin/
```

## Usage

`eip` support public and private IP addresses.

```bash
$ eip --public
216.58.212.110

$ eip --private
192.10.22.33
```

You can also filter the results using `--filter`. The filter flag accepts any filter [accepted by the Go AWS SDK](https://github.com/datacratic/aws-sdk-go/blob/master/service/ec2/api.go#L9532-L9754).

```bash
$ eip --filter tag:SystemGroup=api,app --private
10.11.12.13
```
