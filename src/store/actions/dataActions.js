export const createData = (data) => {
  return (dispatch, getState, { getFirebase, getFirestore }) => {
    //Make async call to database
    const firestore = getFirestore();
    const profile = getState().firebase.profile;
    const authorId = getState().firebase.auth.uid;
    firestore
      .collection("patients")
      .add({
        ...data,
        authorFirstName: profile.firstName,
        authorLastName: profile.lastName,
        authorId: authorId,
        createdAt: new Date(),
      })
      .then(() => {
        dispatch({ type: "CREATE_DATA", data: data });
      })
      .catch((err) => {
        dispatch({ type: "CREATE_DATA_ERROR", err });
      });
  };
};
