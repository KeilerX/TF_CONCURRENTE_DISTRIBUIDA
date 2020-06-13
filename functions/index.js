const functions = require("firebase-functions");
const admin = require("firebase-admin");
admin.initializeApp();

// // Create and Deploy Your First Cloud Functions
// // https://firebase.google.com/docs/functions/write-firebase-functions
//
// exports.helloWorld = functions.https.onRequest((request, response) => {
//  response.send("Hello from Firebase!");
// });

const createNotification = (notification) => {
  return admin
    .firestore()
    .collection("notifications")
    .add(notification)
    .then((doc) => console.log("Notification added", doc));
};

exports.dataCreated = functions.firestore
  .document("patients/{patientId}")
  .onCreate((doc) => {
    const patient = doc.data();
    const notification = {
      content: "Agregó un nuevo registro",
      user: `${patient.authorFirstName} ${patient.authorLastName}`,
      time: admin.firestore.FieldValue.serverTimestamp(),
    };

    return createNotification(notification);
  });

exports.userJoined = functions.auth.user().onCreate(async (user) => {
  const doc = await admin.firestore().collection("users").doc(user.uid).get();
  const newUser = doc.data();
  const notification = {
    content: "Se unió",
    user: `${newUser.firstName} ${newUser.lastName}`,
    time: admin.firestore.FieldValue.serverTimestamp(),
  };
  return createNotification(notification);
});
