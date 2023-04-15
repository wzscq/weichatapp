import { Menu } from 'antd';
import { useNavigate } from 'react-router-dom';

import logo from '../../images/logo/4.png';

import "./index.css";

const items = [
    {
      label: 'VIPGPT',
      key: '/VIPGPTLogin',
    },
    {
      label: 'Midjourney',
      key: '/Midjourney'
    },
    {
        label: '文心一言',
        key: '/Wenxin'
    },
    {
        label: '腾讯智影',
        key: '/Tencent'
    },
    {
        label: '华为盘古',
        key: '/Huawei'
    },
]

export default function PageHeader({current}) {
    const navigate = useNavigate();

    const onClick = (e) => {
        console.log('click ', e);
        navigate(e.key);
    };

    return (
        <div className="page-header">
            <img src={logo} alt='logo' className='logo' />
            <Menu onClick={onClick} selectedKeys={[current]} mode="horizontal" items={items} />
        </div>
    );
}