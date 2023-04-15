import { useEffect, useState } from 'react';
import mqtt from 'mqtt';
import { useNavigate } from "react-router-dom";
import {getTicket} from '../../../utils/gptfunctions';
import Loading from './Loading';

import './index.css';

const initTicket = {
  loaded: false,
  loading: false,
  ticket:{},
  error:false,
  message:""
}

var g_MQTTClient=null;

export default function QRCode(){
  const [ticket,setTicket]=useState(initTicket);
  const navigate = useNavigate();
  
  useEffect(()=>{
    if(ticket.loaded===false&&ticket.loading===false){
      setTicket({...ticket,loading:true});
      getTicket().then((data)=>{
        console.log(data);
        if(data?.error===true){
          setTicket({...ticket,loading:false,loaded:true,error:true,message:data?.message});
          return;
        }
        setTicket({...ticket,loading:false,loaded:true,error:false,ticket:data?.result});
      });
    }
  },[ticket]);

  useEffect(()=>{
    //收到ticket以后，基于sceneid到mqtt做个订阅，接受用户扫码后的消息，可以从中拿到用户ID
    const connectMqtt=(topic)=>{
      console.log("connectMqtt ... ");
      if(g_MQTTClient!==null){
          g_MQTTClient.end();
          g_MQTTClient=null;
      }

      const server='ws://121.37.185.248:9101';
      const options={
          username:'mosquitto',
          password:'123456',
      }
      console.log("connect to mqtt server ... "+server+" with options:",options);
      g_MQTTClient  = mqtt.connect(server,options);
      g_MQTTClient.on('connect', () => {
          console.log("connected to mqtt server "+server+".");
          console.log("subscribe topics :"+topic+"");
          g_MQTTClient.subscribe(topic, (err) => {
              if(!err){
                  console.log("subscribe topics success.");
                  console.log("topic:",topic);
              } else {
                  console.log("subscribe topics error :"+err.toString());
              }
          });
      });
      g_MQTTClient.on('message', (topic, payload, packet) => {
          console.log("receive message topic :"+topic+" content :"+payload.toString());
          //setUserID(payload.toString());
          localStorage.setItem('userID',payload.toString());
          navigate('/VIPGPTMain');
          setTimeout(()=>{
            g_MQTTClient.end();
            g_MQTTClient=null;
          },10);
      });
      g_MQTTClient.on('close', () => {
          console.log("mqtt client is closed.");
      });
    }

    if(ticket?.ticket?.sceneID>0){
      console.log("ticket.sceneID:",ticket?.ticket?.sceneID);
      connectMqtt('qrlogin/'+ticket?.ticket?.sceneID);
    }

  },[ticket,navigate]);

  return (
    <div className='qr-code'>
      {ticket.loading===true?
        <Loading/>:
        <>
        <div className='title'>微信扫码关注公众号登录</div>
        <img alt="" src={"https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket="+ticket.ticket.ticket} />
        </>
      }
    </div>
  )
}