import React, { useEffect, useState } from "react";
import axios from "axios";
import { useDispatch, useSelector } from "react-redux";
import { setPressEnter, setPressSubmit, updateResults} from "../redux/actions";

const Visualization = () => {
  const dispatch = useDispatch();
  const searchQuery = useSelector((state) => state.searchQuery);
  const enterPressed = useSelector((state) => state.enterPressed);
  const submitPressed = useSelector((state) => state.submitPressed);
  const results = useSelector((state) => state.results); 
  const filters = useSelector((state) => state.filters)
  const [page, setPage] = useState(1);
  const [isFetching, setIsFetching] = useState(false)

  useEffect(() => {
    if (enterPressed || submitPressed) {
      // Reset enterPressed to false after making the query
      dispatch(setPressEnter(false));
      dispatch(setPressSubmit(false));
      setIsFetching(true)

      // Check if any field is non-empty
      const hasNonEmptyField = Object.values(filters).some(value => Boolean(value));
      // Make a text query when the searchQuery changes
      if (searchQuery.trim() !== "" || hasNonEmptyField) {
        makeTextQuery(searchQuery, page);
      } else {
        // Clear results when searchQuery is empty
        dispatch(updateResults([]));
      }
    }
  }, [searchQuery, enterPressed, submitPressed, page]);

  const makeTextQuery = (query, pageNum) => {
    try {
      // Set the Content-Type header to application/json
      const config = {
        headers: {
          "Content-Type": "application/json",
        },
      };
      // Make a GET request to localhost:3000 with the search query
      console.log(filters)
      console.log("query", query)
      axios.get(
        `http://localhost:3000/api/search/filters?query=${query}`,
        {params:filters},
        config
      ).then(response => {
        console.log(response.data)
        if (response.data.results){
          dispatch(updateResults([...results, ...response.data.results]));
          setPage(pageNum + 1);
          setIsFetching(false)
        }
        else{
          console.error('Results field is undefined in the response.');
          dispatch(updateResults([]));
          setPage(1);
          setIsFetching(false);
        }
      })
    } catch (error) {
      console.error("Error making text query:", error);
      dispatch(updateResults([]));
    }
  };

  const handleScroll = () => {
    const threshold = 100;
    const scrolled = window.scrollY;
    const windowHeight = window.innerHeight;
    const documentHeight = document.documentElement.scrollHeight;

    if (documentHeight - (scrolled + windowHeight) < threshold) {
      makeTextQuery(searchQuery, page);
    }
  };

  useEffect(() => {
    window.addEventListener('scroll', handleScroll);

    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, [page, searchQuery]);

  return (
    <div className="visualization-container">
      {
        results.length ===0 ? isFetching ? <div className = "loading"> Loading... </div> : <div className = "no-records-found">No Records Found</div> : 
         <ul>
          {
            results.map((result, index) => (
              <li key={index}>{JSON.stringify(result)}</li>
            ))
          }
        </ul>
      }
    </div>
  );
};

export default Visualization;
