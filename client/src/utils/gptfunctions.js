import axios from "axios";

//chatproxy
//const chatProxyApi=process.env.REACT_APP_OPENAI_HOST+"/openaiproxy/openai/v1/chat/completions/GPT3Dot5Turbo";
const chatProxyApi=process.env.REACT_APP_OPENAI_HOST+"/openaiproxy/ChatCompletion";
console.log(chatProxyApi)
const chatCompleteProxy=async (userID,messages)=>{
  try {
    const reponse= await axios({
      url:chatProxyApi,
      timeout:3000000,
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
  } catch (error) {
    if (error.response) {
      // 请求已发送，但状态码不在 2xx 范围内
      console.log(error.response.data);
      console.log(error.response.status);
      console.log(error.response.headers);
      return "服务器返回错误："+error.response.status+"，请稍后再试";
    } else if (error.request) {
      // 请求已发送但没有收到响应
      console.log(error.request);
      return "未收到服务器的响应，请检查网络是否正常后再试";
    } else {
      // 发生了一些意外错误
      console.log('Error', error.message);
      return "发送请求时发生错误："+error.message;
    }
  }
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