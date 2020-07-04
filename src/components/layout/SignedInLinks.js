import React, { Component } from "react";
import { NavLink } from "react-router-dom";
import { connect } from "react-redux";
import { signOut } from "../../store/actions/authActions";

const SignedInLinks = (props) => {
  const open = false;

  const handleButtonTf = () => {
    open = !open;
  };

  return (
    <ul className="right">
      <li>
        <NavLink to="/covid_analysis">Análisis Covid</NavLink>
      </li>
      <li>
        <NavLink to="/group_selection">Grupo Riesgo</NavLink>
      </li>
      <li>
        <NavLink to="/create_pacient">Registar Paciente</NavLink>
      </li>
      <li>
        <NavLink to="/dashboard_pacients">Lista Pacientes</NavLink>
      </li>
      <li>
        <NavLink to="/knn">KNN</NavLink>
      </li>
      <li>
        <NavLink to="/kmeans">KMens</NavLink>
      </li>
      <li>
        <a href="/" onClick={props.signOut}>
          Log Out
        </a>
      </li>
      <li>
        <NavLink to="/" className="btn btn-floating pink lighten-1">
          {props.profile.initials}
        </NavLink>
      </li>
    </ul>
  );
};

const mapDispatchToProps = (dispatch) => {
  return {
    signOut: () => dispatch(signOut()),
  };
};

export default connect(null, mapDispatchToProps)(SignedInLinks);
