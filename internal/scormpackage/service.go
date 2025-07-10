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

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/guilherme-gatti/poc_scorm/internal/storage"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ProcessScormPackage(zipPath string) error {
	dest := strings.TrimSuffix(zipPath, ".zip")

	err := unzip(zipPath, dest)
	if err != nil {
		return fmt.Errorf("erro ao descompactar: %w", err)
	}

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

	manifest, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir imsmanifest.xml: %w", err)
	}
	defer manifest.Close()

	var data Manifest
	decoder := xml.NewDecoder(manifest)
	err = decoder.Decode(&data)
	if err != nil {
		return fmt.Errorf("erro ao parsear XML: %w", err)
	}

	fmt.Printf("Manifest: %+v\n", data)

	digitalCourse, err := mapManifestToDigitalCourse(data)
	if err != nil {
		return fmt.Errorf("erro ao mapear manifest: %w", err)
	}

	err = validate.Struct(digitalCourse)
	if err != nil {
		return fmt.Errorf("erro na validação: %w", err)
	}

	fmt.Println("✅ Dados validados com sucesso!")

	manifestJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("erro ao gerar JSON do manifest: %w", err)
	}

	// Transforma DigitalCourse ➜ JSON (por enquanto não usa)
	// digitalCourseJSON, err := json.Marshal(digitalCourse)
	// if err != nil {
	// 	return fmt.Errorf("erro ao gerar JSON do curso digital: %w", err)
	// }

	// Insere no banco SQLite (sem digital_course_json por enquanto)
	_, err = storage.DB.Exec(`
		INSERT INTO courses (identifier, version, manifest_json, path)
		VALUES (?, ?, ?, ?)
	`, data.Identifier, data.Version, manifestJSON, dest)

	if err != nil {
		return fmt.Errorf("erro ao salvar no banco: %w", err)
	}

	fmt.Println("✅ Manifest e curso digital salvos no banco com sucesso!")

	return nil
}

func mapManifestToDigitalCourse(manifest Manifest) (*DigitalCourse, error) {
	title := manifest.Metadata.LOM.General.Title.Langstring
	if title == "" {
		title = manifest.Identifier
	}
	description := manifest.Metadata.LOM.General.Description.Langstring

	digitalCourse := &DigitalCourse{
		UUID:        uuid.New().String(),
		Name:        title,
		Description: description,
		CourseType:  "SCORM",
		Modules:     []Module{},
	}

	for _, org := range manifest.Organizations.Organization {
		module := Module{
			UUID:   uuid.New().String(),
			Name:   org.Title,
			Order:  0,
			Topics: []Topic{},
		}

		for i, item := range org.Items {
			topic := Topic{
				UUID:                  uuid.New().String(),
				Name:                  item.Title,
				Type:                  inferTopicType(item, manifest.Resources),
				Order:                 i,
				Description:           fmt.Sprintf("Tópico extraído do SCORM: %s", item.Title),
				DigitalCourseId:       digitalCourse.UUID,
				DigitalCourseModuleId: module.UUID,
			}

			if len(item.Items) > 0 {
				subTopics := processSubItems(item.Items, digitalCourse.UUID, module.UUID, len(module.Topics))
				module.Topics = append(module.Topics, subTopics...)
			} else {
				module.Topics = append(module.Topics, topic)
			}
		}

		digitalCourse.Modules = append(digitalCourse.Modules, module)
	}

	return digitalCourse, nil
}

func processSubItems(items []Item, courseUUID, moduleUUID string, startOrder int) []Topic {
	var topics []Topic

	for i, item := range items {
		topic := Topic{
			UUID:                  uuid.New().String(),
			Name:                  item.Title,
			Type:                  "LECTURE",
			Order:                 startOrder + i,
			Description:           fmt.Sprintf("Subtópico extraído do SCORM: %s", item.Title),
			DigitalCourseId:       courseUUID,
			DigitalCourseModuleId: moduleUUID,
		}

		topics = append(topics, topic)

		if len(item.Items) > 0 {
			subTopics := processSubItems(item.Items, courseUUID, moduleUUID, startOrder+len(topics))
			topics = append(topics, subTopics...)
		}
	}

	return topics
}

func inferTopicType(item Item, resources Resources) string {
	for _, resource := range resources.Resource {
		if resource.Identifier == item.IdentifierRef {
			if strings.Contains(strings.ToLower(resource.Href), "assessment") ||
				strings.Contains(strings.ToLower(resource.Href), "quiz") ||
				strings.Contains(strings.ToLower(resource.Href), "test") {
				return "ASSESSMENT"
			}
		}
	}
	return "LECTURE"
}

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

func ValidateDigitalCourse(digitalCourse *DigitalCourse) error {
	return validate.Struct(digitalCourse)
}
