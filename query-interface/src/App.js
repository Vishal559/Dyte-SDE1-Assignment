import React from "react";
import { Provider } from "react-redux";
import store from "./redux/store";
import SearchBar from "./components/SearchBar";
import Visualization from "./components/Visualization";

function App() {
  return (
    <Provider store={store}>
      <div className="App" style={{ display: "flex", width: "100%", height: "100%"}}>
        <SearchBar />
        <Visualization />
      </div>
    </Provider>
  );
}

export default App;
