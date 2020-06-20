import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import axios from "axios";
import { Formik, Field } from "formik";
import * as yup from "yup";
import InputRange from "react-input-range";
import "react-input-range/lib/css/index.css";

const GroupAnalysisSchema = yup.object({
  edad: yup
    .number()
    .required("La edad es requerida.")
    .integer()
    .lessThan(101, "La edad máxima es 100 años.")
    .moreThan(0, "La edad mínima es 1 año."),
  genero: yup
    .number()
    .required("Seleccione su género por favor")
    .integer()
    .moreThan(-1, "Seleccione su género por favor")
    .lessThan(2),
  cardio_disease: yup.bool().oneOf([true, false]),
  diabetes: yup.bool().oneOf([true, false]),
  resp_disease: yup.bool().oneOf([true, false]),
  hipertension: yup.bool().oneOf([true, false]),
  cancer: yup.bool().oneOf([true, false]),
});

class GroupAnalysis extends Component {
  state = {
    data: [],
  };

  render() {
    const { auth } = this.props;
    if (!auth.uid) return <Redirect to="/signin" />;

    return (
      <Formik
        initialValues={{
          edad: 0,
          genero: -1,
          cardio_disease: false,
          diabetes: false,
          resp_disease: false,
          hipertension: false,
          cancer: false,
        }}
        validationSchema={GroupAnalysisSchema}
        onSubmit={(values) => {
          const covid = {
            edad: parseInt(values.edad),
            genero: parseInt(values.genero),
            cardio_disease: values.cardio_disease ? 1 : 0,
            diabetes: values.diabetes ? 1 : 0,
            resp_disease: values.resp_disease ? 1 : 0,
            hipertension: values.hipertension ? 1 : 0,
            cancer: values.cancer ? 1 : 0,
          };
          const serv_url = "http://localhost:8000";
          console.log({ covid });
          axios({
            method: "post",
            url: serv_url + "/group_analysis",
            data: { covid },
          }).then((res) => {
            this.setState({ data: res.data });
            res.data.clase === 1
              ? alert(
                  "Usted puede sufrir de un ataque al corazón. Ocurrencias de la clase sana: " +
                    res.data.ocurs0 +
                    ", ocurrencias de la clase enferma: " +
                    res.data.ocurs1
                )
              : alert(
                  "Usted esta a salvo por ahora. Ocurrencias de la clase sana: " +
                    res.data.ocurs0 +
                    ", ocurrencias de la clase enferma: " +
                    res.data.ocurs1
                );
          });
        }}
      >
        {({
          handleSubmit,
          handleChange,
          values,
          errors,
          touched,
          status,
          setFieldValue,
        }) => (
          <form onSubmit={handleSubmit} className="blue lighten-4">
            <h5 className="grey-text text-darken-3">
              Análisis de Grupo de Riesgo
            </h5>
            <div className="input-field">
              <label htmlFor="edad">Edad (años)</label>
              {/* Edad */}
              <input
                type="number"
                onChange={handleChange}
                value={values.edad || ""}
                name="edad"
              />
              {errors.edad && touched.edad ? (
                <div>
                  <label className="red-text">{errors.edad}</label>
                </div>
              ) : null}
            </div>
            {/* Género */}
            <div className="">
              <label htmlFor="genero">
                <span>Género</span>
                <div>
                  <p>
                    <label>
                      <input
                        type="radio"
                        name="genero"
                        onChange={() => setFieldValue("genero", 0)}
                        value={values.genero || ""}
                        checked={values.genero === 0}
                      />
                      <span>Hombre</span>
                    </label>
                  </p>
                  <p>
                    <label>
                      <input
                        type="radio"
                        name="genero"
                        onChange={() => setFieldValue("genero", 1)}
                        value={values.genero || ""}
                        checked={values.genero === 1}
                      />
                      <span>Mujer</span>
                    </label>
                  </p>
                </div>
              </label>
              {errors.genero && touched.genero ? (
                <div>
                  <label className="red-text">{errors.genero}</label>
                </div>
              ) : null}
            </div>
            {/* Cardiovascular Disease */}
            <div>
              <label htmlFor="cardio_disease">
                <Field
                  type="checkbox"
                  name="cardio_disease"
                  id="cardio_disease"
                />
                <span>¿Tiene algún problema cardiovascular?</span>
              </label>
            </div>
            {/* Diabetes */}
            <div>
              <label htmlFor="diabetes">
                <Field type="checkbox" name="diabetes" id="diabetes" />
                <span>¿Tiene diabetes?</span>
              </label>
            </div>
            {/* Chronic respiratory disease */}
            <div>
              <label htmlFor="resp_disease">
                <Field type="checkbox" name="resp_disease" id="resp_disease" />
                <span>¿Tiene alguna enfermedad crónica respiratoria?</span>
              </label>
            </div>
            {/* Hypertension */}
            <div>
              <label htmlFor="hipertension">
                <Field type="checkbox" name="hipertension" id="hipertension" />
                <span>¿Sufre de hipertensión?</span>
              </label>
            </div>
            {/* Cancer */}
            <div>
              <label htmlFor="cancer">
                <Field type="checkbox" name="cancer" id="cancer" />
                <span>¿Tienes CANCER?</span>
              </label>
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

export default connect(mapStateToProps)(GroupAnalysis);
