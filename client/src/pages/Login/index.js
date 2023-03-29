//import { Input,Space,Button } from 'antd';
//import { useState } from 'react';

import QRCode from './QRCode';

import './index.css';

export default function Login({setUserID}){
  //const [value,setValue]=useState("");
  /*const onLogin=()=>{
    localStorage.setItem("userID",value);
    setUserID(value);
  }*/

  return (
    <div className="login-page">
      <div className="header">ChatGPT Proxy</div>
      <div className='content'>
        <QRCode setUserID={setUserID}/>
      </div>
    </div>
  );
}