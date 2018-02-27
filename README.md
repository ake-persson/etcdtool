# etcdtool

Export/Import/Edit etcd directory as JSON/YAML/TOML and validate directory using JSON schema.

# Use cases

- Backup/Restore in a format which is not database or version specific.
- Migration of data from production to testing.
- Store configuration in Git and use import to populate etcd.
- Validate directory entries using JSON schema.

# Build

```bash
git clone https://github.com/mickep76/etcdtool.git
cd etcdtool
make
```

# Build RPM

Make sure you have Docker configured.

```bash
git clone https://github.com/mickep76/etcdtool.git
cd etcdtool
make rpm
```

# Install using Homebrew on Mac OS X

```bash
brew tap mickep76/funk-gnarge
brew install etcdtool
```

**Update**

```bash
brew update
brew upgrade --all
```

# Example

Make sure you have Docker configured.

**Start etcd:**

```
./init-etcd.sh start
eval "$(./init-etcd.sh env)"
```

**Import some data:**

```
cd examples/host/
etcdtool import /hosts/test1.example.com test1.example.com.json
etcdtool import /hosts/test2.example.com test2.example.com.json
```

**Inspect the content:**

```
etcdtool tree /
etcdtool export /

**Export the content and infer numbers lists to keep original json:**
etcdtool export / --num-infer-list
```

**Validate data with different routes:**

```
etcdtool -d -c etcdtool.toml validate /
etcdtool -d -c etcdtool.toml validate /hosts
etcdtool -d -c etcdtool.toml validate /hosts/test2.example.com
etcdtool -d -c etcdtool.toml validate /hosts/test2.example.com/interfaces
etcdtool -d -c etcdtool.toml validate /hosts/test2.example.com/interfaces/eth0
```

**Import with validation**:

```
etcdtool -d -c etcdtool.toml import -v /hosts/test3.example.com test2.example.com.json
```

**Fix validation error:**

```
etcdtool -d -c etcdtool.toml edit -v -f toml /hosts/test2.example.com
```

```
---
    gw = "1.192.168.0.1"
+++
    gw = "192.168.0.1"
```

**Re-validate data:**

```
etcdtool -d -c etcdtool.toml validate /hosts
```

**Authentication**

These commands will prompt you for the password for the user. Alternatively, you
can pass the password in a file with `--password-file` or `-F`:  

```
cat /path/to/passwordfile
passwordstring

etcdtool -password-file /path/to/passwordfile validate /hosts
```

# Caveats

- etcd doesn't support list's, this is handled by using the index as the key:

**JSON Input:**

```json
{
    "users": [
        { "username": "jblack", "first_name": "John", "last_name": "Blackbeard" },
        { "username": "ltrier", "first_name": "Lars", "last_name": "Von Trier" }
    ]
}
```

**Result in etcd:**

```
users/0/username: jblack
users/0/first_name: John
users/0/last_name: Blackbeard
users/1/username: ltrier
users/1/first_name: Ludwig
users/1/last_name: Von Treimer
```

# TODO
- Add detection of format for import based on file type
