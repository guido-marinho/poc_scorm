package scorm

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/guilherme-gatti/poc_scorm/internal/storage"
)

// ProcessScormPackage processa o pacote SCORM: descompacta, encontra o manifest, parseia e salva no banco.
func ProcessScormPackage(zipPath string) error {
	// Remove .zip pra criar a pasta destino
	dest := strings.TrimSuffix(zipPath, ".zip")

	// Descompacta o ZIP
	err := unzip(zipPath, dest)
	if err != nil {
		return fmt.Errorf("erro ao descompactar: %w", err)
	}

	// Busca recursiva pelo imsmanifest.xml
	var manifestPath string
	err = filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "imsmanifest.xml" {
			manifestPath = path
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("erro ao caminhar: %w", err)
	}

	if manifestPath == "" {
		return fmt.Errorf("imsmanifest.xml não encontrado em %s", dest)
	}

	// Abre o manifest
	manifest, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir imsmanifest.xml: %w", err)
	}
	defer manifest.Close()

	// Faz parse XML ➜ struct Manifest
	var data Manifest
	decoder := xml.NewDecoder(manifest)
	err = decoder.Decode(&data)
	if err != nil {
		return fmt.Errorf("erro ao parsear XML: %w", err)
	}

	fmt.Printf("Manifest: %+v\n", data)

	// Transforma struct ➜ JSON
	manifestJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("erro ao gerar JSON do manifest: %w", err)
	}

	// Insere no banco SQLite
	_, err = storage.DB.Exec(`
		INSERT INTO courses (identifier, version, manifest_json, path)
		VALUES (?, ?, ?, ?)
	`, data.Identifier, data.Version, manifestJSON, dest)

	if err != nil {
		return fmt.Errorf("erro ao salvar no banco: %w", err)
	}

	fmt.Println("✅ Manifest salvo no banco com sucesso!")

	return nil
}

// unzip extrai um arquivo .zip para a pasta destino.
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Detecta prefixo comum
	prefix := ""
	if len(r.File) > 0 {
		prefix = strings.SplitN(r.File[0].Name, "/", 2)[0] + "/"
	}

	for _, f := range r.File {
		path := filepath.Join(dest, strings.TrimPrefix(f.Name, prefix))

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
