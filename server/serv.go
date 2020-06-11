package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"encoding/csv"
	"encoding/json"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/*
type Persona struct {
	Id       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Edad     int    `json:"edad"`
}
*/

/*
func RootEndpointGET(response http.ResponseWriter, request *http.Request) {
	persJson := `[{"id":1,"nombre": "Emmanuel","apellido": "German","edad": 22},{"id":2,"nombre":"Juan Diego","apellido":"Peña","edad":50}]`
	var pers []Persona
	json.Unmarshal([]byte(persJson), &pers)
	data, _ := json.Marshal(pers)
	for _, per := range pers {
		fmt.Println("Nombre: ", per.Nombre)
		fmt.Println("Apellido: ", per.Apellido)
	}
	//fmt.Println(pers)
	fmt.Println("JSON: ", string(data))
	response.Write(data)
}
func RootEndpointPOST(response http.ResponseWriter, request *http.Request) {
	persJson := `[{"id":1,"nombre": "Emmanuel","apellido": "German","edad": 22},{"id":2,"nombre":"Juan Diego","apellido":"Peña","edad":50}]`
	var pers []Persona
	json.Unmarshal([]byte(persJson), &pers)
	data, _ := json.Marshal(pers)
	for _, per := range pers {
		fmt.Println("Nombre: ", per.Nombre)
		fmt.Println("Apellido: ", per.Apellido)
	}
	//fmt.Println(pers)
	fmt.Println("JSON: ", string(data))
	response.Write(data)
}
*/

func readCsvFile(filePath string) ([]string, map[string]int, [][]float32) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	headers := make([]string, len(data[0]))
	copy(headers, data[0])

	col := make(map[string]int)
	for i, header := range headers {
		col[header] = i
	}

	data = data[1:]

	fdata := make([][]float32, len(data))
	for i := range fdata {
		fdata[i] = make([]float32, len(headers))
		for j := range fdata[i] {
			val, _ := strconv.ParseFloat(data[i][j], 32)
			fdata[i][j] = float32(val)
		}
	}

	return headers, col, fdata
}

func sliceCols(headers []string, col map[string]int, data [][]float32, newheaders []string) ([]string, map[string]int, [][]float32) {
	temp := make([]int, len(newheaders))
	newdata := make([][]float32, len(data))
	for i, newh := range newheaders {
		temp[i] = col[newh]
	}

	for i := range newdata {
		newdata[i] = make([]float32, len(temp))
		for j, t := range temp {
			newdata[i][j] = data[i][t]
		}
	}

	newcol := make(map[string]int)
	for i, header := range newheaders {
		newcol[header] = i
	}

	return newheaders, newcol, newdata

}

func head(data [][]float32, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(data[i])
	}
}

func normalize(data [][]float32) [][]float32 {
	mins, maxs := make([]float32, len(data[0])), make([]float32, len(data[0]))
	for i := 0; i < len(data[0]); i++ {
		mins[i], maxs[i] = float32(math.MaxFloat32), float32(0)
		for j := 0; j < len(data); j++ {
			if data[j][i] < mins[i] {
				mins[i] = data[j][i]
			}
			if data[j][i] > maxs[i] {
				maxs[i] = data[j][i]
			}
		}
		for j := 0; j < len(data); j++ {
			data[j][i] = ((data[j][i] - mins[i]) / (maxs[i] - mins[i]))
		}
	}

	return data
}

func monoKNN(xTrain [][]float32, yTrain [][]float32, xTest []float32, k int) float32 {
	distancias := make([][]float32, k)
	for i := 0; i < len(distancias); i++ {
		distancias[i] = []float32{math.MaxFloat32, -1}
	}

	for i := 0; i < len(xTrain); i++ {
		suma := float32(0)
		for j := 0; j < len(xTrain[i]); j++ {
			suma += (xTrain[i][j] - xTest[j]) * (xTrain[i][j] - xTest[j])
		}
		suma = float32(math.Sqrt(float64(suma)))
		j := len(distancias) - 1
		for ; j >= 0; j-- {
			if distancias[j][0] <= suma {
				temp := append(distancias[:j+1], []float32{suma, yTrain[i][0]})
				//fmt.Println(reflect.TypeOf(temp))
				for _, tmp := range distancias[j+1:] {
					temp = append(temp, tmp)
				}
				distancias = temp[:k]
				break
			} else {
				if j == 0 {
					temp := make([][]float32, k+1)
					temp[0] = []float32{suma, yTrain[i][0]}
					for m := range distancias {
						temp[m+1] = distancias[m]
					}
					distancias = temp[:k]
				}
			}
		}
	}
	clases := make(map[float32]int)
	for _, dist := range distancias {
		if _, found := clases[dist[1]]; !found {
			clases[dist[1]] = 1
		} else {
			clases[dist[1]]++
		}
	}
	var res float32
	max := -1
	for key, val := range clases {
		fmt.Printf("Clase '%d': %d ocurrencias\n", int(key), val)
		if val > max {
			max = val
			res = key
		}
	}

	return res
}

// GLOBAL VARIABLES

var colnames []string
var col map[string]int
var data [][]float32
var xTrain [][]float32
var yTrain [][]float32

type person struct {
	Age              float32 `json:"age"`
	Height           float32 `json:"height"`
	Weight           float32 `json:"weight"`
	Gender           float32 `json:"gender"`
	Sbp              float32 `json:"sbp"`
	Dbp              float32 `json:"dbp"`
	Cholesterol      float32 `json:"cholesterol"`
	Glucose          float32 `json:"glucose"`
	Smoking          float32 `json:"smoking"`
	AlcoholConsume   float32 `json:"alcohol_consume"`
	PhysicalActivity float32 `json:"physical_activity"`
}
type body struct {
	Person person
}

func KNNrequest(response http.ResponseWriter, request *http.Request) {

	var bdy body
	err := json.NewDecoder(request.Body).Decode(&bdy)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR CTM")
		return
	}
	per := bdy.Person

	xTest := []float32{per.Age, per.Height, per.Weight, per.Gender, per.Sbp,
		per.Dbp, per.Cholesterol, per.Glucose, per.Smoking, per.AlcoholConsume, per.PhysicalActivity}

	xTrain = append(xTrain, xTest)
	xTrain = normalize(xTrain)
	xTrain, xTest = xTrain[:len(yTrain)-1], xTrain[len(yTrain)-1]
	fmt.Fprintf(response, "%.0f", monoKNN(xTrain, yTrain, xTest, 200))
}

func main() {

	colnames, col, data = readCsvFile("../cardio_train.csv")
	_, _, xTrain = sliceCols(colnames, col, data, colnames[1:len(colnames)-1])
	_, _, yTrain = sliceCols(colnames, col, data, []string{colnames[len(colnames)-1]})

	//Datos de prueba
	//xTest := []float32{22, 2, 178, 60, 110, 65, 1, 1, 0, 0, 0}

	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	//router.HandleFunc("/", RootEndpointGET).Methods("GET")
	//router.HandleFunc("/", RootEndpointPOST).Methods("POST")
	router.HandleFunc("/knn", KNNrequest).Methods("POST")

	fmt.Println("Now server is running on port 8000")
	http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(router))
}
