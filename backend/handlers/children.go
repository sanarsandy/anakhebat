package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"tukem-backend/db"
	"tukem-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateChild(c echo.Context) error {
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	req := new(models.CreateChildRequest)
	if err := c.Bind(req); err != nil {
		c.Logger().Error("Bind error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request: " + err.Error()})
	}

	c.Logger().Info("Creating child with data: ", req)

	// Prepare gestational_age - use nil if not premature or if value is 0
	var gestationalAge *int
	if req.IsPremature && req.GestationalAge != nil && *req.GestationalAge > 0 {
		gestationalAge = req.GestationalAge
	}

	// Insert into DB
	query := `INSERT INTO children (parent_id, name, dob, gender, birth_weight, birth_height, is_premature, gestational_age) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	          RETURNING id, created_at`
	
	var child models.Child
	err := db.DB.QueryRow(query, userID, req.Name, req.DOB, req.Gender, req.BirthWeight, req.BirthHeight, req.IsPremature, gestationalAge).
		Scan(&child.ID, &child.CreatedAt)
	
	if err != nil {
		c.Logger().Error("Database error: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create child: " + err.Error()})
	}

	// Populate response
	child.ParentID = userID
	child.Name = req.Name
	child.DOB = req.DOB
	child.Gender = req.Gender
	child.BirthWeight = req.BirthWeight
	child.BirthHeight = req.BirthHeight
	child.IsPremature = req.IsPremature
	child.GestationalAge = gestationalAge

	c.Logger().Info("Child created successfully: ", child.ID)
	return c.JSON(http.StatusCreated, child)
}

func GetChildren(c echo.Context) error {
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	query := `SELECT id, parent_id, name, dob, gender, birth_weight, birth_height, is_premature, gestational_age, created_at 
	          FROM children WHERE parent_id = $1 ORDER BY created_at DESC`
	
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	children := []models.Child{}
	for rows.Next() {
		var child models.Child
		err := rows.Scan(&child.ID, &child.ParentID, &child.Name, &child.DOB, &child.Gender, 
			&child.BirthWeight, &child.BirthHeight, &child.IsPremature, &child.GestationalAge, &child.CreatedAt)
		if err != nil {
			continue
		}
		children = append(children, child)
	}

	return c.JSON(http.StatusOK, children)
}

func GetChild(c echo.Context) error {
	// Log to check if this handler is being called for denver-ii routes
	path := c.Request().URL.Path
	if strings.Contains(path, "denver-ii") {
		c.Logger().Errorf("GetChild handler called for denver-ii path: %s - THIS SHOULD NOT HAPPEN", path)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Route not found"})
	}
	childID := c.Param("id")
	
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	var child models.Child
	query := `SELECT id, parent_id, name, dob, gender, birth_weight, birth_height, is_premature, gestational_age, created_at 
	          FROM children WHERE id = $1 AND parent_id = $2`
	
	err := db.DB.QueryRow(query, childID, userID).Scan(
		&child.ID, &child.ParentID, &child.Name, &child.DOB, &child.Gender,
		&child.BirthWeight, &child.BirthHeight, &child.IsPremature, &child.GestationalAge, &child.CreatedAt)
	
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, child)
}

func UpdateChild(c echo.Context) error {
	childID := c.Param("id")
	
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	req := new(models.UpdateChildRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Prepare gestational_age - use nil if not premature or if value is 0
	var gestationalAge *int
	if req.IsPremature && req.GestationalAge != nil && *req.GestationalAge > 0 {
		gestationalAge = req.GestationalAge
	}

	query := `UPDATE children SET name = $1, dob = $2, gender = $3, birth_weight = $4, birth_height = $5, 
	          is_premature = $6, gestational_age = $7 
	          WHERE id = $8 AND parent_id = $9`
	
	result, err := db.DB.Exec(query, req.Name, req.DOB, req.Gender, req.BirthWeight, req.BirthHeight, 
		req.IsPremature, gestationalAge, childID, userID)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update child"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Child updated successfully"})
}

func DeleteChild(c echo.Context) error {
	childID := c.Param("id")
	
	// Get user ID from JWT
	user := c.Get("user").(*jwt.Token)
	claims := *user.Claims.(*jwt.MapClaims)
	userID := claims["user_id"].(string)

	query := `DELETE FROM children WHERE id = $1 AND parent_id = $2`
	result, err := db.DB.Exec(query, childID, userID)
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete child"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Child not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Child deleted successfully"})
}
