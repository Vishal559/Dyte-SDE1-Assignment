import { act } from "react-dom/test-utils";

const initialState = {
  searchQuery: "",
  enterPressed: false,
  results: [],
  filters: {
    level: "",
    message: "",
    resourceId: "",
    timestamp: "",
    traceId: "",
    spanId: "",
    commit: "",
    parentResourceId: "",
  },
  submitPressed: false,
  page: 1,
};

const rootReducer = (state = initialState, action) => {
  switch (action.type) {
    case "SET_SEARCH_QUERY":
      return {
        ...state,
        searchQuery: action.payload,
      };
    case "SET_PRESS_ENTER":
      return {
        ...state,
        enterPressed: action.payload,
      };
    case "UPDATE_RESULTS":
      return {
        ...state,
        results: action.payload,
      }
    case "SET_FILTERS":
      return {
        ...state,
        filters: {...state.filters, ...action.payload},
      };
    case "SET_PRESS_SUBMIT":
      return {
        ...state,
        submitPressed: action.payload
      }
    case "SET_PAGE":
      return {
        ...state,
        page: action.payload
      }
    default:
      return state;
  }
};

export default rootReducer;
