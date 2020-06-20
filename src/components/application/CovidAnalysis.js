import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import axios from "axios";
import { Formik, Field } from "formik";
import * as yup from "yup";
import InputRange from "react-input-range";
import "react-input-range/lib/css/index.css";

const CovidAnalysisSchema = yup.object({
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
  tos: yup
    .number()
    .required("El síntoma de tos es requerido.")
    .integer()
    .lessThan(4, "El grado de tos máximo es de 3.")
    .moreThan(-1, "El grado de tos mínimo es de 0."),
  temperatura: yup
    .number()
    .required("La temperatura es requerida.")
    .lessThan(41, "La temperatura máxima es de 40°.")
    .moreThan(35, "La temperatura mínima es de 36°."),
  dolor_garganta: yup
    .number()
    .required("El síntoma de tos es requerido.")
    .integer()
    .lessThan(4, "El grado de dolor de garganta máximo es de 3.")
    .moreThan(-1, "El grado de dolor de garganta mínimo es de 0."),
  malestar_general: yup
    .number()
    .required("Debe elegir una opción")
    .lessThan(4, "El grado de malestar general máximo es de 3.")
    .moreThan(-1, "El grado de malestar general mínimo es de 0."),
});

class CoivdAnalysis extends Component {
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
          tos: 0,
          temperatura: 36,
          dolor_garganta: 0,
          malestar_general: 0,
        }}
        validationSchema={CovidAnalysisSchema}
        onSubmit={(values) => {
          const covid = {
            edad: parseInt(values.edad),
            genero: parseInt(values.genero),
            tos: parseInt(values.tos),
            temperatura:
              values.temperatura % 1 >= 0.5
                ? Math.ceil(values.temperatura)
                : Math.floor(values.temperatura),
            dolor_garganta: parseInt(values.dolor_garganta),
            malestar_general: parseInt(values.malestar_general),
          };
          const serv_url = "http://localhost:8000";
          console.log({ covid });
          axios({
            method: "post",
            url: serv_url + "/covid_analysis",
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
            <h5 className="grey-text text-darken-3">Análisis de COVID-19</h5>
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
            {/* Tos */}
            <div className="field">
              <label htmlFor="tos">Grado de Tos</label>
              <ul style={{ marginLeft: "20px", marginRight: "20px" }}>
                <InputRange
                  maxValue={3}
                  minValue={0}
                  value={values.tos}
                  formatLabel={(value) =>
                    value === 0
                      ? "Sin Tos"
                      : value === 1
                      ? "Un poco de tos"
                      : value === 2
                      ? "Con tos"
                      : "Tos grave"
                  }
                  onChange={(e) => setFieldValue("tos", e)}
                />
                {errors.tos && touched.tos ? (
                  <div>
                    <label className="red-text">{errors.tos}</label>
                  </div>
                ) : null}
              </ul>
            </div>
            {/* Temperatura */}
            <div className="">
              <label htmlFor="temperatura">Temperatura</label>
              <ul style={{ marginLeft: "20px", marginRight: "20px" }}>
                <InputRange
                  step={0.1}
                  maxValue={40}
                  minValue={36}
                  value={values.temperatura}
                  formatLabel={(value) => `${value}°`}
                  onChange={(e) =>
                    setFieldValue("temperatura", Math.round(e * 10) / 10)
                  }
                />
                {errors.temperatura && touched.temperatura ? (
                  <div>
                    <label className="red-text">{errors.temperatura}</label>
                  </div>
                ) : null}
              </ul>
            </div>
            {/* Dolor de Garganta */}
            <div className="">
              <label htmlFor="dolor_garganta">Dolor de Garganta</label>
              <ul style={{ marginLeft: "20px", marginRight: "20px" }}>
                <InputRange
                  maxValue={3}
                  minValue={0}
                  value={values.dolor_garganta}
                  formatLabel={(value) =>
                    value === 0
                      ? "Sin dolor"
                      : value === 1
                      ? "Un poco de dolor"
                      : value === 2
                      ? "Con dolor"
                      : "Mucho dolor"
                  }
                  onChange={(e) => setFieldValue("dolor_garganta", e)}
                />
                {errors.dolor_garganta && touched.dolor_garganta ? (
                  <div>
                    <label className="red-text">{errors.dolor_garganta}</label>
                  </div>
                ) : null}
              </ul>
            </div>
            {/* Malestar General */}
            <div className="">
              <label htmlFor="malestar_general">Malestar General</label>
              <ul style={{ marginLeft: "20px", marginRight: "25px" }}>
                <InputRange
                  maxValue={3}
                  minValue={0}
                  value={values.malestar_general}
                  formatLabel={(value) =>
                    value === 0
                      ? "Sin malestar"
                      : value === 1
                      ? "Un poco de malestar"
                      : value === 2
                      ? "Con malestar"
                      : "Mucho malestar"
                  }
                  onChange={(e) => setFieldValue("malestar_general", e)}
                />
                {errors.malestar_general && touched.malestar_general ? (
                  <div>
                    <label className="red-text">
                      {errors.malestar_general}
                    </label>
                  </div>
                ) : null}
              </ul>
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

export default connect(mapStateToProps)(CoivdAnalysis);
