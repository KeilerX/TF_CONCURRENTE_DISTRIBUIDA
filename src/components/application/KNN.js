import React, { Component } from "react";
import { connect } from "react-redux";
import { Redirect } from "react-router-dom";
import axios from "axios";
import { Formik, Field } from "formik";
import * as yup from "yup";

const PersonSchema = yup.object({
  age: yup
    .number()
    .required("La edad es requerida.")
    .integer()
    .lessThan(101, "La edad máxima es 100 años")
    .moreThan(17, "La edad mínima es 18 años"),
  height: yup
    .number()
    .required("La altura es requerida")
    .integer()
    .moreThan(139, "La altura mínima es 140 (cm)")
    .lessThan(211, "La altura máxima es 210 (cm)"),
  weight: yup
    .number()
    .required("El peso es requerido")
    .integer()
    .moreThan(39, "El peso mínimo es 40 (kg)")
    .lessThan(151, "El peso máximo es 150 (kg)"),
  gender: yup
    .number()
    .required("Seleccione su género por favor")
    .integer()
    .moreThan(0, "Seleccione su género por favor")
    .lessThan(3),
  sbp: yup
    .number()
    .required("Debe elegir una opción")
    .integer()
    .moreThan(99, "La mínima presión arterial sistólica es 100")
    .lessThan(171, "La máxima La presión arterial sistólica 170"),
  dbp: yup
    .number()
    .required("La presión arterial diastólica es requerida")
    .integer()
    .moreThan(59, "La mínima presión arterial diastólica es 60")
    .lessThan(101, "La máxima La presión arterial diastólica 100"),
  cholesterol: yup
    .number()
    .required("Debe elegir una opción")
    .moreThan(0, "Debe elegir una opción")
    .lessThan(4),
  glucose: yup
    .number()
    .required("Debe elegir una opción")
    .moreThan(0, "Debe elegir una opción")
    .lessThan(4),
  smoking: yup.bool().oneOf([true, false]),
  alcohol_consume: yup.bool().oneOf([true, false]),
  physical_activity: yup.bool().oneOf([true, false]),
});

class KNN extends Component {
  state = {
    data: [],
  };
  render() {
    const { auth } = this.props;
    if (!auth.uid) return <Redirect to="/signin" />;

    return (
      <Formik
        initialValues={{
          age: 0,
          height: 0,
          weight: 0,
          gender: 0,
          sbp: 0,
          dbp: 0,
          cholesterol: 0,
          glucose: 0,
          smoking: false,
          alcohol_consume: false,
          physical_activity: false,
        }}
        validationSchema={PersonSchema}
        onSubmit={(values) => {
          const person = {
            age: parseInt(values.age),
            height: parseInt(values.height),
            weight: parseInt(values.weight),
            gender: parseInt(values.gender),
            sbp: parseInt(values.sbp),
            dbp: parseInt(values.dbp),
            cholesterol: parseInt(values.cholesterol),
            glucose: parseInt(values.glucose),
            smoking: values.smoking ? 1 : 0,
            alcohol_consume: values.alcohol_consume ? 1 : 0,
            physical_activity: values.physical_activity ? 1 : 0,
          };
          console.log(person);
          const serv_url = "http://localhost:8000";
          axios({
            method: "post",
            url: serv_url + "/knn",
            data: {
              age: parseInt(values.age),
              height: parseInt(values.height),
              weight: parseInt(values.weight),
              gender: parseInt(values.gender),
              sbp: parseInt(values.sbp),
              dbp: parseInt(values.dbp),
              cholesterol: parseInt(values.cholesterol),
              glucose: parseInt(values.glucose),
              smoking: values.smoking ? 1 : 0,
              alcohol_consume: values.alcohol_consume ? 1 : 0,
              physical_activity: values.physical_activity ? 1 : 0,
            },
          }).then((res) => {
            this.setState({ data: res.data });
            console.log(res);
            console.log(res.data);
          });
          /* axios.post("http://localhost:8000/knn", { person }).then((res) => {
            this.setState({ data: res.data });
            console.log(res);
            console.log(res.data);
          }); */
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
            <h5 className="grey-text text-darken-3">KNN Algorithm</h5>
            <div className="input-field">
              <label htmlFor="age">Edad (años)</label>
              {/* Age */}
              <input
                type="number"
                onChange={handleChange}
                value={values.age || ""}
                name="age"
              />
              {errors.age && touched.age ? (
                <div>
                  <label className="red-text">{errors.age}</label>
                </div>
              ) : null}
            </div>
            {/* Height */}
            <div className="input-field">
              <label htmlFor="height">Altura (cm)</label>
              <input
                type="number"
                onChange={handleChange}
                value={values.height || ""}
                name="height"
              />
              {errors.height && touched.height ? (
                <div>
                  <label className="red-text">{errors.height}</label>
                </div>
              ) : null}
            </div>
            {/* Weight */}
            <div className="input-field">
              <label htmlFor="weight">Peso (kg)</label>
              <input
                type="number"
                onChange={handleChange}
                value={values.weight || ""}
                name="weight"
              />
              {errors.weight && touched.weight ? (
                <div>
                  <label className="red-text">{errors.weight}</label>
                </div>
              ) : null}
            </div>
            {/* Gender */}
            <div className="">
              <label htmlFor="gender">
                <span>Género</span>
                <div>
                  <p>
                    <label>
                      <input
                        type="radio"
                        name="gender"
                        onChange={() => setFieldValue("gender", 2)}
                        value={values.gender || ""}
                        checked={values.gender === 2}
                      />
                      <span>Hombre</span>
                    </label>
                  </p>
                  <p>
                    <label>
                      <input
                        type="radio"
                        name="gender"
                        onChange={() => setFieldValue("gender", 1)}
                        value={values.gender || ""}
                        checked={values.gender === 1}
                      />
                      <span>Mujer</span>
                    </label>
                  </p>
                </div>
              </label>
              {errors.gender && touched.gender ? (
                <div>
                  <label className="red-text">{errors.gender}</label>
                </div>
              ) : null}
            </div>
            {/* Systolic Blood Pressure */}
            <div className="input-field">
              <label htmlFor="sbp">Presión arterial sistólica</label>
              <input
                type="number"
                name="sbp"
                onChange={handleChange}
                value={values.sbp || ""}
              />
              {errors.sbp && touched.sbp ? (
                <div>
                  <label className="red-text">{errors.sbp}</label>
                </div>
              ) : null}
            </div>
            {/* Diastolic Blood Pressure */}
            <div className="input-field">
              <label htmlFor="dbp">Presión arterial diastólica</label>
              <input
                type="number"
                name="dbp"
                onChange={handleChange}
                value={values.dbp || ""}
              />
              {errors.dbp && touched.dbp ? (
                <div>
                  <label className="red-text">{errors.dbp}</label>
                </div>
              ) : null}
            </div>
            {/* Cholesterol */}
            <div className="input-field">
              <div>
                <label htmlFor="cholesterol">Colesterol</label>
                <select
                  className="browser-default"
                  name="cholesterol"
                  onChange={handleChange}
                  value={values.cholesterol}
                >
                  <option disabled value={0} defaultValue>
                    Elige una opción
                  </option>
                  <option value={1}>Normal</option>
                  <option value={2}>Por encima de lo normal</option>
                  <option value={3}>Muy por encima de lo normal</option>
                </select>
              </div>
              {errors.cholesterol && touched.cholesterol ? (
                <div>
                  <label className="red-text">{errors.cholesterol}</label>
                </div>
              ) : null}
            </div>
            {/* Glucose */}
            <div className="input-field">
              <div>
                <label htmlFor="glucose">Glucosa</label>
                <select
                  className="browser-default"
                  name="glucose"
                  onChange={handleChange}
                  value={values.glucose}
                >
                  <option disabled value={0} defaultValue>
                    Elige una opción
                  </option>
                  <option value={1}>Normal</option>
                  <option value={2}>Por encima de lo normal</option>
                  <option value={3}>Muy por encima de lo normal</option>
                </select>
              </div>
              {errors.glucose && touched.glucose ? (
                <div>
                  <label className="red-text">{errors.glucose}</label>
                </div>
              ) : null}
            </div>
            {/* Smoking */}
            <div>
              <label htmlFor="smoking">
                <Field type="checkbox" name="smoking" id="smoking" />
                <span>¿Fuma?</span>
              </label>
            </div>
            {/* Alcohol Consume */}
            <div>
              <label htmlFor="alcohol_consume">
                <Field
                  type="checkbox"
                  name="alcohol_consume"
                  id="alcohol_consume"
                />
                <span>¿Consume alcohol?</span>
              </label>
            </div>
            {/* Physical Activity */}
            <div>
              <label htmlFor="physical_activity">
                <Field
                  type="checkbox"
                  name="physical_activity"
                  id="physical_activity"
                />
                <span>¿Realiza alguna activida física?</span>
              </label>
            </div>
            <button type="submit" className="btn pink lighten-1 z-depth-0">
              Resultado
            </button>
            {/* {this.state.data &&
              this.state.data.map((d) => {
                return (
                  <p key={d.id}>
                    {d.nombre} {d.apellido}
                  </p>
                );
              })} */}
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

export default connect(mapStateToProps)(KNN);
