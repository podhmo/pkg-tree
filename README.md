# pkg-tree
show pkg tree

```console
$ pkg-tree
pkgtree: error: required argument 'pkg' not provided
usage: pkgtree [<flags>] <pkg>

dump pkg dependencies

Flags:
  --help                 Show context-sensitive help (also try --help-long and --help-man).
  --ignore-std-pkg       
  --ignore-internal-pkg  
  --disable-show-id      

Args:
  <pkg>  pkg
```

## examples

```console
$  pkg-tree --ignore-std-pkg github.com/golang/dep
github.com/golang/dep #=0
  github.com/golang/dep/vendor/github.com/pkg/errors #=14
  github.com/golang/dep/internal/fs #=28
    github.com/golang/dep/vendor/github.com/pkg/errors #=14
  github.com/golang/dep/gps #=30
    github.com/golang/dep/vendor/github.com/boltdb/bolt #=32
    github.com/golang/dep/vendor/github.com/sdboyer/constext #=37
    github.com/golang/dep/vendor/golang.org/x/sync/errgroup #=39
      github.com/golang/dep/vendor/golang.org/x/net/context #=40
    github.com/golang/dep/vendor/github.com/pkg/errors #=14
    github.com/golang/dep/internal/fs #=28
    github.com/golang/dep/vendor/github.com/nightlyone/lockfile #=102
    github.com/golang/dep/vendor/github.com/jmank88/nuts #=103
      github.com/golang/dep/vendor/github.com/boltdb/bolt #=32
    github.com/golang/dep/gps/paths #=104
    github.com/golang/dep/gps/internal/pb #=105
      github.com/golang/dep/vendor/github.com/golang/protobuf/proto #=106
    github.com/golang/dep/gps/pkgtree #=110
      github.com/golang/dep/vendor/github.com/pkg/errors #=14
      github.com/golang/dep/vendor/github.com/armon/go-radix #=121
    github.com/golang/dep/vendor/github.com/golang/protobuf/proto #=106
    github.com/golang/dep/vendor/github.com/armon/go-radix #=121
    github.com/golang/dep/vendor/github.com/Masterminds/semver #=122
    github.com/golang/dep/vendor/github.com/Masterminds/vcs #=125
  github.com/golang/dep/gps/paths #=104
  github.com/golang/dep/gps/pkgtree #=110
  github.com/golang/dep/vendor/github.com/pelletier/go-toml #=127
```
