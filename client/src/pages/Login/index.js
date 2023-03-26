import { Input,Space,Button } from 'antd';
import { useState } from 'react';

import './index.css';

export default function Login({setUserID}){
  const [value,setValue]=useState("");

  const onLogin=()=>{
    localStorage.setItem("userID",value);
    setUserID(value);
  }

  return (
    <div className="login-page">
      <div className="header">ChatGPT Proxy</div>
      <div className='content'>
      <Space.Compact
        style={{
          width: '100%',
        }}
      >
        <Input value={value} onChange={(e)=>setValue(e.target.value)} placeholder='请输入用户标识并点击右侧按钮进入聊天'/>
        <Button type="primary" onClick={onLogin}>Login</Button>
      </Space.Compact>
      </div>
    </div>
  );
}