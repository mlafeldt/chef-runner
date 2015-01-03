package omnibus

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _assets_install_wrapper_sh = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xa4\x92\x4d\x6a\xeb\x30\x10\xc7\xf7\x3a\xc5\x3c\x25\x8b\x97\x82\x6d\x92\x6d\x70\xa0\xb4\x9b\xae\x7a\x84\xa2\xc8\xa3\x4a\xa0\x8c\x8c\x34\x4e\x5b\xc8\xe1\xab\x24\x76\x48\x5c\x17\x02\xf5\x46\x86\xf9\x7f\xf8\xe7\xd1\xec\x5f\xb5\x75\x54\x25\x2b\x66\xf0\x08\x69\xa7\x22\xc3\x47\x54\x6d\x8b\x11\x54\x0c\x1d\x35\xf0\xba\x23\xb7\xed\x12\xbc\x50\x62\xe5\x3d\x46\x21\x12\x32\x14\x98\x4f\x1d\x5d\xcb\xf5\x7c\x29\xf6\x18\x93\x0b\x54\xcf\x57\x42\xec\x14\x39\x83\x89\xeb\x2a\xb4\x5c\x69\x8b\xa6\xea\xc7\xc5\x30\x2a\xf9\x93\x85\xee\x62\x44\xe2\xb7\x8b\xf7\xbf\x45\xd5\x40\x41\x4b\x90\xf3\x41\x29\x61\xb5\xa9\x1a\xdc\x57\xd4\x79\x0f\x07\xd0\x5d\xae\x6e\x24\x48\x28\xcc\x6a\x21\x84\x56\x09\xb3\xbc\xcf\x90\xe0\x48\x48\x79\x30\xca\x27\x5c\x08\xc8\x0f\x6a\x1b\x40\xd6\xf5\x06\x9e\x83\xa3\x77\xa0\xc0\x36\x9f\xa5\x3c\x4d\xd7\x6b\xe1\x15\xe7\x9e\xb1\xb8\xa7\x3d\x3a\xce\x02\xe8\x2b\x20\x18\x78\xca\x4c\x65\xd9\x47\x24\x9b\xfb\xcf\x7f\xe2\x92\xc9\xb1\xeb\xeb\x9d\x81\x93\xbb\xa0\xac\x1a\x11\xcb\x35\xb0\x45\x3a\xe9\x6e\xdb\x8f\x05\x97\xc2\xb1\x2d\x33\x9e\x37\xd1\x94\x93\x4c\x98\xd9\x27\x32\xef\x27\x9a\xa4\x32\x6e\x80\x7b\xb8\x25\xfb\x89\x05\xf5\xd5\x46\xee\x64\x1c\x5e\x94\x8f\xf9\x12\x7c\xfd\x9d\x71\x3a\xfe\x57\x46\x28\xf6\x57\x1f\x3d\x22\xc6\xa4\xb4\xf8\x0e\x00\x00\xff\xff\x83\x2a\x65\xdd\x2b\x03\x00\x00")

func assets_install_wrapper_sh_bytes() ([]byte, error) {
	return bindata_read(
		_assets_install_wrapper_sh,
		"assets/install-wrapper.sh",
	)
}

func assets_install_wrapper_sh() (*asset, error) {
	bytes, err := assets_install_wrapper_sh_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "assets/install-wrapper.sh", size: 811, mode: os.FileMode(420), modTime: time.Unix(1420311258, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"assets/install-wrapper.sh": assets_install_wrapper_sh,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"assets": &_bintree_t{nil, map[string]*_bintree_t{
		"install-wrapper.sh": &_bintree_t{assets_install_wrapper_sh, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

