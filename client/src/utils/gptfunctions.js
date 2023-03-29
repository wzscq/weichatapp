import axios from "axios";

//chatproxy
//const chatProxyApi=process.env.REACT_APP_OPENAI_HOST+"/openaiproxy/openai/v1/chat/completions/GPT3Dot5Turbo";
const chatProxyApi=process.env.REACT_APP_OPENAI_HOST+"/openaiproxy/ChatCompletion";
console.log(chatProxyApi)
const chatCompleteProxy=async (userID,messages)=>{
  
  const reponse= await axios({
    url:chatProxyApi,
    timeout:300000,
    method:"post",
    headers: {
      'Content-Type': 'application/json'
    },
    data:{
      sessionid:userID,
      messages:messages
    }});

  if(reponse.data?.error===true){
    return reponse.data?.message;
  }
  
  return reponse.data?.result;
}

const getTicketUrl = process.env.REACT_APP_OPENAI_HOST+"/public/getTicket";

const scene_id=787;

const getTicketRequest = {
  "expire_seconds": 600, 
  "action_name": "QR_SCENE", 
  "action_info": 
  {"scene": {"scene_id": scene_id}}
}

const getTicket = async ()=>{
  const response = await axios({
    url:getTicketUrl,
    timeout:300000,
    method:"post",
    headers: {
      'Content-Type': 'application/json'
    },
    data:getTicketRequest
  });
  return response.data;
}

export {
  chatCompleteProxy,
  getTicket
}