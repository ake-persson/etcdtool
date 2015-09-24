# etcd-Export/Import/Validate

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

# Usage etcd Export

```bash
Usage of bin/etcd-export:
  -dir="/": etcd directory
  -format="JSON": Data serialization format YAML, TOML or JSON
  -node="": etcd node
  -output="": Output file
  -port="2379": etcd port
  -version=false: Version
```

> You can also set an env. variable for the etcd node and port.

```bash
export ETCD_CONN="http://etcd1.example.com:2379"
```

# Usage etcd Import

```bash
Usage of bin/etcd-import:
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
    	No validate using JSON schema
  -node string
    	etcd node
  -port string
    	etcd port (default "2379")
  -schema string
    	etcd key for JSON schema
  -version
    	Version
```

> You can also provide input by using STDIN.

## Example

```
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

# Usage etcd Validate

```bash
Usage of bin/etcd-validate:
  -dir string
    	etcd directory
  -node string
    	etcd node
  -port string
    	etcd port (default "2379")
  -schema string
    	etcd key for JSON schema
  -version
    	Version
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
