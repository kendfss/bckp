bckp
---

compress and/or extract files/trees from zip archives  

## install
```zsh
make
```
puts it at $GOPATH/bin/bckp and converts this readme into a manpage on osx/linux  

## usage
```
  -a    Create archive instead of directory [default: false]
  -d string
        place you want to keep your archive/backup [default: .] (default ".")
  -n    Unzip each argument into its own directory [default: false]
  -u    Unzip arguments [default: false]
```

## todo
- nested archiving 
