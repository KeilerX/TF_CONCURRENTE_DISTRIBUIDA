import React, { Component } from "react";
import PacientList from "../pacient/PacientList";
import { connect } from "react-redux";
import { firestoreConnect } from "react-redux-firebase";
import { compose } from "redux";
import { Redirect } from "react-router-dom";

class PacientDashboard extends Component {
  render() {
    const { data8010, data8011, data8012, data8013, auth } = this.props;
    if (!auth.uid) return <Redirect to="/signin" />;
    return (
      <div className="dashboard container">
        <div className="row">
          <div className="col s12 m3">
            <h5 className="white-text">Block Chain 1</h5>
            <PacientList data={data8010} />
          </div>
          <div className="col s12 m3">
            <h5 className="white-text">Block Chain 2</h5>
            <PacientList data={data8011} />
          </div>
          <div className="col s12 m3">
            <h5 className="white-text">Block Chain 3</h5>
            <PacientList data={data8012} />
          </div>
          <div className="col s12 m3">
            <h5 className="white-text">Block Chain 4</h5>
            <PacientList data={data8013} />
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  //console.log(state);
  return {
    data8010: state.firestore.ordered.pacientes8010,
    data8011: state.firestore.ordered.pacientes8011,
    data8012: state.firestore.ordered.pacientes8012,
    data8013: state.firestore.ordered.pacientes8013,
    auth: state.firebase.auth,
  };
};

export default compose(
  connect(mapStateToProps),
  firestoreConnect([
    { collection: "pacientes8010", orderBy: ["CreatedAt", "desc"] },
    { collection: "pacientes8011", orderBy: ["CreatedAt", "desc"] },
    { collection: "pacientes8012", orderBy: ["CreatedAt", "desc"] },
    { collection: "pacientes8013", orderBy: ["CreatedAt", "desc"] },
  ])
)(PacientDashboard);
