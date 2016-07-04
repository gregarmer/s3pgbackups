# s3pgbackups

PostgreSQL cluster backups to Amazon S3.

[![Build Status](https://travis-ci.org/gregarmer/s3pgbackups.svg?branch=master)](https://travis-ci.org/gregarmer/s3pgbackups)

## Overview

This application will backup your PostgreSQL cluster to Amazon S3, which includes:

* Separate SQL pg_dump of each database
* Compressed data (gzip)
* Automatically rotated backups
* Ability to exclude specific tables, or entire databases

## Installation

### Debian

```bash
$ curl -s https://packagecloud.io/install/repositories/gregarmer/packages/script.deb.sh | sudo bash
```

### CentOS (or RHEL / other rpm based distro's)

```bash
$ curl -s https://packagecloud.io/install/repositories/gregarmer/packages/script.rpm.sh | sudo bash
```

## Getting Started

1. Follow the installation steps above.
2. Configure ~/.s3pgbackups or your-config.cfg - something like this:

    ```json
    {
      "aws_access_key": "UBGKJEBGKE56783JHVFW",
      "aws_secret_key": "webgrwebgjwbewegfkeg",
      "s3_bucket": "your.db.hostname",
      "s3_rotate_old": true,
      "pg_username": "username",
      "pg_password": "password",
      "pg_sslmode": true,
      "excludes": ["postgres", "template0", "*.django_session", "db1.table3"],
    }
    ```

    > Note: Running s3pgbackups without a config will create a blank config.
3. Run `s3pgbackups -c your-config.cfg` or alternatively `s3pgbackups`

## Usage

```
$ s3pgbackups -h
Usage of dist/s3pgbackups:
  -c string  path to the config file (default "~/.s3pgbackups")
  -n         don't actually do anything, just print what would be done
  -v         be verbose
```

`-c config.cfg` is optional. If you don't specify a config it'll default to
`~/.s3pgbackups` and will create this file if it doesn't already exist.

You can use the `-n` parameter to put s3pgbackups into no-op mode, where no
actions will actually be performed.

The `-v` parameter will make s3pgbackups be verbose about what is actually
happening.

## TODO

- [ ] Support for other RDBMS's
- [ ] Support for other targets, besides Amazon S3
- [ ] If we do the above, we should probably rename the project ?
- [x] Using S3 behind a proxy - depends on [this PR](https://github.com/goamz/goamz/pull/33/files#diff-7db3cb93944d57b2ffc803281c906018R1004)
