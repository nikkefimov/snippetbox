package nfs

import (
	"net/http"
	"path/filepath"
)

// create type NeuteredFileSystem which contains http.FileSystem
type NeuteredFileSystem struct {
	Fs http.FileSystem
}

// create method Open()
// we open request's path with method IsDir()
// check if the called path is a folder or not
// if its folder, use method Stat() for checking if the index.html file exists inside this folder
// if file is not exist, then method will return an os.ErrNotExist error
// which in turn will be converted via http.FileServer to 404 page not found response
func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.Fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil

}
