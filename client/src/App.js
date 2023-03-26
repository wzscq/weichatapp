import {useState} from "react";
import Main from "./pages/Main";
import Login from "./pages/Login";

function App() {
  const [userID, setUserID] = useState(localStorage.getItem("userID") || "");
  return userID!==""?<Main userID={userID} setUserID={setUserID}/>:<Login setUserID={setUserID}/>;
}

export default App;
