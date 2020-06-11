const initState = {
  authError: null,
};
const authReducer = (state = initState, action) => {
  switch (action.type) {
    case "LOGIN_SUCCESS":
      console.log("Login success");
      return {
        ...state,
        authError: null,
      };
    case "LOGIN_ERROR":
      console.log("Login error");
      return {
        ...state,
        authError: "Login failed",
      };
    case "SIGN_OUT_SUCCESS":
      console.log("Sign out");
      return state;
    case "SIGN_UP_SUCCESS":
      console.log("Sign up success");
      return {
        ...state,
        authError: null,
      };
    case "SIGN_UP_ERROR":
      console.log("Sign up error");
      return {
        ...state,
        authError: "Sign up failed",
      };
    default:
      return state;
  }
};

export default authReducer;
