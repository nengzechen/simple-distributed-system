package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers() {
	handler := new(studentsHandler)
	http.Handle("/students", handler)
	http.Handle("/students/", handler)
}

type studentsHandler struct{}

func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/") // 分割的返回值是分出多少块，一个/分出两个块
	switch len(pathSegments) {
	case 2:
		sh.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func (sh studentsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	data, err := sh.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}
func (sh studentsHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	data, err := sh.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to serialize student: %q", err)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}

func (sh studentsHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	var g Grade
	dec := json.NewDecoder(r.Body) // 解码body中的内容
	err = dec.Decode(&g)           // 放到g中
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	fmt.Println("芝士雪豹")
	student.Grades = append(student.Grades, g)
	fmt.Print(student.Grades)
	w.WriteHeader(http.StatusCreated) // 201状态
	data, err := sh.toJSON(g)
	fmt.Println(data)
	fmt.Println(g)
	if err != nil {
		log.Println(err)
		return // 如果append成功，toJSON失败，后面都就不该执行了
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data) // 数据写回
}

func (sh studentsHandler) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b) // 建立编码器
	err := enc.Encode(obj)     // 对传进来的变量进行编码
	if err != nil {
		return nil, fmt.Errorf("failed to serialize students %q", err)
	}
	return b.Bytes(), nil
}
