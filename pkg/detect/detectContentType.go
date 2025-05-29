package detect

import (
	"io"
	"net/http"
)

func DetectContentType(file io.ReadSeeker) (string, error) {
	buf := make([]byte, 512) // Для детекции нужно первые 512 байт
	_, err := file.Read(buf)
	if err != nil {
		return "", err
	}
	_, err = file.Seek(0, io.SeekStart) // Возвращаем указатель в начало
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buf), nil
}
