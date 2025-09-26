package inner

import (
	"io"
	"io/fs"
)

type FS interface {
	SaveImage	(source io.Reader, newFileName string) 	(string, error)
	DeleteImage	(filename string) 						error
	Open		(name string)							(fs.File, error)

	Path 		()										string
}