import React from "react";
import { connect } from "react-redux";
import { firestoreConnect } from "react-redux-firebase";
import { compose } from "redux";
import { Redirect } from "react-router-dom";
import moment from "moment";

const DataDetails = (props) => {
  const { data, auth } = props;
  if (!auth.uid) return <Redirect to="/signin" />;
  if (data) {
    return (
      <div className="container section data-details">
        <div className="card z-depth-0">
          <div className="card-cotent">
            <span className="card-title">{data.id}</span>
            <p>{data.age}</p>
            <p>{data.has_disease ? "has disease" : "no disease"}</p>
          </div>
          <div className="card-action grey lighten-4 grey-text">
            <div>
              Registered by {data.authorFirstName} {data.authorLastName}
            </div>
            <div>{moment(data.createdAt.toDate()).calendar()}</div>
          </div>
        </div>
      </div>
    );
  } else {
    return (
      <div className="container center">
        <p>Loading data...</p>
      </div>
    );
  }
};

const mapStateToProps = (state, ownProps) => {
  const id = ownProps.match.params.id;
  const data = state.firestore.data.patients;
  const d = data ? data[id] : null;
  return {
    data: d,
    auth: state.firebase.auth,
  };
};

export default compose(
  connect(mapStateToProps),
  firestoreConnect([{ collection: "patients" }])
)(DataDetails);
