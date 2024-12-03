package main

import (
	"Go-test/db"
	"encoding/json"
	"net/http"
)

// Todo represents a task in the database
type Todo struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Complete bool   `json:"complete"`
}

// Response structure for generic API responses
type Response struct {
	Message string `json:"message"`
}

func main() {

	// DSN: Replace with your MySQL credentials
	dsn := "root:P@ssw0rd@tcp(127.0.0.1:3306)/todo_go"

	// Initialize the database connection
	db.Init(dsn)

	// Run migrations to create the todo table
	db.DB.AutoMigrate(&Todo{})

	// Define API routes
	http.HandleFunc("/todos", getTodosHandler)
	http.HandleFunc("/add-todo", addTodoHandler)
    http.HandleFunc("/delete-todo/", deleteTodoHandler) // Add this line
    http.HandleFunc("/update-todo/", updateTodoHandler) // Add this line
	// Start the server
	http.ListenAndServe(":8080", nil)
}


// Handler to get all todos
func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todos []Todo
	result := db.DB.Find(&todos) // Fetch all todos from the database
	if result.Error != nil {
		http.Error(w, "Error fetching todos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos) // Return todos as JSON
}

// Handler to add a new todo
func addTodoHandler(w http.ResponseWriter, r *http.Request) {
    // Set response content type
    w.Header().Set("Content-Type", "application/json")

    // Ensure the request method is POST
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Decode the JSON payload into a Todo struct
    var newTodo Todo
    if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }

    // Save the new todo to the database
    result := db.DB.Create(&newTodo)
    if result.Error != nil {
        http.Error(w, "Error inserting todo into database", http.StatusInternalServerError)
        return
    }

    // Respond with a custom success message
    response := map[string]interface{}{
        "message": "Todo created successfully",
        "todo":    newTodo,
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}
// Handler to delete a todo by ID
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
    // Set response content type
    w.Header().Set("Content-Type", "application/json")

    // Ensure the request method is DELETE
    if r.Method != http.MethodDelete {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Get the ID from the URL path
    id := r.URL.Path[len("/delete-todo/"):]
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    // Try to delete the todo
    result := db.DB.Delete(&Todo{}, id)
    if result.Error != nil {
        http.Error(w, "Error deleting todo", http.StatusInternalServerError)
        return
    }

    // Check if any row was actually deleted
    if result.RowsAffected == 0 {
        http.Error(w, "Todo not found", http.StatusNotFound)
        return
    }

    // Respond with success message
    response := map[string]string{
        "message": "Todo deleted successfully",
    }
    json.NewEncoder(w).Encode(response)
}


// Handler to update a todo by ID
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
    // Set response content type
    w.Header().Set("Content-Type", "application/json")

    // Ensure the request method is PUT
    if r.Method != http.MethodPut {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Get the ID from the URL path
    id := r.URL.Path[len("/update-todo/"):]
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    // Fetch the todo by ID
    var todo Todo
    result := db.DB.First(&todo, id)
    if result.Error != nil {
        if result.RowsAffected == 0 {
            http.Error(w, "Todo not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error fetching todo", http.StatusInternalServerError)
        }
        return
    }

    // Decode the updated data from the request body
    var updatedData Todo
    if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }

    // Update the fields
    todo.Title = updatedData.Title
    todo.Complete = updatedData.Complete

    // Save the changes to the database
    saveResult := db.DB.Save(&todo)
    if saveResult.Error != nil {
        http.Error(w, "Error updating todo", http.StatusInternalServerError)
        return
    }

    // Respond with the updated todo
    response := map[string]interface{}{
        "message": "Todo updated successfully",
        "todo":    todo,
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
