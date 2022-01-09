import {VscTable} from 'react-icons/vsc';
import {FcIdea} from 'react-icons/fc';
import {BiCube} from 'react-icons/bi';

const getIcon = (type, size = "5xl") => {
    switch (type) {
        case 'featureset':
            return <FcIdea className={`text-${size}`} />;
        case 'table':
            return <VscTable className={`text-${size}`} />;
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
