import { React, useState } from "react";
import axios from "axios";
import { useDispatch, useSelector } from "react-redux";
import { setSearchQuery, setPressEnter, updateResults , setFilters} from "../redux/actions";
import "../App.css";

const SearchBar = () => {
  const dispatch = useDispatch();
  const searchQuery = useSelector((state) => state.searchQuery);
  const filters = useSelector((state) => state.filters)


  const handleInputChange = (e) => {
    dispatch(setSearchQuery(e.target.value));
    dispatch(setPressEnter(false));
  };

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    dispatch(setFilters({ [name]: value}));
  };

  const handleKeyPress = (e) => {
    // Trigger text query on Enter key press
    if (e.key === "Enter" && searchQuery.trim() !== "") {
      console.log("Text query triggered:", searchQuery);
      dispatch(setPressEnter(true));
    }
  };

  const handleSubmit = async () => {
    dispatch(updateResults([]));
    dispatch(setPressEnter(true));
  };


  return (
    <div className="search-page-container">
      <div className="input-bar">
        <input
          type="text"
          value={searchQuery}
          onChange={handleInputChange}
          onKeyPress={handleKeyPress}
          placeholder="Enter your search query..."
        />
      </div>
      {/* Filters */}
      <div className="filters-container">
        <div>
          <label>Level:</label>
          <input
            type="text"
            name="level"
            value={filters.level}
            onChange={handleFilterChange}
          />
        </div>

        <div>
          <label>Message:</label>
          <input
            type="text"
            name="message"
            value={filters.message}
            onChange={handleFilterChange}
          />
        </div>

        <div>
          <label>Resource ID:</label>
          <input
            type="text"
            name="resourceId"
            value={filters.resourceId}
            onChange={handleFilterChange}
          />
        </div>

        <div>
          <label>Timestamp:</label>
          <input
            type="text"
            name="timestamp"
            value={filters.timestamp}
            onChange={handleFilterChange}
          />
        </div>

        <div>
          <label>Trace ID:</label>
          <input
            type="text"
            name="traceId"
            value={filters.traceId}
            onChange={handleFilterChange}
          />
        </div>

        <div>
          <label>Span ID:</label>
          <input
            type="text"
            name="spanId"
            value={filters.spanId}
            onChange={handleFilterChange}
          />
        </div>

        <div>
          <label>Commit:</label>
          <input
            type="text"
            name="commit"
            value={filters.commit}
            onChange={handleFilterChange}
          />
        </div>

        <div>
          <label>Parent Resource ID:</label>
          <input
            type="text"
            name="parentResourceId"
            value={filters.parentResourceId}
            onChange={handleFilterChange}
          />
        </div>

        <button className="button" onClick={handleSubmit} style={{
          backgroundColor: "gray",
          border: "none",
          borderRadius: "2px",
          color: "white",
          padding: "15px 32px",
          textAlign: "center",
          textDecoration: "none",
          display: "inline-block",
          fontSize: "16px",
          margin: "4px 2px",
          cursor: "pointer",
        }}>Submit</button>
      </div>
    </div>
  );
};

export default SearchBar;
