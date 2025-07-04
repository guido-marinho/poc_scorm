package scorm

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ProcessScormPackage processa o pacote SCORM: descompacta, encontra o manifest e faz parse.
func ProcessScormPackage(zipPath string) error {
	// Remove o .zip para criar a pasta destino
	dest := strings.TrimSuffix(zipPath, ".zip")

	// Descompacta
	err := unzip(zipPath, dest)
	if err != nil {
		return fmt.Errorf("erro ao descompactar o arquivo: %w", err)
	}

	// Busca recursiva pelo imsmanifest.xml dentro da pasta destino
	var manifestPath string
	err = filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // para o Walk se der problema
		}

		if !info.IsDir() && info.Name() == "imsmanifest.xml" {
			manifestPath = path
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("erro ao caminhar pelos arquivos: %w", err)
	}

	if manifestPath == "" {
		return fmt.Errorf("imsmanifest.xml não encontrado em %s", dest)
	}

	// Abre o arquivo encontrado
	manifest, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir imsmanifest.xml: %w", err)
	}
	defer manifest.Close()

	// Faz o parse
	var data Manifest
	decoder := xml.NewDecoder(manifest)
	err = decoder.Decode(&data)
	if err != nil {
		return fmt.Errorf("erro ao fazer parse do XML: %w", err)
	}

	fmt.Printf("Manifest: %+v\n", data)

	return nil
}

// unzip extrai um arquivo .zip para a pasta destino.
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(dstFile, fileInArchive)

		dstFile.Close()
		fileInArchive.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// Manifest representa a estrutura básica do imsmanifest.xml
type Manifest struct {
	XMLName    xml.Name `xml:"manifest"`
	Identifier string   `xml:"identifier,attr"`
	Version    string   `xml:"version,attr"`
}
