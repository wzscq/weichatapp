import QRCode from './QRCode';
import PageHeader from '../../components/PageHeader';
import TopTip from './TopTip';
import BottomTip from './BottomTip';
import BottomBanner from '../BottomBanner';

import './index.css';

export default function Login(){
  return (
    <div className="login-page">
      <PageHeader current="/VIPGPTLogin"/>
      <div className='content'>
        <TopTip/>
        <QRCode/>
        <BottomTip/>
        <BottomBanner/>
      </div>
    </div>
  );
}