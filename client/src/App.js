import {
  createHashRouter,
  RouterProvider,
} from "react-router-dom";

import VIPGPTMain from "./pages/VIPGPTMain";
import VIPGPTLogin from "./pages/VIPGPTLogin";
import Other from "./pages/Other";


const router = createHashRouter([
  {
    path: "/",
    element: <VIPGPTLogin/>,
  },
  {
    path: "/VIPGPTLogin",
    element: <VIPGPTLogin/>,
  },
  {
    path: "/VIPGPTMain",
    element: <VIPGPTMain/>,
  },
  {
    path: "/Midjourney",
    element: <Other/>,
  },
  {
    path: "/Wenxin",
    element: <Other/>,
  },
  {
    path: "/Tencent",
    element: <Other/>,
  },
  {
    path: "/Huawei",
    element: <Other/>,
  }
]);

function App() {
  return (<RouterProvider router={router} />);
}

export default App;
