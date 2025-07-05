package course

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gatti/poc_scorm/internal/storage"
)

func CoursesHandler(c *gin.Context) {
	rows, err := storage.DB.Query(`SELECT id, identifier, version, path FROM courses`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar cursos"})
		return
	}
	defer rows.Close()

	var result []gin.H

	for rows.Next() {
		var id int
		var identifier, version, path string

		if err := rows.Scan(&id, &identifier, &version, &path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler cursos"})
			return
		}

		result = append(result, gin.H{
			"id":         id,
			"identifier": identifier,
			"version":    version,
			"path":       path,
		})
	}

	c.JSON(http.StatusOK, result)
}

func DeleteCourseHandler(c *gin.Context) {
	courseID := c.Param("id")

	var path string
	err := storage.DB.QueryRow(`SELECT path FROM courses WHERE id = ?`, courseID).Scan(&path)
	if err != nil {
		c.JSON(404, gin.H{"error": "Curso não encontrado"})
		return
	}

	_, err = storage.DB.Exec(`DELETE FROM progress WHERE course_id = ?`, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover progressos"})
		return
	}

	_, err = storage.DB.Exec(`DELETE FROM courses WHERE id = ?`, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover curso"})
		return
	}

	err = os.RemoveAll(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover arquivos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Curso removido"})
}
