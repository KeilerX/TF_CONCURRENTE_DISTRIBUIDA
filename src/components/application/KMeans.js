import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import axios from "axios";

class KMeans extends Component {
  state = {
    data: [],
    k: 0,
    centroids: [],
    items_centroid: [],
    max_it: 0,
  };

  handleChange = (e) => {
    this.setState({
      [e.target.id]: parseInt(e.target.value),
    });
    console.log(e.target.value);
  };

  generateCentroids = (e) => {
    e.preventDefault();
    this.setState({ centroids: [] });
    for (let i = 0; i < this.state.k; i++) {
      const centroid = {
        id: parseInt(i),
        age: parseInt(30 + Math.random() * (100 - 30)),
        height: parseInt(140 + Math.random() * (210 - 140)),
        weight: parseInt(40 + Math.random() * (150 - 40)),
        gender: parseInt(Math.random() * 2),
        sbp: parseInt(100 + Math.random() * (170 - 100)),
        dbp: parseInt(60 + Math.random() * (100 - 60)),
        cholesterol: parseInt(1 + Math.random() * (3 - 1)),
        glucose: parseInt(1 + Math.random() * (3 - 1)),
        smoking: Math.random() >= 0.5 ? 1 : 0,
        alcohol_consume: Math.random() >= 0.5 ? 1 : 0,
        physical_activity: Math.random() >= 0.5 ? 1 : 0,
      };
      this.setState((prevState) => {
        prevState.centroids = [...prevState.centroids, centroid];
      });
      console.log(centroid);
    }
  };

  handleSubmit = (e) => {
    this.setState({ centroids: [], items_centroid: [] });
    e.preventDefault();
    const serv_url = "http://localhost:8000";
    const k = this.state.k;
    const max_it = this.state.max_it;
    console.log({ k });
    axios({
      method: "post",
      url: serv_url + "/kmeans",
      data: { k, max_it },
    }).then((res) => {
      console.log(res);
      this.setState({
        centroids: res.data.centroids,
        items_centroid: res.data.ncentroid,
      });
    });
  };

  render() {
    const { auth } = this.props;
    if (!auth.uid) return <Redirect to="/signin" />;

    const { k } = this.state;
    const message_k = k <= 0 || k > 5 ? true : false;
    return (
      <form onSubmit={this.handleSubmit} className="blue lighten-4">
        <h5 className="grey-text text-darken-3">Algoritmo KMeans</h5>
        <div className="input-field">
          <label htmlFor="age">Número de Centroides</label>
          {/* Número de Centroides */}
          <input
            type="number"
            id="k"
            value={this.state.k || ""}
            onChange={this.handleChange}
          />
        </div>
        {message_k ? (
          <label className="red-text">
            El número de centroides debe estar entre 1 y 5
          </label>
        ) : null}

        <div>
          <div className="input-field">
            <label htmlFor="max_it">Máximo de iteraciones</label>
            {/* Máximo de iteraciones */}
            <input
              type="number"
              id="max_it"
              value={this.state.max_it || ""}
              onChange={this.handleChange}
            />
          </div>
          <div>
            <button className="btn pink lighten-1 z-depth-0">Resultado</button>
          </div>
          {/* <div>
            <button
              onClick={this.generateCentroids}
              className="btn green lighten-1 z-depth-0"
            >
              Generar Centroides
            </button>
          </div> */}
          <div>
            {this.state.centroids &&
              this.state.centroids.map((c, i) => {
                return (
                  <div key={i}>
                    <p>
                      Centroide {i + 1} con {this.state.items_centroid[i]}{" "}
                      registros
                    </p>
                    <p>
                      Edad: {parseInt(c.age)}(años), Altura:{" "}
                      {parseInt(c.height)}(cm), Peso: {parseInt(c.weight)}
                      (kg), Género:{" "}
                      {parseInt(c.gender) === 1 ? "Mujer" : "Hombre"}, Presión
                      arterial sitólica: {parseInt(c.sbp)}, Presión arterial
                      diastólica: {parseInt(c.dbp)}, Colesterol:{" "}
                      {parseInt(c.cholesterol) === 1
                        ? "Normal"
                        : c.cholesterol === 2
                        ? "Por encima de lo normal"
                        : "Muy por encima de lo normal"}
                      , Glucosa:{" "}
                      {c.glucose === 1
                        ? "Normal"
                        : parseInt(c.cholesterol) === 2
                        ? "Por encima de lo normal"
                        : "Muy por encima de lo normal"}
                      , {c.smoking === 1 ? "Fuma" : "No fuma"},{" "}
                      {c.alcohol_consume === 1
                        ? "Consume alcohol"
                        : "No consume alcohol"}
                      ,{" "}
                      {c.physical_activity === 1
                        ? "Realiza actividad física"
                        : "No realiza actividad física"}
                    </p>
                  </div>
                );
              })}
          </div>
        </div>
      </form>
    );
  }
}

const mapStateToProps = (state) => {
  return {
    auth: state.firebase.auth,
  };
};

export default connect(mapStateToProps)(KMeans);
