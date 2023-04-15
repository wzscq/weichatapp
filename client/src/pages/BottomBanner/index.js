
import banner1 from '../../images/bottombanner/1.jpg';
import banner2 from '../../images/bottombanner/2.jpg';
import banner3 from '../../images/bottombanner/3.jpg';
import banner4 from '../../images/bottombanner/4.jpg';
import banner5 from '../../images/bottombanner/5.jpg';
import banner6 from '../../images/bottombanner/6.jpg';
import banner7 from '../../images/bottombanner/7.jpg';

import './index.css';

export default function BottomBanner(){
  return (
    <div className="bottom-banner">
      <img src={banner1} alt="banner1"/>
      <img src={banner2} alt="banner2"/>
      <img src={banner3} alt="banner3"/>
      <img src={banner4} alt="banner4"/>
      <img src={banner5} alt="banner5"/>
      <img src={banner6} alt="banner6"/>
      <img src={banner7} alt="banner7"/>
    </div>
  );
}