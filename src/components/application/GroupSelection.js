import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import axios from "axios";
import { Formik, Field } from "formik";
import * as yup from "yup";

const KMeansSchema = yup.object({
  k: yup
    .number()
    .required("El número de clusters es requerido.")
    .integer()
    .lessThan(101, "El máximo número de clusters es de 100.")
    .moreThan(1, "El mínimo número de clusters es de 2."),
  max_it: yup
    .number()
    .required("El número máximo de iteraciones es requerido.")
    .integer()
    .moreThan(1, "El número de iteraciones debe ser mayor a 1."),
});

class GroupSelection extends Component {
  state = {
    data: [],
    centroids: [],
    items_centroid: [],
  };

  continueToGroupAnalysis = (e) => {
    localStorage.setItem("centroids", JSON.stringify(this.state.centroids));
    this.props.history.push("/group_analysis");
  };

  render() {
    const { auth } = this.props;
    if (!auth.uid) return <Redirect to="/signin" />;

    return (
      <Formik
        initialValues={{
          k: 0,
          max_it: 0,
        }}
        validationSchema={KMeansSchema}
        onSubmit={(values) => {
          this.setState({ centroids: [], items_centroid: [] });
          const serv_url = "http://localhost:8000";
          const k = values.k;
          const max_it = values.max_it;
          console.log({ k, max_it });
          axios({
            method: "post",
            url: serv_url + "/group_selection",
            data: { k, max_it },
          }).then((res) => {
            console.log(res);
            this.setState({
              centroids: res.data.centroids,
              items_centroid: res.data.ncentroid,
            });
          });
        }}
      >
        {({
          handleSubmit,
          handleChange,
          values,
          errors,
          touched,
          setFieldValue,
        }) => (
          <form onSubmit={handleSubmit} className="blue lighten-4">
            <h5 className="grey-text text-darken-3">
              Selección de Grupos de Riesgo
            </h5>
            <div className="input-field">
              <label htmlFor="k">Número de clusters o centroides</label>
              {/* K */}
              <input
                type="number"
                onChange={handleChange}
                value={values.k || ""}
                name="k"
              />
              {errors.k && touched.k ? (
                <div>
                  <label className="red-text">{errors.k}</label>
                </div>
              ) : null}
            </div>
            {/* Max Iterations */}
            <div className="input-field">
              <label htmlFor="max_it">Máximo de Iteraciones</label>
              <input
                type="number"
                onChange={handleChange}
                value={values.max_it || ""}
                name="max_it"
              />
              {errors.max_it && touched.max_it ? (
                <div>
                  <label className="red-text">{errors.max_it}</label>
                </div>
              ) : null}
            </div>
            <button type="submit" className="btn pink lighten-1 z-depth-0">
              Resultado
            </button>
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
                        Edad: {Math.round(c.edad)}(años), Género:{" "}
                        {Math.round(c.gender) === 1 ? "Mujer" : "Hombre"},{" "}
                        {Math.round(c.cardio_disease) === 1
                          ? "Problemas cardiovasculares: SÍ"
                          : "Problemas cardiovasculares: NO"}
                        ,{" "}
                        {Math.round(c.diabetes) === 1
                          ? "Diabetes: SÍ"
                          : "Diabetes: NO"}
                        ,{" "}
                        {Math.round(c.resp_disease) === 1
                          ? "Enfermedad respiratoria crónica: SÍ"
                          : "Enfermedad respiratoria crónica: NO"}
                        ,{" "}
                        {Math.round(c.hipertension) === 1
                          ? "Hipertensión: SÍ"
                          : "Hipertensión: NO"}
                        ,{" "}
                        {Math.round(c.cancer) === 1
                          ? "Cáncer: SÍ"
                          : "Cáncer: NO"}
                      </p>
                    </div>
                  );
                })}
            </div>
            {this.state.centroids && this.state.centroids.length > 0 ? (
              <div>
                <button
                  className="btn blue lighten-1 z-depth-0"
                  onClick={this.continueToGroupAnalysis}
                >
                  Siguiente
                </button>
              </div>
            ) : null}
          </form>
        )}
      </Formik>
    );
  }
}

const mapStateToProps = (state) => {
  return {
    auth: state.firebase.auth,
  };
};

export default connect(mapStateToProps)(GroupSelection);
