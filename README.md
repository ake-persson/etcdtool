# Usage Etcd Export

```bash
Usage:
  etcd-export [OPTIONS]

Application Options:
  -v, --verbose    Verbose
      --version    Version
  -f, --format=    Data serialization format YAML, TOML or JSON (YAML)
  -o, --output=    Output file (STDOUT)
  -n, --etcd-node= Etcd Node
  -p, --etcd-port= Etcd Port (2379)
  -d, --etcd-dir=  Etcd Dir (/)

Help Options:
  -h, --help       Show this help message
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

Configure Docker on your Linux or Mac OS X host.

```bash
./init-etcd.sh start
eval "$(./init-etcd.sh env)"
etcd-export
```

# Install using Homebrew

```bash
brew tap mickep76/funk-gnarge
brew install etcd-export
```
