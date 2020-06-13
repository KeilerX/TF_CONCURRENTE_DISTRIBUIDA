package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

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

////////////           MULTI HILO
func monoKNN(xTrain [][]float32, yTrain [][]float32, xTest []float32, k int) (float32, map[float32]int) {
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
				temp := make([][]float32, k+1)
				for m := 0; m < j+1; m++ {
					temp[m] = make([]float32, 2)
					copy(temp[m], distancias[m])
				}
				temp[j+1] = []float32{suma, yTrain[i][0]}
				for m := j + 1; m < k; m++ {
					temp[m+1] = make([]float32, 2)
					copy(temp[m+1], distancias[m])
				}
				distancias = temp[:k]
				break
			} else {
				if j == 0 {
					temp := make([][]float32, k+1)
					temp[0] = []float32{suma, yTrain[i][0]}
					for m := range distancias {
						temp[m+1] = make([]float32, 2)
						copy(temp[m+1], distancias[m])
					}
					distancias = temp[:k]
				}
			}
		}
	}
	clases := make(map[float32]int)
	for _, dist := range distancias {
		if dist[1] == -1 {
			continue
		}
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

	return res, clases
}

////////////           MULTI HILO
func routineKNN(xTrain [][]float32, yTrain [][]float32, xTest []float32, k int, outCh chan []float32) {
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
				temp := make([][]float32, k+1)
				for m := 0; m < j+1; m++ {
					temp[m] = make([]float32, 2)
					copy(temp[m], distancias[m])
				}
				temp[j+1] = []float32{suma, yTrain[i][0]}
				for m := j + 1; m < k; m++ {
					temp[m+1] = make([]float32, 2)
					copy(temp[m+1], distancias[m])
				}
				distancias = temp[:k]
				break
			} else {
				if j == 0 {
					temp := make([][]float32, k+1)
					temp[0] = []float32{suma, yTrain[i][0]}
					for m := range distancias {
						temp[m+1] = make([]float32, 2)
						copy(temp[m+1], distancias[m])
					}
					distancias = temp[:k]
				}
			}
		}
	}

	for _, dist := range distancias {
		outCh <- dist
	}
}
func multiKNN(xTrain [][]float32, yTrain [][]float32, xTest []float32, k int, routines int) (float32, map[float32]int) {

	outCh := make(chan []float32)
	xsize := len(xTrain) / routines
	for i := 0; i < routines; i++ {
		if i < routines-1 {
			go routineKNN(xTrain[i*xsize:(i+1)*xsize], yTrain[i*xsize:(i+1)*xsize], xTest, k, outCh)
		} else {
			go routineKNN(xTrain[i*xsize:], yTrain[i*xsize:], xTest, k, outCh)
		}
	}

	distancias := make([][]float32, k)
	for i := 0; i < len(distancias); i++ {
		distancias[i] = []float32{math.MaxFloat32, -1}
	}

	for i := 0; i < k*routines; i++ {
		j := len(distancias) - 1
		candidate := <-outCh
		for ; j >= 0; j-- {
			if distancias[j][0] <= candidate[0] {
				temp := make([][]float32, k+1)
				for m := 0; m < j+1; m++ {
					temp[m] = make([]float32, 2)
					copy(temp[m], distancias[m])
				}
				temp[j+1] = []float32{candidate[0], candidate[1]}
				for m := j + 1; m < k; m++ {
					temp[m+1] = make([]float32, 2)
					copy(temp[m+1], distancias[m])
				}
				distancias = temp[:k]
				break
			} else {
				if j == 0 {
					temp := make([][]float32, k+1)
					temp[0] = []float32{candidate[0], candidate[1]}
					for m := range distancias {
						temp[m+1] = make([]float32, 2)
						copy(temp[m+1], distancias[m])
					}
					distancias = temp[:k]
				}
			}
		}
	}
	close(outCh)

	clases := make(map[float32]int)
	for _, dist := range distancias {
		if dist[1] == -1 {
			continue
		}
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

	return res, clases
}

func centroidesCercanos(id int, dftemp [][]float32, centers [][]float32, GCh chan []int, idCh chan int) {
	n := len(dftemp)
	ncols := len(dftemp[0])
	k := len(centers)
	dist := make([][]float32, k)
	for i := 0; i < k; i++ {
		dist[i] = make([]float32, n)
	}
	G := make([]int, n)
	for i, point := range dftemp {
		c := 0
		for cent := 0; cent < k; cent++ {
			suma := float32(0)
			for col := 0; col < ncols; col++ {
				suma += (point[col] - centers[cent][col]) * (point[col] - centers[cent][col])
			}
			dist[cent][i] = float32(math.Sqrt(float64(suma)))
			if dist[cent][i] < dist[c][i] {
				c = cent
			}
		}
		G[i] = c
	}
	GCh <- G
	idCh <- id
}

func multiKMeans(dftemp [][]float32, k int, maxIt int) ([]int, [][]float32, int) {
	n := len(dftemp)
	ncols := len(dftemp[0])
	centers := make([][]float32, k)
	G := make([]int, n)
	rc := make(map[int]struct{})
	it := 0
	for len(rc) != k {
		r := rand.Intn(n - 1)
		rc[r] = struct{}{}
	}
	for i := range rc {
		centers[i] = make([]float32, ncols)
		for j := 0; j < ncols; j++ {
			centers[i][j] = dftemp[i][j]
		}
	}

	for repeat := false; !repeat && it < maxIt; {
		GCh := make(chan []int, 1)
		idCh := make(chan int, 1)
		tsize := n / k
		for i := 0; i < k; i++ {
			if i < k-1 {
				go centroidesCercanos(i, dftemp[i*tsize:i*tsize+tsize], centers, GCh, idCh)
			} else {
				go centroidesCercanos(i, dftemp[i*tsize:], centers, GCh, idCh)
			}
		}
		for i := 0; i < k; i++ {
			Gpiece := <-GCh
			id := <-idCh
			for j := 0; j < len(Gpiece); j++ {
				G[id*tsize+j] = Gpiece[j]
			}
		}
		close(GCh)
		end := make(chan bool)
		newcenters := make([][]float32, k)
		for i := 0; i < k; i++ {
			newcenters[i] = make([]float32, ncols)
		}
		counters := make([]int, k)
		for i := 0; i < k; i++ {
			counters[i] = 0
		}
		for i := 0; i < k; i++ {
			go func(id int) {
				for i := id; i < n; i += k {
					for j := 0; j < ncols; j++ {
						newcenters[G[i]][j] += dftemp[i][j]
					}
					counters[G[i]]++
				}
				end <- true
			}(i)
		}
		for i := 0; i < k; i++ {
			<-end
		}
		for i := 0; i < k; i++ {
			for j := 0; j < ncols; j++ {
				newcenters[i][j] /= float32(counters[i])
			}
		}
		it++
		repeat = true
		for i := 0; i < k; i++ {
			for j := 0; j < ncols; j++ {
				if centers[i][j] != newcenters[i][j] {
					repeat = false
					break
				}
			}

			copy(centers[i], newcenters[i])
		}
	}

	return G, centers, it
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
	Person    person `json:"person"`
	Algorithm int    `json:"algorithm"`
	K         int    `json:"k"`
	Threads   int    `json:"n_threads"`
}

type body2 struct {
	K int `json:"k"`
}

type res struct {
	Clase  int `json:"clase"`
	Ocurs0 int `json:"ocurs0"`
	Ocurs1 int `json:"ocurs1"`
}

type res2 struct {
	Centroids []person `json:"centroids"`
	Ncentroid []int    `json:"ncentroid"`
}

func knnRequest(r http.ResponseWriter, request *http.Request) {
	var bdy body
	err := json.NewDecoder(request.Body).Decode(&bdy)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	per := bdy.Person

	xTest := []float32{per.Age, per.Height, per.Weight, per.Gender, per.Sbp,
		per.Dbp, per.Cholesterol, per.Glucose, per.Smoking, per.AlcoholConsume, per.PhysicalActivity}

	dftemp := make([][]float32, len(xTrain))
	for i := range xTrain {
		dftemp[i] = make([]float32, len(xTrain[i]))
		copy(dftemp[i], xTrain[i])
	}
	dftemp = append(dftemp, xTest)
	dftemp = normalize(dftemp)
	dftemp, xTest = dftemp[:len(yTrain)-1], dftemp[len(yTrain)-1]
	if bdy.Algorithm == 1 {
		//fmt.Fprintf(response, "%.0f", monoKNN(dftemp, yTrain, xTest, bdy.K))
		clase, ocurs := monoKNN(dftemp, yTrain, xTest, bdy.K)
		response := res{Clase: int(clase), Ocurs0: ocurs[0], Ocurs1: ocurs[1]}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(r, err.Error(), http.StatusBadRequest)
			fmt.Println("ERROR GAA")
			return
		}
		fmt.Fprintf(r, "%s", jsonResponse)
	} else if bdy.Algorithm == 2 {
		//fmt.Fprintf(response, "%.0f", monoKNN(dftemp, yTrain, xTest, bdy.K))
		clase, ocurs := multiKNN(dftemp, yTrain, xTest, bdy.K, bdy.Threads)
		response := res{Clase: int(clase), Ocurs0: ocurs[0], Ocurs1: ocurs[1]}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(r, err.Error(), http.StatusBadRequest)
			fmt.Println("ERROR GAA")
			return
		}
		fmt.Fprintf(r, "%s", jsonResponse)
	}
}

func kmeansRequest(r http.ResponseWriter, request *http.Request) {
	var bdy body2
	err := json.NewDecoder(request.Body).Decode(&bdy)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	k := bdy.K

	dftemp := make([][]float32, len(xTrain))
	for i := range xTrain {
		dftemp[i] = make([]float32, len(xTrain[i]))
		copy(dftemp[i], xTrain[i])
	}
	maxIt := 100
	G, centers, _ := multiKMeans(dftemp, k, maxIt)
	ocurs := make(map[int]int)
	ncentroid := make([]int, k)
	centroids := make([]person, k)
	for _, clase := range G {
		if _, found := ocurs[clase]; !found {
			ocurs[clase] = 1
		} else {
			ocurs[clase]++
		}
	}
	for key, val := range ocurs {
		ncentroid[key] = val
	}

	for i := 0; i < k; i++ {
		centroids[i] = person{Age: centers[i][0], Height: centers[i][1],
			Weight: centers[i][2], Gender: centers[i][3], Sbp: centers[i][4],
			Dbp: centers[i][5], Cholesterol: centers[i][6], Glucose: centers[i][7],
			Smoking: centers[i][8], AlcoholConsume: centers[i][9], PhysicalActivity: centers[i][10]}
	}
	response := res2{Centroids: centroids, Ncentroid: ncentroid}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	fmt.Fprintf(r, "%s", jsonResponse)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	colnames, col, data = readCsvFile("cardio_train.csv")
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
	router.HandleFunc("/knn", knnRequest).Methods("POST")
	router.HandleFunc("/kmeans", kmeansRequest).Methods("POST")

	fmt.Println("Now server is running on port 8000")
	http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(router))
}
