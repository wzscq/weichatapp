import { Card } from 'antd';

import './index.css';

export default function TopTip(){
    return (
        <Card className='top-tip'>
            <div className='header1'>你好</div>
            <div className='header1'>我是VIPGPT，你的专属AI集合平台</div>
            <div className='header2'>作为一个人工智能语言模型，我可以回答你的问题，为你提供有用信息，帮助你完成创作。</div>
        </Card>
    )
}