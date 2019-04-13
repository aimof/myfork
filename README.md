# Where is my fork?

## Usage

```sh
export MYFORK_TOKEN="your personal access token"
myfork $OWNER/$REPO $OUTPUT_FILE
```

## example

```sh
$ myfork antonholmquist/jason /dev/stdout
2019-04-12 14:16:17 +0000 UTC,0,https://api.github.com/repos/aimof/jason
2016-04-09 13:19:24 +0000 UTC,0,https://api.github.com/repos/grubern/jason
```

### tsv format

```tsv
timestamp	stars	url
```