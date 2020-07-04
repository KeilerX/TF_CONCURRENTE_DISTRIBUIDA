import React from "react";
import moment from "moment";

const PacientSummary = ({ data }) => {
  var CreatedAt = data.CreatedAt;
  let edad;
  let genero;
  let cardio_disease;
  let diabetes;
  let resp_disease;
  let hipertension;
  let cancer;
  if (data.Data !== "Genesis") {
    var pacientJson = JSON.parse(data.Data);
    edad = pacientJson.edad;
    genero = pacientJson.genero;
    cardio_disease = pacientJson.cardio_disease;
    diabetes = pacientJson.diabetes;
    resp_disease = pacientJson.resp_disease;
    hipertension = pacientJson.hipertension;
    cancer = pacientJson.cancer;
  }
  let hash = data.Hash;

  //var utf8 = unescape(encodeURIComponent(hash));
  var arr = [];
  for (var i = 0; i < hash.length; i++) {
    arr.push(hash.charCodeAt(i));
  }
  let hashPromise = decodeURIComponent(escape(hash));

  //console.log(arr, hashPromise);

  let prev_hash = data.PrevHash;
  let prevHashPromise = decodeURIComponent(escape(prev_hash));

  return (
    <div className="card z-depth-0 data-summary">
      <div className="card-content grey-text text-darken-3">
        {data.Data !== "Genesis" ? (
          <div>
            <p>Edad: {edad}</p>
            <p>Género: {genero === 1 ? "Mujer" : "Hombre"}</p>
            <p>Problema Cardiovascular: {cardio_disease === 1 ? "Sí" : "No"}</p>
            <p>Diabetes: {diabetes === 1 ? "Sí" : "No"}</p>
            <p>
              Enfermedad crónica respiratoria:{" "}
              {resp_disease === 1 ? "Sí" : "No"}
            </p>
            <p>Hipertensión: {hipertension === 1 ? "Sí" : "No"}</p>
            <p>Cáncer: {cancer === 1 ? "Sí" : "No"}</p>
            <ul>
              Hash Previo: <p style={{ fontSize: "11px" }}>{prevHashPromise}</p>
            </ul>
            <ul>
              Mi Hash: <p style={{ fontSize: "11px" }}>{hashPromise}</p>
            </ul>
            <p className="grey-text">{moment(CreatedAt.toDate()).calendar()}</p>
          </div>
        ) : (
          <div>
            <p>Nombre: {data.Data}</p>
            <ul>
              Hash Previo: <p style={{ fontSize: "11px" }}>{prevHashPromise}</p>
            </ul>
            <ul>
              Mi Hash: <p style={{ fontSize: "11px" }}>{hashPromise}</p>
            </ul>
          </div>
        )}
      </div>
    </div>
  );
};

export default PacientSummary;
