package student

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/amit8889/golangCRUDApi/internal/storage"
	"github.com/amit8889/golangCRUDApi/internal/types"
	"github.com/amit8889/golangCRUDApi/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")

		var stu types.Student

		err := json.NewDecoder(r.Body).Decode(&stu)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
				"message": fmt.Sprintf("Invalid request body: %v", err.Error()),
				"success": false,
			})
			return

		}
		slog.Info("parsed a student")
		fmt.Println(stu)
		validationErrors := response.ValidateStruct(stu)
		if validationErrors != nil {
			response.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
				"message": "validation failed",
				"errors":  validationErrors,
				"success": false,
			})
			return
		}
		//fmt.Println(stu)
		// create a new student
		student, err := storage.CreateStudent(stu.Name, stu.Email, stu.Age)
		if err != nil {
			slog.Info("Error in creating student :", err)
			response.WriteJson(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to create student",
				"success": false,
			})
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]interface{}{
			"message": "student created successfully",
			"success": true,
			"id":      student,
		})

	}
}

func GetStudentByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting a student by id")
		//vars := mux.Vars(r) // Get the path variables
		//id := vars["id"]
		id := r.PathValue("id") // Get the path variables
		slog.Info("Id of user is : ", slog.String("id", id))
		if id == "" {
			response.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
				"message": "id is required",
				"success": false,
			})
			return
		}
		student, err := storage.GetStudent(id)
		if err != nil {
			slog.Info("Error in getting student by id :", err)
			response.WriteJson(w, http.StatusNotFound, map[string]interface{}{
				"message": "Failed to get student",
				"success": false,
			})
			return
		}

		//stu, err := storage.GetStudent(id)
		response.WriteJson(w, http.StatusAccepted, map[string]interface{}{
			"message": "student details",
			"success": true,
			"student": student,
		})

	}
}
