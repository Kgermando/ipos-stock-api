package reservations

import (
	"strconv"

	"github.com/kgermando/ipos-stock-api/database"
	"github.com/kgermando/ipos-stock-api/models"
	"github.com/kgermando/ipos-stock-api/utils"

	"github.com/gofiber/fiber/v2"
)

// Synchronisation Send data to Local
func GetDataSynchronisation(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	sync_created := c.Query("sync_created", "2023-01-01")
	var data []models.Reservation

	if posUUID == "-" {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("created_at > ?", sync_created).
			Order("reservations.updated_at DESC").
			Preload("Pos").
			Find(&data)
	} else {
		db.Unscoped().Where("entreprise_uuid = ?", entrepriseUUID).
			Where("pos_uuid = ?", posUUID).
			Where("created_at > ?", sync_created).
			Order("reservations.updated_at DESC").
			Preload("Pos").
			Find(&data)
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All reservations sync data",
		"data":    data,
	})
}

// Paginate by entreprise
func GetPaginatedReservationEntreprise(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))
	search := c.Query("search", "")
	sort := c.Query("sort", "client_name")
	order := c.Query("order", "asc")

	offset := (page - 1) * limit

	var dataList []models.Reservation
	var totalRecords int64

	query := db.Model(&models.Reservation{}).Where("entreprise_uuid = ?", entrepriseUUID)

	if search != "" {
		query = query.Where("client_name LIKE ? OR table LIKE ? OR status LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	query.Count(&totalRecords)

	// Get paginated results
	err := query.Order(sort + " " + order).
		Offset(offset).
		Limit(limit).
		Preload("Pos").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch reservations",
			"error":   err.Error(),
		})
	}

	// Calculate total pages
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	// Prepare pagination metadata
	pagination := map[string]interface{}{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     limit,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All reservations",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Paginate by posUUID
func GetPaginatedReservationByPosUUID(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))
	search := c.Query("search", "")
	sort := c.Query("sort", "client_name")
	order := c.Query("order", "asc")

	offset := (page - 1) * limit

	var dataList []models.Reservation
	var totalRecords int64

	query := db.Model(&models.Reservation{}).Where("entreprise_uuid = ?", entrepriseUUID).Where("pos_uuid = ?", posUUID)

	if search != "" {
		query = query.Where("client_name LIKE ? OR table LIKE ? OR status LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	query.Count(&totalRecords)

	// Get paginated results
	err := query.Order(sort + " " + order).
		Offset(offset).
		Limit(limit).
		Preload("Pos").
		Find(&dataList).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch reservations",
			"error":   err.Error(),
		})
	}

	// Calculate total pages
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))

	// Prepare pagination metadata
	pagination := map[string]interface{}{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"current_page":  page,
		"page_size":     limit,
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "All reservations",
		"data":       dataList,
		"pagination": pagination,
	})
}

// Get All data
func GetAllReservations(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")

	var data []models.Reservation
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All reservations",
		"data":    data,
	})
}

// Get All data by search
func GetAllReservationBySearch(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	search := c.Query("search", "")

	var data []models.Reservation

	query := db.Where("entreprise_uuid = ?", entrepriseUUID).Where("pos_uuid = ?", posUUID)

	if search != "" {
		query = query.Where("client_name LIKE ? OR table LIKE ? OR status LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Preload("Pos").Find(&data)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All reservations",
		"data":    data,
	})
}

// Get one data
func GetReservation(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var reservation models.Reservation
	db.Where("uuid = ?", uuid).Preload("Pos").First(&reservation)
	if reservation.ClientName == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No reservation found",
				"data":    nil,
			},
		)
	}
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Reservation found",
			"data":    reservation,
		},
	)
}

// Create data
func CreateReservation(c *fiber.Ctx) error {
	p := &models.Reservation{}

	if err := c.BodyParser(&p); err != nil {
		return err
	}

	if p.ClientName == "" || p.Table == "" || p.ReservationDate == "" || p.ReservationTime == "" {
		return c.Status(400).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Form not complete",
				"data":    nil,
			},
		)
	}

	p.UUID = utils.GenerateUUID()
	p.Sync = true

	database.DB.Create(p)

	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Reservation created successfully",
			"data":    p,
		},
	)
}

// Update data
func UpdateReservation(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var reservation models.Reservation

	// Find the reservation
	db.Where("uuid = ?", uuid).First(&reservation)
	if reservation.ClientName == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No reservation found",
				"data":    nil,
			},
		)
	}

	// Parse request body
	if err := c.BodyParser(&reservation); err != nil {
		return err
	}

	reservation.Sync = true

	// Save to database
	database.DB.Save(&reservation)
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Reservation updated successfully",
			"data":    reservation,
		},
	)
}

// Delete data
func DeleteReservation(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	db := database.DB

	var reservation models.Reservation
	db.Where("uuid = ?", uuid).First(&reservation)
	if reservation.ClientName == "" {
		return c.Status(404).JSON(
			fiber.Map{
				"status":  "error",
				"message": "No reservation found",
				"data":    nil,
			},
		)
	}

	db.Delete(&reservation)
	return c.JSON(
		fiber.Map{
			"status":  "success",
			"message": "Reservation deleted successfully",
			"data":    nil,
		},
	)
}

// Get reservations by status
func GetReservationsByStatus(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	status := c.Params("status")

	var data []models.Reservation
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("status = ?", status).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All reservations by status",
		"data":    data,
	})
}

// Get reservations by date
func GetReservationsByDate(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	date := c.Params("date")

	var data []models.Reservation
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("reservation_date = ?", date).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All reservations by date",
		"data":    data,
	})
}

// Get reservations by table
func GetReservationsByTable(c *fiber.Ctx) error {
	db := database.DB
	entrepriseUUID := c.Params("entreprise_uuid")
	posUUID := c.Params("pos_uuid")
	table := c.Params("table")

	var data []models.Reservation
	db.Where("entreprise_uuid = ?", entrepriseUUID).
		Where("pos_uuid = ?", posUUID).
		Where("table = ?", table).
		Preload("Pos").
		Find(&data)
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All reservations by table",
		"data":    data,
	})
}
