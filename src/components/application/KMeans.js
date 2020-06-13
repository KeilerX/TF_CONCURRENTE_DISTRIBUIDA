import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import axios from "axios";

class KMeans extends Component {
  state = {
    data: [],
    k: 0,
    centroids: [],
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
      console.log(centroid);
      this.setState((prevState) => ({
        centroids: [...prevState.centroids, centroid],
      }));
    }
  };

  handleSubmit = (e) => {
    e.preventDefault();
    console.log(this.state);
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
          <button
            onClick={this.generateCentroids}
            className="btn blue darken-2 z-depth-0"
          >
            Generar Centroides
          </button>
          <div>
            {this.state.centroids.map((c) => {
              return (
                <div key={c.id}>
                  <p>Centroide {c.id + 1}</p>
                  <p>
                    Edad: {c.age}(años), Altura: {c.height}(cm), Peso:{" "}
                    {c.weight}(kg), Género:{" "}
                    {c.gender === 1 ? "Mujer" : "Hombre"}, Presión arterial
                    sitólica: {c.sbp}, Presión arterial diastólica: {c.dbp},
                    Colesterol:{" "}
                    {c.cholesterol === 1
                      ? "Normal"
                      : c.cholesterol === 2
                      ? "Por encima de lo normal"
                      : "Muy por encima de lo normal"}
                    , Glucosa:{" "}
                    {c.glucose === 1
                      ? "Normal"
                      : c.cholesterol === 2
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
        <div>
          <button className="btn pink lighten-1 z-depth-0">Resultado</button>
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
