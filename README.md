
# Golang, MongoDB and Gin TODO API

  

### Overview

  

This API was developed to handle a simple todo app with simple functionality. It was created as a programming class final project which requiered a fullstack app to be developed. We were free to choose whatever technologies we wanted, but database actions was a must to get a good grade.

  

You are free to download and thereafter adjust the API to your needs. This project could be used as a boilerplate for an API that requires basic authentication (password hashing implemented) and basic CMS functionality.

  

&nbsp;

  

### Example Usage

  

Use as boilerplate for a CMS platform

Use as boilerplate for a user managment system

  

&nbsp;

  

### Getting Started

  

#### Required dependencies

  

Go programming language (Developed and tested on version: 1.18.1)

MongoDB local server or connection string for an external server (Has to have a database with two collections inside, "users" and "todos")

  

&nbsp;

  

#### Usage

  

**1. Clone the repository to your local machine**

  

Open Git Bash (if installed. Otherwise install) and type the following command in a directory that you are comfortable with using:

> `git clone https://github.com/ArseniSkobelev/golang-mongodb-todo-app-api .`

  

**2. Install required dependencies**

  

Enter the following command to install all required go dependencies:

> `go mod tidy`

  

**3. Adjusting environment variables in the code**

  

Open the directory in the text editor of choice and open the `main.go` file. Then adjust the following variables for your needs:

  

[![Variable screenshot](https://lh3.googleusercontent.com/pw/AM-JKLVsBnnPGDC0I25UGFdwNs6QZFt75YdsQURD3qa7r-O37jQksfrmZt85I103u85w7-yx8Bf8cwf7ngGiWSBbpsjn_yq4YdgOOW13JdoGNyTH_KXZZgsADoASHZC_iv_4QHpm8imE4jwjd9QNnJqBOeg=w491-h76-no "Variable screenshot")](https://lh3.googleusercontent.com/pw/AM-JKLVsBnnPGDC0I25UGFdwNs6QZFt75YdsQURD3qa7r-O37jQksfrmZt85I103u85w7-yx8Bf8cwf7ngGiWSBbpsjn_yq4YdgOOW13JdoGNyTH_KXZZgsADoASHZC_iv_4QHpm8imE4jwjd9QNnJqBOeg=w491-h76-no  "Variable screenshot")

  

_HASHING_LENGTH: Defines how many times a password should be hashed before it will be sent to the database.

_HOST: Defines where the server should listen for requests.

_DATABASE: (MOST IMPORTANT) Defines what database a connection will be executed for

  

**3.5 Important information**

  

If you want to use an external MongoDB database you will have to implement your own database connection. The reason for that is that this project is using a module which i had created earlier for prototyping. It involves a simple MongoDB connection without being able to change the connection URI

  

&nbsp;

**4. Running the code**

  

Enter the following command to run the code using and host the server at an earlier defined location:

> `go run .`

  

&nbsp;

  

#### Endpoints and their usage

  

**(localhost:3001)/createUser** - (POST, {"username": USERNAME, "email": EMAIL, "password": PASSWORD}) Gets data bundled with the POST request, hash the given password and then insert everything to the database in the 'users' collection. Returns - positive or negative message based on the outcome of the request.

  &nbsp;

**(localhost:3001)/createTodo** - (POST, {"title": TITLE, "status": STATUS, "owner": USERNAME}) Gets data bundled with the POST request and insert it into the 'todos' collection. Returns - positive or negative message based on the outcome of the request.

  &nbsp;

**(localhost:3001)/checkLogin** - (POST, {"username": USERNAME, "password": PASSWORD}) Gets data bundled with the POST request and finds a user with the provided username in 'users' collection. Thereafter check whether the password in the database matches the password provided in the POST requests body. Returns - boolean true or false.

  &nbsp;

**(localhost:3001)/getTodos** - (POST, {"owner": USERNAME}) Finds all todos currently active for a given user. Returns - (positive message and all todos in an array) or negative message based on the outcome of the request.

&nbsp;

**(localhost:3001)/changeTodoStatus** - (POST, {"Owner": USERNAME, "Title": TITLE, "NewStatus": NEWSTATUS}) Changes the given todo's status to the given new status variable. Deletes the todo afterwards to keep a clean db. Returns - positive or negative message based on the outcome of the request.



&nbsp;

  

### Technologies used

  

Go (used as the main programming language)

Gin (HTTP framework for Go)

MongoDB (used as the database provider)