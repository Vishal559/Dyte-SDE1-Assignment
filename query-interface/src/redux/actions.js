export const setSearchQuery = (query) => {
  return {
    type: "SET_SEARCH_QUERY",
    payload: query,
  };
};

export const setPressEnter = (isPressed) => {
  return {
    type: "SET_PRESS_ENTER",
    payload: isPressed,
  };
};

export const updateResults = (results) => { 
  return {
    type: "UPDATE_RESULTS",
    payload: results,
  };
};

export const setFilters = (filters) => {
  return {
    type: "SET_FILTERS",
    payload: filters,
  };
};

export const setPressSubmit = (filters) => {
  return {
    type: "SET_PRESS_SUBMIT",
    payload: filters,
  }
}