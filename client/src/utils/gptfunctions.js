import axios from "axios";

var g_userid="";

const getUserID=()=>{
  return g_userid;
}

const setUserID=(userid)=>{
  g_userid=userid;
}

//chatproxy
//const chatProxyApi=process.env.REACT_APP_OPENAI_HOST+"/openaiproxy/openai/v1/chat/completions/GPT3Dot5Turbo";
const chatProxyApi=process.env.REACT_APP_OPENAI_HOST+"/openaiproxy/ChatCompletion";
console.log(chatProxyApi)
const chatCompleteProxy=async (messages)=>{
  
  const reponse= await axios({
    url:chatProxyApi,
    timeout:300000,
    method:"post",
    headers: {
      'Content-Type': 'application/json'
    },
    data:{
      sessionid:g_userid,
      messages:messages
    }});

  if(reponse.data?.error===true){
    return reponse.data?.message;
  }
  
  return reponse.data?.result;
}

export {
  chatCompleteProxy,
  getUserID,
  setUserID
}