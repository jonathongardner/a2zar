# File
## golden-archive
### Xar
```
podman run --rm -v $PWD:/foo:z -it --user=root registry.fedoraproject.org/fedora-toolbox:latest /bin/bash
sudo dfn install xar
cd /foo/lfs/golden-archive
find * | sort | xargs xar -cf ../xar/golden-archive.xar # not sure order is mantained...
```