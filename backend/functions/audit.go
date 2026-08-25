package functions

import (
	"database/sql"
	"log"
	"net/http"
	"product-catalogue/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// LogAudit inserts an audit log entry in the background (non-blocking).
func LogAudit(userID int, action, entityType string, entityID int, details string) {
	go func() {
		_, err := config.DB.Exec(
			"INSERT INTO audit_logs (user_id, action, entity_type, entity_id, details) VALUES (?, ?, ?, ?, ?)",
			userID, action, entityType, entityID, details,
		)
		if err != nil {
			log.Printf("[AUDIT ERROR] user=%d action=%s entity=%s/%d err=%v\n", userID, action, entityType, entityID, err)
		}
	}()
}

func GetAuditLogs(c *gin.Context) {
	role := c.GetString("role")
	if role != "system_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only system_admin can view audit logs"})
		c.Abort()
		return
	}

	search := strings.ToLower(strings.TrimSpace(c.Query("q")))
	searchPattern := "%" + search + "%"

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	rows, err := config.DB.QueryContext(ctx,
		`SELECT al.id, al.user_id, COALESCE(u.name,'') as username, al.action, al.entity_type, 
		 COALESCE(al.entity_id, 0), COALESCE(al.details,''), al.created_at 
		 FROM audit_logs al 
		 LEFT JOIN users u ON al.user_id = u.user_id 
		 WHERE ?='' OR lower(al.action) LIKE ? OR lower(al.entity_type) LIKE ? 
		 OR lower(COALESCE(al.details,'')) LIKE ? OR lower(COALESCE(u.name,'')) LIKE ?
		 ORDER BY al.created_at DESC LIMIT 200`,
		search, searchPattern, searchPattern, searchPattern, searchPattern,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch audit logs"})
		c.Abort()
		return
	}
	defer rows.Close()

	type AuditEntry struct {
		ID         int    `json:"id"`
		UserID     int    `json:"user_id"`
		Username   string `json:"username"`
		Action     string `json:"action"`
		EntityType string `json:"entity_type"`
		EntityID   int    `json:"entity_id"`
		Details    string `json:"details"`
		CreatedAt  string `json:"created_at"`
	}

	entries := []AuditEntry{}
	for rows.Next() {
		var e AuditEntry
		var createdAt sql.NullTime
		if err := rows.Scan(&e.ID, &e.UserID, &e.Username, &e.Action, &e.EntityType, &e.EntityID, &e.Details, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error scanning audit logs"})
			c.Abort()
			return
		}
		if createdAt.Valid {
			e.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		}
		entries = append(entries, e)
	}

	c.JSON(http.StatusOK, gin.H{"audit_logs": entries})
}
