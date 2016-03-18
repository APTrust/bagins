# CHANGELOG

## 0.9.0

* Added support for custom tag files.

* Added support for tag manifests.

* Custom tag files may be anywhere outside the data directory, including custom directories.

* Custom tag files may be omitted from tagmanifests.

### New Methods and Breaking Changes

This version includes the following breaking changes:

* bagins.NewBag now takes a slice of hashNames, and will create manifests and optionally tag manifests for each named algorithm. The function signature was:

```go
NewBag(location string, name string, hashName string) (*Bag, error)
```

It is now:

```go
NewBag(location string, name string, hashNames []string, createTagManifests bool) (*Bag, error)
```

* bagins.ReadBag now automatically discovers payload manifests, and will support multiple payload manifests in the same bag. The function signature was:

```go
func ReadBag(pth string, tagfiles []string, manifest string) (*Bag, error)
```

It is now:

```go
func ReadBag(pathToFile string, tagfiles []string) (*Bag, error)
```

Note that the BagIt spec supports both parsed and unparsed tag files. The tagFiles param for ReadBag describes which tag files you want to parse when reading the bag. Other tag files in the bag will not be parsed.

* The Bag.AddTagfile() method has not changed, but developers should understand that it adds *managed* tag files in which you add name-value pairs to the TagFile.Data map (which is map[string]string). The bagins library writes these name-value pairs to the tag file internally, ensuring they conform to the correct tag file format.

* The new method Bag.AddCustomTagfile() adds unmanaged tag files to the bag. "Unmanaged" just means the bagins library makes no attempt to parse these custom tag files. You can use this method to add tag files of any type, text or binary, and bagins just copies them into the bag without question.

* New function Bag.UnparsedTagFiles() returns a list of tag files that bagins found but did not try to parse. When you call ReadBag with ["bag-info.txt", "my-info.txt"] as the tagFiles param, bagins will parse the two named tag files, and they will not be in the list of files returned by Bag.UnparsedTagFiles(). All other files in the bag that are not manifests or tagmanifests or part of the payload will be returned by Bag.UnparsedTagFiles().

* New constants bagins.PayloadManifest and bagins.TagManifest define the two types of manifests.

* New Bag fuction `GetManifest(manifestType, algorithm string) (*Manifest)` returns the manifest of the specified type and algorithm, if it exists.

* New Bag function `GetManifests(manifestType string) ([]*Manifest)` returns all manifests of the specified type, where type is either bagins.PayloadManifest or bagins.TagManifest.

* The signature for function NewManifest has changed from this:

```go
func NewManifest(pth string, hashName string) (*Manifest, error)
```

To this:

```go
func NewManifest(pathToFile string, hashName string, manifestType string) (*Manifest, error)
```

The new manifestType param should be either bagins.PayloadManifest or bagins.TagManifest.

* The new function Manifest.Type() returns either bagins.PayloadManifest or bagins.TagManifest.

* The new function Manifest.Algorithm() returns the name of the manifest's checksum algorithm, in all lower-case.

* The signature for function Payload.Add has changed from this:

```go
func (p *Payload) Add(srcPath string, dstPath string, m *Manifest) (string, error)
```

To this:

```go
func (p *Payload) Add(srcPath string, dstPath string, manifests []*Manifest) (map[string]string, error)
```

This allows the checksum of the newly added file to be written to multiple manifests.

* Similarly, the signature of Payload.AddAll has changed from this:

```go
func (p *Payload) AddAll(src string, m *Manifest) (fxs map[string]string, errs []error)
```

To this:

```go
func (p *Payload) AddAll(src string, manifests []*Manifest) (checksums map[string]map[string]string, errs []error)
```

See the inline documentation for an explanation of the return values.



## 0.8.0

* Added ability to open and read a bag directory on disk, tag files and manifests.

* Reduced number of concurrent files processed in checksums from 100 to 5

* Removed a number of unneeded methods.

* Added support for multiple tag fields with the same name and tag files respect
  field order.


# 0.7.0

* Added a Bag.Contents method that lists all the files found in the bag directory
  regardless to weather they are managed by the bag object or not.

* Added a Bag.FileManifest method to list all the files in a bag object it manages
  and can work on.

* Added a Bag.Invetory method that confirms that all files lin Bag.FileManifest are
  indeed written inside the bag.

* Added a Bag.Orphans method that lists any files in the bag that are not found
  in the Bag.FileManifest.


## 0.6.1

* bagmaker runs with throttled go routines to avoid a too many open files error.


## 0.6.0

* Can compile a command line executable basic bagger.  See README.rst for info


## 0.5.0

* Initial release.  Library works to build basic bags.
