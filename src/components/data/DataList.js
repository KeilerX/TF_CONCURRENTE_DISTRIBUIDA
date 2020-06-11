import React from "react";
import DataSummary from "./DataSummary";
import { Link } from "react-router-dom";

const DataList = ({ data }) => {
  return (
    <div className="data-list section">
      {data &&
        data.map((d) => {
          return (
            <Link to={"/data/" + d.id} key={d.id}>
              <DataSummary data={d} />;
            </Link>
          );
        })}
    </div>
  );
};

export default DataList;
