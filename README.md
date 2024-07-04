# MagicMirror

Usage of the tool is relatively simple, given you know RegEx.

```
Usage:
  magicmirror [OPTIONS]

Application Options:
  -v, --verbose           Show verbose information (twice for debug)
  -q, --quiet             Suppress output
  -c, --strip-color       Strip colorized output
  -t, --trim-leading=     Trim leading directory components
  -p, --process-workers=  Number of process workers to start (default: 25)
  -m, --match-workers=    Number of match workers to start (default: 2)
  -d, --download-workers= Number of download workers to start (default: 2)
  -r, --regex=            Regex pattern to match (can be specified multiple times)
  -u, --url=              URLs to fetch (can be specified multiple times)
  -V, --version           Show version information

Help Options:
  -h, --help              Show this help message
```

### Some examples might be:

- Downloading the xz compressed linux-X.Y.0 and patch-X.Y.Z files only, not every release,
  and every patch and anything else that may be in the directory. You also want to download
  5 files at a time, because you have a big pipe and you wanna show it off. You also don't
  need that hostname and pub directory component. Well, your command might look like this:

```
$ magicmirror -r '^.*/(linux-[0-9]*\.[0-9]*\.0|patch-.*)\.tar\.xz' -u http://cdn.kernel.org/pub/linux/kernel/ -t 2
```

Now that in and of it'self isn't as useful as one might hope. It is useful, but can be more so.
Let's say, for sake of argument that you are a site admin, and one of your tasks is keeping the
mirrors hosted at your site up to date. Traditionally this has been scripted and uses something
like, rsync, httrack, wget, or some such thing. I aim to simplify that process if possible. Now
that being said, wrapping the tool up in a script is still the best way to go, as it allows you
to maintain a single file of mirrors and patterns and let the automation handle the rest.

Well, to that end, let's add some more mirror jobs to our previous command.

```
$ magicmirror -t 2 \
  -r '^.*/(linux-[0-9]*\.[0-9]*\.0|patch-.*)\.tar\.xz' -u http://cdn.kernel.org/pub/linux/kernel/ \
  -r '^.*/([0-9]*\.[0-9]*|snapshots)/.*/.*\.(iso|img)$' -u http://cdn.openbsd.org/pub/OpenBSD/ \
  -r '^.*/NetBSD-[0-9]*\.[0-9]/iso/.*\.(iso|img.gz)$' -u http://cdn.netbsd.org/pub/NetBSD/
```

This works out because our paths all look like: `*/pub/*/...` and there is the current limitation
of the tool. That is being worked around and should be resolved soon.
