# Usage Etcd Export

```bash
Usage:
  etcd-export [OPTIONS]

Application Options:
  -v, --verbose    Verbose
      --version    Version
  -f, --format=    Data serialization format YAML, TOML or JSON (JSON)
  -o, --output=    Output file (STDOUT)
  -n, --etcd-node= Etcd Node
  -p, --etcd-port= Etcd Port (2379)
  -d, --etcd-dir=  Etcd Dir (/)

Help Options:
  -h, --help       Show this help message
```

You can also set an env. variable for the Etcd node and port.

```bash
export ETCD_CONN="http://etcd1.example.com:2379"
```

# Usage Etcd Import

```bash
Usage:
  etcd-import [OPTIONS]

Application Options:
  -v, --verbose    Verbose
      --version    Version
  -f, --format=    Data serialization format YAML, TOML or JSON (JSON)
  -i, --input=     Input file (STDOUT)
  -n, --etcd-node= Etcd Node
  -p, --etcd-port= Etcd Port (2379)
  -d, --etcd-dir=  Etcd Dir (/)

Help Options:
  -h, --help       Show this help message
```

# Build

```bash
git clone https://github.com/mickep76/etcd-export.git
cd etcd-export
./build
bin/etcd-export --version
```

## Test

First configure Docker on your Linux or Mac OS X host.

```bash
./init-etcd.sh start
eval "$(./init-etcd.sh env)"
bin/etcd-import -i example.json
bin/etcd-export -f toml
```

# Install using Homebrew

```bash
brew tap mickep76/funk-gnarge
brew install etcd-export
```
