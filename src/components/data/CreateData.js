import React, { Component } from "react";
import { connect } from "react-redux";
import { createData } from "../../store/actions/dataActions";
import { Redirect } from "react-router-dom";

class CreateData extends Component {
  state = {
    age: "",
    has_disease: false,
  };

  handleChange = (e) => {
    const id = e.target.id;
    const value = e.target.value;
    console.log(value);
    if (id === "age") {
      parseInt(value);
    }
    this.setState({
      [id]: value,
    });
  };

  handleCheck = (e) => {
    this.setState({
      [e.target.id]: e.target.checked ? 1 : 0,
    });
  };

  handleSubmit = (e) => {
    e.preventDefault();
    //console.log(this.state);
    this.props.createData(this.state);
    this.props.history.push("/");
  };

  render() {
    const { auth } = this.props;
    if (!auth.uid) return <Redirect to="/signin" />;

    return (
      <div className="container">
        <form onSubmit={this.handleSubmit} className="white">
          <h5 className="grey-text text-darken-3">Register New Data</h5>
          <div className="input-field">
            <label htmlFor="age">Age</label>
            <input type="number" id="age" onChange={this.handleChange} />
          </div>
          <div className="">
            <label htmlFor="has_disease">
              <input
                type="checkbox"
                id="has_disease"
                defaultChecked={this.state.has_disease}
                onChange={this.handleCheck}
              />
              <span>Has Disease</span>
            </label>
          </div>
          <div className="input-field">
            <button className="btn pink lighten-1 z-depth-0">Register</button>
          </div>
        </form>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  return {
    auth: state.firebase.auth,
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    createData: (data) => dispatch(createData(data)),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(CreateData);
