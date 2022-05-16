// telling the interpreter that this is an external package
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	dbConnection "github.com/ArseniSkobelev/go-mongodb-connection"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var _HASHING_LENGTH = 14
var _HOST = "localhost:3001"
var _DATABASE = "programming-finals-2022"

// structs
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
}

type Todo struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Status int                `bson:"status,omitempty"` // status "0" = Not active, status "1" = Active
	Owner  string             `bson:"owner,omitempty"`  // username that owns this todo
}

type LoginTest struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bsom:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

type TodoStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Owner     string             `bson:"username,omitempty"`
	NewStatus int                `bson:"newstatus,omitempty"`
	Title     string             `bson:"title,omitempty"`
}

type DeleteTodo struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Owner string             `bson:"username,omitempty"`
	Title string             `bson:"title,omitempty"`
}

type GetTodos struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Title  string             `bson:"title,omitempty"`
	Status int                `bson:"status,omitempty"` // status "0" = Not active, status "1" = Active
	Owner  string             `bson:"owner,omitempty"`  // username that owns this todo
}

type UserData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bsom:"username,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
}

// functions

// -- utility
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), _HASHING_LENGTH)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// -- endpoints

func greeting(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Welcome to the todo API. More information at: github.com/ArseniSkobelev/golang-mongodb-todo-app-api"})
}

func createUser(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var tempJson User
	json.Unmarshal([]byte(body), &tempJson)

	hashedPassword, _ := HashPassword(tempJson.Password)

	var tempUser = User{
		Username: strings.ToLower(tempJson.Username),
		Email:    strings.ToLower(tempJson.Email),
		Password: hashedPassword,
	}

	usersCollection := dbConnection.MongoConnection().Database(_DATABASE).Collection("users")

	result, err := usersCollection.InsertOne(context.TODO(), tempUser)
	if err != nil {
		fmt.Println(err)
	}

	var resultString = "Inserted document with _id: " + fmt.Sprintf("%v", result.InsertedID)
	c.IndentedJSON(http.StatusCreated, gin.H{"message": resultString})
}

func checkLogin(c *gin.Context) {
	usersCollection := dbConnection.MongoConnection().Database(_DATABASE).Collection("users")

	var tempJson LoginTest
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&tempJson)
	if err != nil {
		panic(err)
	}

	filter := &bson.M{
		"username": strings.ToLower(tempJson.Username),
	}

	var user User

	usersCollection.FindOne(c, filter).Decode(&user)

	userPassword := user.Password
	result := CheckPasswordHash(tempJson.Password, userPassword)
	if result == true {
		c.IndentedJSON(http.StatusOK, gin.H{"message": result})
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": result})
	}
}

func createTodo(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var tempJson Todo
	json.Unmarshal([]byte(body), &tempJson)

	var tempTodo = Todo{
		Title:  tempJson.Title,
		Status: tempJson.Status,
		Owner:  tempJson.Owner,
	}

	if tempTodo.Title == "" {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "❌ Missing a todo title!"})
	} else {
		TodosCollection := dbConnection.MongoConnection().Database(_DATABASE).Collection("todos")

		result, err := TodosCollection.InsertOne(context.TODO(), tempTodo)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(result)

		var resultString = "✅ Succesfully created a todo!"
		c.IndentedJSON(http.StatusCreated, gin.H{"message": resultString})
	}

}

func changeTodoStatus(c *gin.Context) {
	TodosCollection := dbConnection.MongoConnection().Database(_DATABASE).Collection("todos")

	body, _ := ioutil.ReadAll(c.Request.Body)
	var tempJson TodoStatus
	json.Unmarshal([]byte(body), &tempJson)

	selectedUser := bson.M{"owner": tempJson.Owner, "title": tempJson.Title}
	newStatus := tempJson.NewStatus
	result, err := TodosCollection.UpdateOne(c, selectedUser, bson.D{
		{"$set", bson.D{{"status", newStatus}}},
	},
	)
	if err != nil {
		log.Fatal(err)
	}

	selectedTodo := bson.M{"owner": tempJson.Owner, "title": tempJson.Title}
	deleteResult, err := TodosCollection.DeleteOne(context.TODO(), selectedTodo)
	if deleteResult.DeletedCount == 0 {
		log.Fatal("Error on deleting one todo", err)
	}

	if result.ModifiedCount > 0 {
		var resultString = "✅ Todo removed succesfully!"
		c.IndentedJSON(http.StatusOK, gin.H{"message": resultString})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "❌ Something went wrong!"})
	}
}

func deleteTodo(c *gin.Context) {
	TodosCollection := dbConnection.MongoConnection().Database(_DATABASE).Collection("todos")

	body, _ := ioutil.ReadAll(c.Request.Body)
	var tempJson DeleteTodo
	json.Unmarshal([]byte(body), &tempJson)

	selectedTodo := bson.M{"owner": tempJson.Owner, "title": tempJson.Title}
	deleteResult, err := TodosCollection.DeleteOne(context.TODO(), selectedTodo)
	if deleteResult.DeletedCount == 0 {
		log.Fatal("Error on deleting one todo", err)
	}

	if deleteResult.DeletedCount > 0 {
		var resultString = "Deleted documents: " + fmt.Sprintf("%v", deleteResult.DeletedCount)
		c.IndentedJSON(http.StatusOK, gin.H{"message": resultString})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
	}
}

func getTodos(c *gin.Context) {
	TodosCollection := dbConnection.MongoConnection().Database(_DATABASE).Collection("todos")

	body, _ := ioutil.ReadAll(c.Request.Body)
	var tempJson GetTodos
	json.Unmarshal([]byte(body), &tempJson)

	fmt.Println(tempJson.Owner)

	filter := bson.M{"owner": tempJson.Owner}
	cursor, err := TodosCollection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	var allTodos []GetTodos
	var tempTodos GetTodos
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(output, &tempTodos)
		allTodos = append(allTodos, tempTodos)
	}
	if len(allTodos) > 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "✅ Todos loaded succesfully!", "data": allTodos})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "❌ No todos found for the given user!"})
	}
}

func getUserData(c *gin.Context) {
	usersCollection := dbConnection.MongoConnection().Database(_DATABASE).Collection("users")

	body, _ := ioutil.ReadAll(c.Request.Body)
	var tempJson UserData
	json.Unmarshal([]byte(body), &tempJson)

	filter := bson.M{"username": tempJson.Username}
	cursor, err := usersCollection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []bson.M
	var allUsers []UserData
	var tempUsers UserData
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(output, &tempUsers)
		allUsers = append(allUsers, tempUsers)
	}
	if len(allUsers) > 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"message": allUsers})
	} else {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
	}
}

// -- main

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/", greeting)
	router.POST("/createUser", createUser)
	router.POST("/createTodo", createTodo)
	router.POST("/checkLogin", checkLogin)
	router.POST("/changeTodoStatus", changeTodoStatus)
	router.POST("/deleteTodo", deleteTodo)
	router.POST("/getTodos", getTodos)
	router.POST("/getUserData", getUserData)

	router.Run(_HOST)
}
