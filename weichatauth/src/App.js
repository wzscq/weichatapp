import React, { useEffect, useState } from 'react';
import axios from "axios";

import './App.css';

const params = new URLSearchParams(window.location.search);

function App() {
  const [loaded,setLoaded]=useState(false);
  const [loading,setLoading]=useState(true);
  const [customerInfo,setCustomerInfo]=useState({});

  useEffect(()=>{
    if(loaded===false){
      const code=params.get("code");
      if(code){
        setLoaded(true);
        const url=process.env.REACT_APP_HOST+"/customer/getCustomerInfo";
        const reqConfig={
          url:url,
          method:"post",
          data:{
            code:code
          }
        };
        axios(reqConfig).then((response)=>{
          setLoading(false);
          setCustomerInfo(response.data);
        });
      }
    }
  },[setLoaded,loaded,setLoading,setCustomerInfo]);

  let content=(<div>loading</div>);

  if(loading===false){
    if(customerInfo?.nickname||customerInfo?.headimgurl){
      content=(
        <div>
          <div style={{marginTop:20}}>欢迎</div>
          <div>{customerInfo?.nickname}</div>
          <div><img style={{width:100,height:100,borderRadius:5,marginTop:10,marginBottom:10}} src={customerInfo?.headimgurl} alt="" /></div>
          <div>授权成功，赶快打开页面体验你的私人助理</div>
          <div>web地址:http://vipgpt.top/client</div>
        </div>);
    } else {
      content=(
        <div style={{marginTop:20}}>授权失败，请稍后重新尝试</div>
      );
    }
  }

  return (
    <div className="App">
      <div style={{display:"none"}}>{params.get("code")}</div>
      {content}
    </div>
  );
}

export default App;
