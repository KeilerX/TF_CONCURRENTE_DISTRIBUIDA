import React from "react";
import moment from "moment";

const PacientSummary = ({ data }) => {
  return (
    <div className="card z-depth-0 data-summary">
      <div className="card-content grey-text text-darken-3">
        <p>Edad: {data.edad}</p>
        <p>Género: {data.genero === 1 ? "Mujer" : "Hombre"}</p>
        <p>
          Problema Cardiovascular: {data.cardio_disease === 1 ? "Sí" : "No"}
        </p>
        <p>Diabetes: {data.diabetes === 1 ? "Sí" : "No"}</p>
        <p>
          Enfermedad crónica respiratoria:{" "}
          {data.resp_disease === 1 ? "Sí" : "No"}
        </p>
        <p>Hipertensión: {data.hipertension === 1 ? "Sí" : "No"}</p>
        <p>Cáncer: {data.cancer === 1 ? "Sí" : "No"}</p>
        <p>Hash Previo: {data.prev_hash}</p>
        <p>Mi Hash: {data.my_hash}</p>
        <p className="grey-text">
          {moment(data.createdAt.toDate()).calendar()}
        </p>
      </div>
    </div>
  );
};

export default PacientSummary;
