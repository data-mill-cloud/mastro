import {VscTable} from 'react-icons/vsc';
import {FcIdea} from 'react-icons/fc';
import {BiCube} from 'react-icons/bi';

const getIcon = (type) => {
    switch (type) {
        case 'featureset':
            return <FcIdea className="text-5xl" />;
        case 'table':
            return <VscTable className="text-5xl" />;
        default:
            return <BiCube className="text-5xl" />;
    }
}

function EntryIcon({type}) {
    return (
        <div className="content-center">
            {getIcon(type)}
        </div>
    )
}

export default EntryIcon
