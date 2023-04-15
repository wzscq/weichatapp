import { Card } from 'antd';

import './index.css';

export default function BottomTip(){
    return (
        <div className='bottom-tip'>
            <div className='title'>VIPGPT</div>
            <Card >
                <div className='header2'>请输入问题 ...</div>
            </Card>
        </div>
    );
}