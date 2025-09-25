# A to Z archiver

Pure go readers for those forgotten archives. It is inspired by the native go tar package.

## Examples
Examples can be found [here](example/README.md)

## Dev
### Setup
Git LFS is used. Ensure it is installed with:
```
git lfs --version
```
To pull lfs files:
```
git lfs install
```
## Test
### Files
All archive reader should have a golden-archive test. This is a set of "known" files to confirm the reader is working as expected (i.e. we get known shas, symlinks, etc). Each reader can also have its own test for special formats.

To build the golden archives `podman` must be installed. From `testdata/lfs/test` run:
```
./build.sh
```
