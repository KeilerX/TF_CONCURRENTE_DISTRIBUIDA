import React from "react";
import PacientSummary from "./PacientSummary";
import { Link } from "react-router-dom";

const PacientList = ({ data }) => {
  return (
    <div className="data-list section">
      {data &&
        data.map((d) => {
          return <PacientSummary data={d} />;
        })}
    </div>
  );
};

export default PacientList;
