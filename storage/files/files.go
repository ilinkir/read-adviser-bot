package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"read-adviser-bot/lib/myError"
	"read-adviser-bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultperm = 0774

var ErrNoSavedPages = errors.New("no saved Pages")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	const errMsg = "can`t save page"
	path := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(path, defaultperm); err != nil {
		return myError.Wrap(errMsg, err)
	}
	fName, err := fileName(page)
	if err != nil {
		return myError.Wrap(errMsg, err)
	}

	path = filepath.Join(path, fName)

	file, err := os.Create(path)
	if err != nil {
		return myError.Wrap(errMsg, err)
	}

	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return myError.Wrap(errMsg, err)
	}
	return nil
}

func (s Storage) Remove(page *storage.Page) error {
	const errMsg = "can`t remove page"

	fName, err := fileName(page)
	if err != nil {
		return myError.Wrap(errMsg, err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	if err := os.Remove(path); err != nil {
		return myError.Wrap(fmt.Sprintf(errMsg+"%s", path), err)
	}
	return nil
}

func (s Storage) IfExists(page *storage.Page) (bool, error) {
	const errMsg = "can`t check exist page"

	fName, err := fileName(page)
	if err != nil {
		return false, myError.Wrap(errMsg, err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, myError.Wrap(fmt.Sprintf(errMsg+"%s", path), err)
	}

	return true, nil

}

func (s Storage) pickRandom(userName string) (*storage.Page, error) {
	const errMsg = "can`t pick random file"
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, myError.Wrap(errMsg, err)
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, myError.Wrap("cant Decode Page", err)
	}

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, myError.Wrap("cant Decode Page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
