package scorm

import (
	"encoding/xml"
	"time"
)

// Estruturas para parsing do imsmanifest.xml
type Manifest struct {
	XMLName       xml.Name      `xml:"manifest"`
	Identifier    string        `xml:"identifier,attr"`
	Version       string        `xml:"version,attr"`
	Metadata      Metadata      `xml:"metadata"`
	Organizations Organizations `xml:"organizations"`
	Resources     Resources     `xml:"resources"`
}

type Metadata struct {
	Schema        string `xml:"schema"`
	SchemaVersion string `xml:"schemaversion"`
	LOM           LOM    `xml:"lom"`
}

type LOM struct {
	General General `xml:"general"`
}

type General struct {
	Title       Title       `xml:"title"`
	Description Description `xml:"description"`
}

type Title struct {
	Langstring string `xml:"langstring"`
}

type Description struct {
	Langstring string `xml:"langstring"`
}

type Organizations struct {
	Default      string         `xml:"default,attr"`
	Organization []Organization `xml:"organization"`
}

type Organization struct {
	Identifier string `xml:"identifier,attr"`
	Title      string `xml:"title"`
	Items      []Item `xml:"item"`
}

type Item struct {
	Identifier    string `xml:"identifier,attr"`
	IdentifierRef string `xml:"identifierref,attr"`
	Title         string `xml:"title"`
	Items         []Item `xml:"item"`
}

type Resources struct {
	Resource []Resource `xml:"resource"`
}

type Resource struct {
	Identifier string `xml:"identifier,attr"`
	Type       string `xml:"type,attr"`
	Href       string `xml:"href,attr"`
	Files      []File `xml:"file"`
}

type File struct {
	Href string `xml:"href,attr"`
}

// Estruturas para dados processados
type ProcessedCourse struct {
	ID           int               `json:"id"`
	Identifier   string            `json:"identifier"`
	Version      string            `json:"version"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Modules      []ProcessedModule `json:"modules"`
	CreatedAt    time.Time         `json:"created_at"`
	ManifestJSON string            `json:"manifest_json"`
	Path         string            `json:"path"`
}

type ProcessedModule struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	UUID     string           `json:"uuid"`
	Order    int              `json:"order"`
	Topics   []ProcessedTopic `json:"topics"`
	CourseID int              `json:"course_id"`
}

type ProcessedTopic struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	UUID         string `json:"uuid"`
	Order        int    `json:"order"`
	Description  string `json:"description"`
	ModuleID     int    `json:"module_id"`
	ResourceHref string `json:"resource_href"`
}
