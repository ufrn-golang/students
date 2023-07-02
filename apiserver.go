package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Student struct {
	Name string		`json:"name"`
	ID string		`json:"id"`
	Year uint		`json:"year"`
	Email string	`json:"email"`
	Course string	`json:"course-name"`
	IsActive bool	`json:"active"`
}

var students []Student
func init() {
	students = []Student{
		{"Paul Terry", "2021070323", 2021, "paul.terry@university.com", "Information Technology", true},
		{"Amanda Ross", "2020212548", 2020, "amanda.ross@university.com", "Computer Science", true},
		{"John Kent", "2022232227", 2022, "john.kent@university.com", "Information Technology", false},
		{"Diana Carter", "2020212333", 2020, "diana.carter@university.com", "Computer Science", true},
		{"Carla Becker", "2022236452", 2022, "carla.becker", "Information Technology", true},
		{"Antony Scott", "2023651143", 2023, "antony.scott", "Software Engineering", true},
	}
}

func getStudentById(res http.ResponseWriter, req *http.Request) {
	urlPathElements := strings.Split(req.URL.Path, "/")
	id := urlPathElements[2]
	student := Student{}
	for _, s := range students {
		if s.ID == id {
			student = s
		}
	}

	if student != (Student{}) {
		if resBody, err := json.Marshal(student); err != nil {
			log.Fatal(err)
		} else {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusOK)
			res.Write(resBody)
		}
	} else {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "404 Not Found")
	}
}

func getStudentsByCourse(res http.ResponseWriter, req *http.Request) {
	courseName := req.URL.Query().Get("course-name")

	var courseStudents []Student
	for _, student := range students {
		if student.Course == courseName {
			courseStudents = append(courseStudents, student)
		}
	}

	if resBody, err := json.Marshal(courseStudents); err != nil {
		log.Fatal(err)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(resBody)
	}
}

func getStudents(res http.ResponseWriter, req *http.Request) {
	if resBody, err := json.Marshal(students); err != nil {
		log.Fatal(err)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(resBody)
	}
}

func addStudent(res http.ResponseWriter, req *http.Request) {
	contents, _ := io.ReadAll(req.Body)
    defer req.Body.Close()

    var s Student
    err := json.Unmarshal(contents, &s)
    if err != nil {
        log.Fatal(err)
    } else {
		students = append(students, s)
        res.WriteHeader(http.StatusCreated)
    }
}

func routing(res http.ResponseWriter, req *http.Request) {
	urlPathElements := strings.Split(req.URL.Path, "/")
	if urlPathElements[1] == "students" {
		queryParams := req.URL.Query()
		if len(queryParams) == 0 {
			switch urlPathElements[2] {
			case "": 
				getStudents(res, req)				// /students
			case "add" : 
				addStudent(res, req)				// /students/add
			default:
				getStudentById(res, req)			// /students/{id}
			}
		} else {									// students/?course-name={courseName}
			getStudentsByCourse(res, req)
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "400 Bad Request")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", routing)

	server := http.Server{
		Addr: "localhost:8081", 
		Handler: mux, 
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
