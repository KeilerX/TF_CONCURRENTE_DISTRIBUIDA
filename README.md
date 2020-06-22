[TA2](#ta2) <br/>
[TF Hito 2](#tfhito2)

<a name="ta2"></a>
## Para la TA2:
Projecto realizado con [React App](https://github.com/facebook/create-react-app) y [Golang](https://golang.org/).
## Scripts disponibles

### npm install
Para instalar las dependencias.

### npm start
Inicia la aplicación en modo de desarrollador.<br />
Abrir [http://localhost:3000](http://localhost:3000) para verla en el browser.
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/npm_start.png)

### go run serv.go
En la carpeta **serv** se encuentra una aplicación en Golang que simula un servidor.<br />
Instalar las dependencias:
- go get github.com/gorilla/mux
- go get github.com/gorilla/handlers
- go get github.com/gonum/stat <br />
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/go_run_serv.png)

## Interfaces
### knn
Ir a [http://localhost:3000/knn](http://localhost:3000/knn) para probar el algoritmo KNN.<br />
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/knn.png)
Se deben completar los parámetros y seleccionar un algoritmo.

### kmeans
Ir a [http://localhost:3000/kmeans](http://localhost:3000/kmeans) para probar el algoritmo KMeans.<br />
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/kmeans.png)
Se debe especificar el número de centroides (clusters) a generar y el máximo de iteraciones a realizar.

<a name="tfhito2"></a>
## Para el Trabajo Final Hito 2
Projecto realizado con [React App](https://github.com/facebook/create-react-app) y [Golang](https://golang.org/).
## Scripts disponibles
En este projecto:

### npm install
Para instalar las dependencias.

### npm start
Inicia la aplicación en modo de desarrollador.<br />
Abrir [http://localhost:3000](http://localhost:3000) para verla en el browser.
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/npm_start.png)

### go run serv.go
En la carpeta **serv** se encuentra una aplicación en Golang que simula un servidor.<br />
Instalar las dependencias:
- go get github.com/gorilla/mux
- go get github.com/gorilla/handlers
- go get github.com/gonum/stat <br />
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/go_run_serv.png)

## Interfaces
### Análisis de Covid-19
Ir a [http://localhost:3000/covid_analysis](http://localhost:3000/covid_analysis) para probar el Análisis de COVID-19.<br />
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/covid_analysis.png)
Se deben completar los parámetros y se le indicará si puede o no estar contagiado de COVID-19.

### Selección de Grupos de Riesgo
Ir a [http://localhost:3000/group_selection](http://localhost:3000/group_selection) para probar la Selección de Grupos de Riesgo.<br />
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/group_selection.png)
Se debe especificar el número de grupos de riesgo (clusters) a generar y el número máximo de iteraciones que se realizarán.

### Análisis de Grupos de Riesgo
Ir a [http://localhost:3000/group_analysis](http://localhost:3000/group_analysis) para probar el Análisis de Grupos de Riesgo.<br />
![](https://github.com/KeilerX/TF_CONCURRENTE_DISTRIBUIDA/blob/master/imgs_test/group_analysis.png)
Se debe completar los datos y se le indicará a que grupo, generados en **Selección de Grupos de Riesgo**, se encuentra.
