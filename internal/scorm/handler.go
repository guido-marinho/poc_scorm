package scorm

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/jung-kurt/gofpdf"

	"github.com/gin-gonic/gin"
	"github.com/guilherme-gatti/poc_scorm/internal/storage"
)

// c *gin.Context - c é um ponteiro para o contexto da requisição que é imutável e isso é importante para evitar que a função modifique o contexto original e economizar memória. caso c nao fosse um ponteiro, a funcao teria que criar uma copia do contexto e isso seria mais custoso em termos de memoria e performance.
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	filePath := "./storage/" + file.Filename

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = ProcessScormPackage(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "uploaded and processed",
	})
}

// TrackHandler recebe tracking SCORM multi-SCO
func TrackHandler(c *gin.Context) {
	var payload struct {
		UserID  int    `json:"userId"`
		ScormID string `json:"scormId"`
		ScoID   string `json:"scoId"`
		Status  string `json:"status"`
		Score   int    `json:"score"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	var courseID int
	err := storage.DB.QueryRow(`SELECT id FROM courses WHERE identifier = ?`, payload.ScormID).Scan(&courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Curso não encontrado"})
		return
	}

	_, err = storage.DB.Exec(`
		INSERT INTO progress (user_id, course_id, sco_id, status, score)
		VALUES (?, ?, ?, ?, ?)
	`, payload.UserID, courseID, payload.ScoID, payload.Status, payload.Score)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar progresso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Progresso salvo",
		"scoId":  payload.ScoID, // útil pra debug
	})
}

// ProgressHandler lista progresso por userId
func ProgressHandler(c *gin.Context) {
	userID := c.Param("userId")

	rows, err := storage.DB.Query(`
		SELECT p.id, p.course_id, c.identifier, p.status, p.score, p.updated_at
		FROM progress p
		JOIN courses c ON p.course_id = c.id
		WHERE p.user_id = ?
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar progresso"})
		return
	}
	defer rows.Close()

	var result []gin.H
	for rows.Next() {
		var id, courseID, score int
		var identifier, status, updatedAt string

		if err := rows.Scan(&id, &courseID, &identifier, &status, &score, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar linhas"})
			return
		}

		result = append(result, gin.H{
			"id":        id,
			"course_id": courseID,
			"scorm_id":  identifier,
			"status":    status,
			"score":     score,
			"updatedAt": updatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}

// CSV
func ExportCSVHandler(c *gin.Context) {
	userID := c.Param("userId")

	rows, err := storage.DB.Query(`
		SELECT c.identifier, p.sco_id, p.status, p.score, p.updated_at
		FROM progress p JOIN courses c ON p.course_id = c.id
		WHERE p.user_id = ?`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar CSV"})
		return
	}
	defer rows.Close()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=progress.csv")

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	w.Write([]string{"Course", "SCO", "Status", "Score", "UpdatedAt"})
	for rows.Next() {
		var id, sco, status, updated string
		var score int
		rows.Scan(&id, &sco, &status, &score, &updated)
		w.Write([]string{id, sco, status, fmt.Sprint(score), updated})
	}
}

// PDF
func ExportPDFHandler(c *gin.Context) {
	userID := c.Param("userId")

	rows, err := storage.DB.Query(`
		SELECT c.identifier, p.sco_id, p.status, p.score, p.updated_at
		FROM progress p JOIN courses c ON p.course_id = c.id
		WHERE p.user_id = ?`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar PDF"})
		return
	}
	defer rows.Close()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Relatório de Progresso")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	for rows.Next() {
		var id, sco, status, updated string
		var score int
		rows.Scan(&id, &sco, &status, &score, &updated)
		line := fmt.Sprintf("%s - %s - %s - %d - %s", id, sco, status, score, updated)
		pdf.Cell(0, 10, line)
		pdf.Ln(8)
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment;filename=progress.pdf")
	err = pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao exportar PDF"})
	}
}
