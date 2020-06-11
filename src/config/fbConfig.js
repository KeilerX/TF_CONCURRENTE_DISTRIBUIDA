import firebase from "firebase/app";
import "firebase/firestore";
import "firebase/auth";

// Your web app's Firebase configuration
var firebaseConfig = {
  apiKey: "AIzaSyAphiW4RQ6zVY7QLNkQlEmYFVcR-z_C0R8",
  authDomain: "go-react-pcd.firebaseapp.com",
  databaseURL: "https://go-react-pcd.firebaseio.com",
  projectId: "go-react-pcd",
  storageBucket: "go-react-pcd.appspot.com",
  messagingSenderId: "1082275678586",
  appId: "1:1082275678586:web:baae013dea807f21a9aec3",
  measurementId: "G-XP3NYM4R9S",
};
// Initialize Firebase
firebase.initializeApp(firebaseConfig);
firebase.firestore();

export default firebase;
