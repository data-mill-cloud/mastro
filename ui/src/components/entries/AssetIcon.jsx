import {VscTable} from 'react-icons/vsc';
import {FcIdea} from 'react-icons/fc';
import {BiCube} from 'react-icons/bi';
import {FaStream} from 'react-icons/fa';
import {ImDatabase} from 'react-icons/im';
import {BsSpeedometer2} from 'react-icons/bs';

const getIcon = (type, size = "5xl") => {
    switch (type) {
        case 'metricset':
            return <BsSpeedometer2 className={`text-${size}`} />;
        case 'featureset':
            return <FcIdea className={`text-${size}`} />;
        case 'stream':
            return <FaStream className={`text-${size}`} />;
        case 'table':
            return <VscTable className={`text-${size}`} />;
        case 'db':
            return <ImDatabase className={`text-${size}`} />;
        default:
            return <BiCube className={`text-${size}`} />;
    }    
}

function AssetIcon({type, size}) {
    return (
        <div className="content-center">
            {getIcon(type, size)}
        </div>
    )
}

export default AssetIcon
