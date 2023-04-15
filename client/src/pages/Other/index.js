import { useLocation } from 'react-router-dom';
import PageHeader from "../../components/PageHeader";

import "./index.css";

export default function Other() {
    const location = useLocation();
    
    return (
        <div className="other-page">
            <PageHeader current={location.pathname}/>
            <div className='content'>
                功能开发中 ...
            </div>
        </div>  
    )
}