import axios from "axios";
import React from "react";

type User = {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  activated: boolean;
  created_at: string | Date;
};

type State = {
  isAuthenticated: boolean;
  loading: boolean;
  user: User | null;
};

type Action =
  | { type: "LOG_IN"; payload: User }
  | { type: "LOGOUT" }
  | { type: "AUTH_ERROR"; payload: Error }
  | { type: "USER_LOADED"; payload: User }
  | { type: "FETCHING_USER" };

type Dispatch = (action: Action) => void;

function authReducer(state: State, action: Action) {
  switch (action.type) {
    case "LOG_IN":
      return {
        ...state,
        user: action.payload,
        loading: false,
        isAuthenticated: true,
        error: null,
      };
    case "LOGOUT":
      return {
        ...state,
        isAuthenticated: false,
        user: null,
        loading: false,
        error: null,
      };
    case "USER_LOADED":
      return {
        ...state,
        isAuthenticated: true,
        user: action.payload,
        loading: false,
        error: null,
      };
    case "AUTH_ERROR":
      return {
        ...state,
        isAuthenticated: false,
        user: null,
        loading: false,
        error: action.payload,
      };
    case "FETCHING_USER":
      return {
        ...state,
        loading: true,
      };
  }
}

const AuthContext = React.createContext<
  { state: State; dispatch: Dispatch } | undefined
>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const initialState = {
    user: null,
    isAuthenticated: false,
    error: null,
    loading: false,
  };
  const [state, dispatch] = React.useReducer(authReducer, initialState);

  const authenticate = async () => {
    try {
      dispatch({ type: "FETCHING_USER" });
      const res = await axios.get("/api/v1/user/me");
      dispatch({
        type: "USER_LOADED",
        payload: res.data.data,
      });
    } catch (error) {
      dispatch({
        type: "AUTH_ERROR",
        payload: error as Error,
      });
    }
  };

  React.useEffect(() => {
    //@ts-ignore
    const Component = children.type;
    // If it doesn't require auth, everything is good
    if (!Component.requiresAuth) return;

    if (state.isAuthenticated) return;

    //TODO: If no cookie is present logout and redirect to login
    console.log("Checking auth");
    if (!state.loading) {
      authenticate();
    }
    //@ts-ignore
  }, [state.loading, state.isAuthenticated, children.type]);

  const value = { state, dispatch };
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = React.useContext(AuthContext);

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
