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

export {
  chatCompleteProxy
}