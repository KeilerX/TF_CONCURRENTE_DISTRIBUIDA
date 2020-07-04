package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gonum/stat"

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

func standarization(data [][]float32) ([][]float32, []float64, []float64) {
	newdata := make([][]float32, len(data))
	for i := 0; i < len(data); i++ {
		newdata[i] = make([]float32, len(data[i]))
	}
	mean := make([]float64, len(data[0]))
	std := make([]float64, len(data[0]))

	for i := 0; i < len(data[0]); i++ {
		column := make([]float64, len(data))
		mean[i] = float64(0)
		for j := 0; j < len(data); j++ {
			column[j] = float64(data[j][i])
			mean[i] += float64(data[j][i])
		}
		mean[i] = mean[i] / float64(len(data))
		std[i] = stat.StdDev(column, nil)
		for j := 0; j < len(data); j++ {
			newdata[j][i] = (data[j][i] - float32(mean[i])) / float32(std[i])
		}
	}

	return newdata, mean, std
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
	for i := 0; i < k; i++ {
		centers[i] = make([]float32, ncols)
	}
	G := make([]int, n)
	rc := make(map[int]struct{})
	it := 0
	for len(rc) != k {
		r := rand.Intn(n - 1)
		rc[r] = struct{}{}
	}
	temp := 0
	for i := range rc {
		for j := 0; j < ncols; j++ {
			centers[temp][j] = dftemp[i][j]
		}
		temp++
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
var colnamesCovid []string
var colnamesDiseases []string
var col map[string]int
var colCovid map[string]int
var colDiseases map[string]int

var data [][]float32
var dfCovid [][]float32
var dfDiseases [][]float32

var xTrain [][]float32
var xTrainCovid [][]float32

var yTrain [][]float32
var yTrainCovid [][]float32

var groupsSelected []pacient

type BlockChain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

type tmsg struct {
	Code string
	Addr string
	Blk  Block
	Bc   BlockChain
}

////////////////           KNN

type person struct {
	Age              float32 `json:"age"`
	Gender           float32 `json:"gender"`
	Height           float32 `json:"height"`
	Weight           float32 `json:"weight"`
	Sbp              float32 `json:"sbp"`
	Dbp              float32 `json:"dbp"`
	Cholesterol      float32 `json:"cholesterol"`
	Glucose          float32 `json:"glucose"`
	Smoking          float32 `json:"smoking"`
	AlcoholConsume   float32 `json:"alcohol_consume"`
	PhysicalActivity float32 `json:"physical_activity"`
}

type testInput struct {
	Edad            float32 `json:"edad"`
	Genero          float32 `json:"genero"`
	Tos             float32 `json:"tos"`
	Temperatura     float32 `json:"temperatura"`
	DolorGarganta   float32 `json:"dolor_garganta"`
	MalestarGeneral float32 `json:"malestar_general"`
}

type pacient struct {
	Edad          float32 `json:"edad"`
	Genero        float32 `json:"genero"`
	CardioDisease float32 `json:"cardio_disease"`
	Diabetes      float32 `json:"diabetes"`
	RespDisease   float32 `json:"resp_disease"`
	Hipertension  float32 `json:"hipertension"`
	Cancer        float32 `json:"cancer"`
}
type body struct {
	Person    person `json:"person"`
	Algorithm int    `json:"algorithm"`
	K         int    `json:"k"`
	Threads   int    `json:"n_threads"`
}

type bodyCovid struct {
	TestInput testInput `json:"covid"`
}

type res struct {
	Clase  int `json:"clase"`
	Ocurs0 int `json:"ocurs0"`
	Ocurs1 int `json:"ocurs1"`
}

/////////////////////           KMEANS

type body2 struct {
	K     int `json:"k"`
	MaxIt int `json:"max_it"`
}

type res2 struct {
	Centroids []person `json:"centroids"`
	Ncentroid []int    `json:"ncentroid"`
}

type resDiseases struct {
	Centroids []pacient `json:"centroids"`
	Ncentroid []int     `json:"ncentroid"`
}

/////////////////////          REGISTRAR PACIENTE
type body3 struct {
	Paciente pacient `json:"pacient"`
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
	//dftemp = normalize(dftemp)
	dftemp, _, _ = standarization(dftemp)
	newdftemp := make([][]float32, len(yTrain))
	for i := 0; i < len(newdftemp); i++ {
		newdftemp[i] = make([]float32, len(dftemp[0]))
		copy(newdftemp[i], dftemp[i])
	}
	copy(xTest, dftemp[len(dftemp)-1])

	if bdy.Algorithm == 1 {
		//fmt.Fprintf(response, "%.0f", monoKNN(dftemp, yTrain, xTest, bdy.K))
		clase, ocurs := monoKNN(newdftemp, yTrain, xTest, bdy.K)
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
		clase, ocurs := multiKNN(newdftemp, yTrain, xTest, bdy.K, bdy.Threads)
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
	fmt.Println("EMPIEZA DECODER")
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
	dftemp, mean_scales, std_scales := standarization(dftemp)
	maxIt := bdy.MaxIt
	fmt.Println("EMPIEZA KMEANS")
	head(dftemp, 1)
	G, centers, _ := multiKMeans(dftemp, k, maxIt)
	head(centers, 1)
	ocurs := make(map[int]int)
	var ncentroid []int
	var centroids []person
	for _, clase := range G {
		if _, found := ocurs[clase]; !found {
			ocurs[clase] = 1
		} else {
			ocurs[clase]++
		}
	}
	for _, val := range ocurs {
		ncentroid = append(ncentroid, val)
	}

	for i := 0; i < k; i++ {
		if !math.IsNaN(float64(centers[i][0])) {
			for j := 0; j < len(centers[i]); j++ {
				centers[i][j] = float32(float64(centers[i][j])*std_scales[j] + mean_scales[j])
			}
			centroids = append(centroids, person{Age: centers[i][0], Gender: centers[i][1], Height: centers[i][2],
				Weight: centers[i][3], Sbp: centers[i][4],
				Dbp: centers[i][5], Cholesterol: centers[i][6], Glucose: centers[i][7],
				Smoking: centers[i][8], AlcoholConsume: centers[i][9], PhysicalActivity: centers[i][10]})
		}
	}
	fmt.Println(centroids)
	fmt.Println(ncentroid)
	response := res2{Centroids: centroids, Ncentroid: ncentroid}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	fmt.Fprintf(r, "%s", jsonResponse)
}
func covidAnalysis(r http.ResponseWriter, request *http.Request) {
	var bdy bodyCovid
	err := json.NewDecoder(request.Body).Decode(&bdy)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	test := bdy.TestInput

	K := 1000
	Threads := 8

	xTest := []float32{test.Edad, test.Genero, test.Tos, test.Temperatura, test.DolorGarganta,
		test.MalestarGeneral}
	fmt.Println(xTest)

	dftemp := make([][]float32, len(xTrainCovid))
	for i := range xTrainCovid {
		dftemp[i] = make([]float32, len(xTrainCovid[i]))
		copy(dftemp[i], xTrainCovid[i])
	}
	dftemp = append(dftemp, xTest)
	//dftemp = normalize(dftemp)
	dftemp, _, _ = standarization(dftemp)

	newdftemp := make([][]float32, len(yTrainCovid))
	for i := 0; i < len(newdftemp); i++ {
		newdftemp[i] = make([]float32, len(dftemp[0]))
		copy(newdftemp[i], dftemp[i])
	}
	copy(xTest, dftemp[len(dftemp)-1])

	//fmt.Fprintf(response, "%.0f", monoKNN(dftemp, yTrain, xTest, bdy.K))
	clase, ocurs := multiKNN(newdftemp, yTrainCovid, xTest, K, Threads)
	response := res{Clase: int(clase), Ocurs0: ocurs[0], Ocurs1: ocurs[1]}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	fmt.Fprintf(r, "%s", jsonResponse)
}

func groupSelection(r http.ResponseWriter, request *http.Request) {
	var bdy body2
	fmt.Println("EMPIEZA DECODER")
	err := json.NewDecoder(request.Body).Decode(&bdy)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	k := bdy.K

	dftemp := make([][]float32, len(dfDiseases))
	for i := range dfDiseases {
		dftemp[i] = make([]float32, len(dfDiseases[i]))
		copy(dftemp[i], dfDiseases[i])
	}
	dftemp, mean_scales, std_scales := standarization(dftemp)
	maxIt := bdy.MaxIt
	fmt.Println("EMPIEZA KMEANS")
	head(dftemp, 1)
	G, centers, _ := multiKMeans(dftemp, k, maxIt)
	head(centers, 1)
	ocurs := make(map[int]int)
	var ncentroid []int
	groupsSelected = nil
	for _, clase := range G {
		if _, found := ocurs[clase]; !found {
			ocurs[clase] = 1
		} else {
			ocurs[clase]++
		}
	}
	for _, val := range ocurs {
		ncentroid = append(ncentroid, val)
	}

	for i := 0; i < k; i++ {
		if !math.IsNaN(float64(centers[i][0])) {
			for j := 0; j < len(centers[i]); j++ {
				z := float64(centers[i][j])*std_scales[j] + mean_scales[j]
				centers[i][j] = float32(z)
			}
			groupsSelected = append(groupsSelected, pacient{Edad: centers[i][0], Genero: centers[i][1], CardioDisease: centers[i][2],
				Diabetes: centers[i][3], RespDisease: centers[i][4],
				Hipertension: centers[i][5], Cancer: centers[i][6]})
		}
	}
	fmt.Println(groupsSelected)
	fmt.Println(ncentroid)
	response := resDiseases{Centroids: groupsSelected, Ncentroid: ncentroid}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	fmt.Fprintf(r, "%s", jsonResponse)
}

func registrarPaciente(r http.ResponseWriter, request *http.Request) {
	// REGISTRANDO PACIENTE

	var b3 body3
	fmt.Println("EMPIEZA DECODER")
	err := json.NewDecoder(request.Body).Decode(&b3)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		fmt.Println("ERROR GAA")
		return
	}
	dato := b3.Paciente
	datobytes, _ := json.Marshal(dato)

	newBlock := Block{}
	newBlock.Data = datobytes
	msg := tmsg{"serv", "localhost:8001", newBlock, BlockChain{}}

	remoteAddr := "localhost:8010"
	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("Can't dial", remoteAddr)
	} else {
		defer conn.Close()

		chRespuesta := make(chan string, 1)

		fmt.Println("Sending to", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
		go bcserver(chRespuesta, r)
		fmt.Fprintf(r, "%s", <-chRespuesta)
	}

}
func bcserver(chr chan string, r http.ResponseWriter) {
	localAddr := "localhost:8001"
	if ln, err := net.Listen("tcp", localAddr); err != nil {
		log.Panicln("Can't start listener on", localAddr)
	} else {
		defer ln.Close()
		fmt.Println("Listeing on", localAddr)
		if conn, err := ln.Accept(); err != nil {
			log.Println("Can't accept", conn.RemoteAddr())
		} else {
			go handle(chr, r, conn)
		}
	}
}
func handle(chr chan string, r http.ResponseWriter, conn net.Conn) {
	defer conn.Close()
	r2 := bufio.NewReader(conn)
	str, err := r2.ReadString('\n')
	if err == nil {
		fmt.Println("Recibido: ", str)
		chr <- str
	} else {
		fmt.Println("Error al leer")
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	// HEART DISEASE DATASET
	colnames, col, data = readCsvFile("cardio_train.csv")
	_, _, xTrain = sliceCols(colnames, col, data, colnames[1:len(colnames)-1])
	_, _, yTrain = sliceCols(colnames, col, data, []string{colnames[len(colnames)-1]})

	// COVID TESTING DATASET
	colnamesCovid, colCovid, dfCovid = readCsvFile("artifitial_covid_testing.csv")
	_, _, xTrainCovid = sliceCols(colnamesCovid, colCovid, dfCovid, colnamesCovid[:len(colnamesCovid)-1])
	_, _, yTrainCovid = sliceCols(colnamesCovid, colCovid, dfCovid, []string{colnamesCovid[len(colnamesCovid)-1]})

	// DISEASES DATASET
	colnamesDiseases, colDiseases, dfDiseases = readCsvFile("artifitial_pacients2.csv")

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
	router.HandleFunc("/covid_analysis", covidAnalysis).Methods("POST")
	router.HandleFunc("/group_selection", groupSelection).Methods("POST")
	router.HandleFunc("/register_pacient", registrarPaciente).Methods("POST")

	fmt.Println("Now server is running on port 8000")
	http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(router))
}
