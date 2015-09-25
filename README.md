# etcd-Export/Import/Validate/Delete/Tree

Export/Import etcd directory as JSON/YAML/TOML and validate directory using JSON schema.

# Use cases

- Backup/Restore in a format which is not database or version specific.
- Migration of data from production to testing.
- Store authorative configuration in Git and use etcd-import to populate etcd.
- Copy data from one directory to another.
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

# Usage etcd-export

```bash
Usage of ./bin/etcd-export:
  -dir string
    	etcd directory (default "/")
  -format string
    	Data serialization format YAML, TOML or JSON (default "JSON")
  -output string
    	Output file
  -peers string
    	Comma separated list of etcd nodes (default "127.0.0.1:4001,127.0.0.1:2379")
  -version
    	Version
```

> You can also set an env. variable for the etcd node and port.

```bash
export ETCD_PEERS="http://etcd1.example.com:2379"
```

# Usage etcd-import

```bash
Usage of ./bin/etcd-import:
  -delete
    	Delete entry before import
  -dir string
    	etcd directory
  -force
    	Force delete without asking
  -format string
    	Data serialization format YAML, TOML or JSON (default "JSON")
  -input string
    	Input file
  -no-validate
    	Skip validation using JSON schema
  -peers string
    	Comma separated list of etcd nodes (default "127.0.0.1:4001,127.0.0.1:2379")
  -schema string
    	etcd key for JSON schema
  -version
    	Version
```

> You can also provide input by using STDIN.

## Example usin schemas

```bash
./init-etcd.sh start
eval $(./init-etcd.sh env)
etcdctl mkdir /schemas
etcdctl set /schemas/ntp "$(cat examples/ntp/ntp_schema.json)"
bin/etcd-import -input examples/ntp/routes.json -dir /routes -no-validate
bin/etcd-import -input examples/ntp/ntp-site1.json -dir /ntp/site1
bin/etcd-import -input examples/ntp/ntp-site2.json -dir /ntp/site2
bin/etcd-export -dir /ntp
bin/etcd-validate -dir /ntp
```

## Example using templates

```bash
etcdctl mkdir /templates
etcdctl set /templates/ntp "$(cat examples/ntp/ntp_template.json)"
bin/etcd-edit -dir /ntp/site3 -new
```

# Usage etcd-validate

```bash
Usage of ./bin/etcd-validate:
  -dir string
    	etcd directory
  -peers string
    	Comma separated list of etcd nodes (default "127.0.0.1:4001,127.0.0.1:2379")
  -schema string
    	etcd key for JSON schema
  -version
    	Version
```

# Usage etcd-delete

```bash
Usage of ./bin/etcd-delete:
  -dir string
    	etcd directory
  -force
    	Force delete without asking
  -peers string
    	Comma separated list of etcd nodes (default "127.0.0.1:4001,127.0.0.1:2379")
  -version
    	Version
```

# Usage etcd-tree

```bash
Usage of bin/etcd-tree:
  -dir string
    	etcd directory
  -peers string
    	Comma separated list of etcd nodes (default "127.0.0.1:4001,127.0.0.1:2379")
  -version
    	Version
```

# Usage etcd-edit

```
Usage of etcd-edit:
  -delete
    	Delete entry before import
  -dir string
    	etcd directory (default "/")
  -editor string
    	Editor (default "vim")
  -force
    	Force delete without asking
  -format string
    	Data serialization format YAML, TOML or JSON (default "JSON")
  -new
    	Create new directory entry using template
  -no-validate
    	Skip validation using JSON schema
  -peers string
    	Comma separated list of etcd nodes (default "http://etcd1:5001")
  -schema string
    	etcd key for JSON schema
  -template string
    	etcd key for template
  -tmp-file string
    	Temporary file (default ".etcd-edit.swp")
  -version
    	Version
```

> You can also set an env. variable for the editor.

```bash
export EDITOR=emacs
```

# Build

```bash
git clone https://github.com/mickep76/etcd-export.git
cd etcd-export
./build
bin/etcd-export --version
```

# Build RPM

```bash
sudo yum install -y rpm-build
make rpm
sudo rpm -i etcd-export-<version>-<release>.rpm
```

## Test

First configure Docker on your Linux or Mac OS X host.

```bash
./init-etcd.sh start
eval "$(./init-etcd.sh env)"
bin/etcd-import -i example.json
bin/etcd-export -f toml
bin/etcd-export | bin/etcd-import -dir /test
```

# Install using Homebrew on Mac OS X

```bash
brew tap mickep76/funk-gnarge
brew install etcd-export
```
