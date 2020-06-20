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
    .lessThan(6, "El máximo número de clusters es de 5.")
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
    k: 0,
    items_centroid: [],
    max_it: 0,
    submit: false,
  };

  handleChange = (e) => {
    this.setState({
      [e.target.id]: parseInt(e.target.value),
    });
    console.log(e.target.value);
  };

  handleSubmit = (e) => {
    e.preventDefault();
    if (this.state.is_correct === true) {
      /* this.setState({ centroids: [], items_centroid: [] });
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
      }); */
      console.log("Correcto");
    } else {
      console.log("Complete los datos");
    }
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
          const serv_url = "http://localhost:8000";
          const algorithm = this.state.algorithm;
          const k = this.state.k;
          const max_it = this.state.max_it;
          console.log({ k, max_it });
          this.props.history.push("/group_analysis");
          axios({
            method: "post",
            url: serv_url + "/group_selection",
            data: { k, max_it },
          })
            .then((res) => {
              console.log(res);
              this.setState({
                centroids: res.data.centroids,
                items_centroid: res.data.ncentroid,
              });
            })
            .then(() => {
              this.props.history.push("/group_analysis");
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
