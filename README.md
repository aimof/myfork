# Where is my fork?

Search forked repositories which has Issues.
It is mvp so it may be changed distructively.

## Usage

```sh
export MYFORK_TOKEN="your personal access token"
myfork $OWNER/$REPO $OUTPUT_FILE
```

## example

```sh
myfork boltdb/bolt bolt.out 
```

bolt.out

```tsv
2018-07-23 20:12:10 +0000 UTC	0	drancom	https://github.com/drancom/bolt
2019-04-12 16:15:08 +0000 UTC	1956	etcd-io	https://github.com/etcd-io/bbolt
2016-11-07 13:33:27 +0000 UTC	0	ef37	https://github.com/ef37/bolt
2016-08-12 08:21:44 +0000 UTC	0	din982	https://github.com/din982/bolt
2016-03-22 15:37:26 +0000 UTC	0	grubern	https://github.com/grubern/bolt
2015-08-20 14:33:46 +0000 UTC	0	nett55	https://github.com/nett55/bolt-1
2015-07-31 08:51:23 +0000 UTC	0	nettan20	https://github.com/nettan20/bolt-1

```

It seems to be etcd-io/bbolt is the repo you should see.

### tsv format

```tsv
timestamp(updatedAt)	stars	size(kB)	url
```

Its order is newest forked to older forked.

## Filter

* Has Issue