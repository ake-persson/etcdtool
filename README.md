# etcdtool

Export/Import/Edit etcd directory as JSON/YAML/TOML and validate directory using JSON schema.

# Use cases

- Backup/Restore in a format which is not database or version specific.
- Migration of data from production to testing.
- Store configuration in Git and use import to populate etcd.
- Validate directory entries using JSON schema.

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

# Usage etcdtool

```bash
NAME:
   etcdtool - Command line tool for etcd to import, export, edit or validate data in either JSON, YAML or TOML format.

USAGE:
   ./bin/etcdtool [global options] command [command options] [arguments...]
   
VERSION:
   2.6
   
COMMANDS:
   import	import a directory
   export	export a directory
   edit		edit a directory
   validate	validate a directory
   tree		List directory as a tree
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --peers, -p "http://127.0.0.1:4001,http://127.0.0.1:2379"	Comma-delimited list of hosts in the cluster [$ETCDTOOL_PEERS]
   --cert 							Identify HTTPS client using this SSL certificate file [$ETCDTOOL_CERT]
   --key 							Identify HTTPS client using this SSL key file [$ETCDTOOL_KEY]
   --ca 							Verify certificates of HTTPS-enabled servers using this CA bundle [$ETCDTOOL_CA]
   --user, -u 							User
   --timeout, -t "1s"						Connection timeout
   --command-timeout, -T "5s"					Command timeout
   --help, -h							show help
   --version, -v						print the version
```

> You can also set an env. variable for the etcd node and port.

```bash
export ETCDTOOL_PEERS="http://etcd1.example.com:2379"
```

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
