import React from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";
import Navbar from "./components/layout/Navbar";
import Dashboard from "./components/dashboard/Dashboard";
import PacientDashboard from "./components/dashboard/PacientDashboard";
import DataDetails from "./components/data/DataDetails";
import SignIn from "./components/auth/SignIn";
import SignUp from "./components/auth/SignUp";
import CreateData from "./components/data/CreateData";
import KNN from "./components/application/KNN";
import KMeans from "./components/application/KMeans";
import GroupSelection from "./components/application/GroupSelection";
import GroupAnalysis from "./components/application/GroupAnalysis";
import CovidAnalysis from "./components/application/CovidAnalysis";
import CreatePacient from "./components/application/CreatePacient";

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
          <Route path="/covid_analysis" component={CovidAnalysis} />
          <Route path="/group_selection" component={GroupSelection} />
          <Route path="/group_analysis" component={GroupAnalysis} />
          <Route path="/create_pacient" component={CreatePacient} />
          <Route path="/dashboard_pacients" component={PacientDashboard} />
          <Route path="/knn" component={KNN} />
          <Route path="/kmeans" component={KMeans} />
        </Switch>
      </div>
    </BrowserRouter>
  );
}

export default App;
