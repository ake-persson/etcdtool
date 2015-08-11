# Usage

```bash
Usage:
  etcd-export [OPTIONS]

Application Options:
  -v, --verbose    Verbose
      --version    Version
  -f, --format=    Data serialization format YAML, TOML or JSON (YAML)
  -o, --output=    Output file (STDOUT)
  -H, --etcd-host= Etcd Host
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

# Install using Homebrew

```bash
brew tap mickep76/funk-gnarge
brew install mickep76/funk-gnarge/etcd-export
```
