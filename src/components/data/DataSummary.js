import React from "react";
import moment from "moment";

const DataSummary = ({ data }) => {
  return (
    <div className="card z-depth-0 data-summary">
      <div className="card-content grey-text text-darken-3">
        <span className="card-title">Person</span>
        <p>
          Registered by {data.authorFirstName} {data.authorLastName}
        </p>
        <p className="grey-text">
          {moment(data.createdAt.toDate()).calendar()}
        </p>
      </div>
    </div>
  );
};

export default DataSummary;
