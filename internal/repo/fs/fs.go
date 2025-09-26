package fs

import (
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
)

type LocalFS struct {
	basePath string
}

func NewFS(shutdownChannel inner.ShutdownChannel, basePath string) *LocalFS {
	mode := 6400
	pathToImages := path.Join(basePath, "images")
	if _, err:= os.Stat(pathToImages); os.IsNotExist(err) {
		err := os.Mkdir(pathToImages, os.FileMode(mode))
		if err != nil {
			shutdownChannel.Send(inner.ShutdownMessage, "mkdirFiles", "failed to make direction `images`")
			return nil
		}
	}

	return &LocalFS{basePath: pathToImages}
}

func (fs *LocalFS) SaveImage(source io.Reader, newFileName string) (string, error) { // source as multipart.File, because using request.FormFile we parse file and its headers
	pathToImage := path.Join(fs.basePath, newFileName)
	f, err := os.OpenFile(pathToImage, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", err
	}

	defer f.Close()

	if _, err := io.Copy(f, source); err != nil {
		return "", err
	}

	return pathToImage, nil
}

func (fs *LocalFS) DeleteImage(filename string) error {
	pathToImage := path.Join(fs.basePath, filename)
	if _, err := os.Stat(pathToImage); os.IsNotExist(err) {
		return err
	}

	return os.Remove(pathToImage)
}

func (fs *LocalFS) Open(name string) (fs.File, error) {
	return os.Open(path.Join(fs.basePath, name))
}

func (fs *LocalFS) Path() string {
	return fs.basePath
}