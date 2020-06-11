import React from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";
import Navbar from "./components/layout/Navbar";
import Dashboard from "./components/dashboard/Dashboard";
import DataDetails from "./components/data/DataDetails";
import SignIn from "./components/auth/SignIn";
import SignUp from "./components/auth/SignUp";
import CreateData from "./components/data/CreateData";
import KNN from "./components/application/KNN";
import KMeans from "./components/application/KMeans";

function App() {
  return (
    <BrowserRouter>
      <div className="App">
        <Navbar />
        <Switch>
          <Route exact path="/" component={Dashboard} />
          <Route path="/data/:id" component={DataDetails} />
          <Route path="/signin" component={SignIn} />
          <Route path="/signup" component={SignUp} />
          <Route path="/create" component={CreateData} />
          <Route path="/knn" component={KNN} />
          <Route path="/kmeans" component={KMeans} />
        </Switch>
      </div>
    </BrowserRouter>
  );
}

export default App;
