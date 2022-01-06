import PropTypes from 'prop-types';
import {Link} from 'react-router-dom';
import {AiOutlineTags} from 'react-icons/ai';
import {FaHashtag} from 'react-icons/fa';
import EntryIcon from './EntryIcon';

function Entry({entry}) {
    return (
        <div data-tip={`last-discovered: ${entry["last-discovered-at"]}`} className="tooltip">
            <div className="card shadow-2xl compact side bg-base-100 border-2 hover:border-gray-400">
                <div className="flex-row items-start space-x-4 card-body">
                    <div>
                        {/*
                        <div className="avatar">
                            <div className="rounded-full shadow w-14 h-14">
                                <img src="https://picsum.photos/200/200" alt="asset" />
                            </div>
                        </div>
                        */}
                        
                        <EntryIcon type={entry.type} />
                        
                    </div>
                    <div>
                        <Link className="card-title" to={`/entries/${entry.name}`} >{entry.name}</Link>
                        <p className="text-base-content text-opacity-40">{entry["last-discovered-at"]}</p>
                        {/*<Link className="text-base-content text-opacity-40" to={`/entries/${entry.name}`} >Open</Link>*/}
                    </div>
                </div>
                <div className="flex-row" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                    <p style={{textAlign: 'justify'}}>{entry.description}</p>
                </div>
                <div className="py-4 artboard artboard-demo bg-base-200">
                    <ul className="menu items-stretch px-3 bg-base-100 horizontal rounded-box">
                        <li data-tip="aaa" className="tooltip tooltip-open tooltip-bottom"><a><FaHashtag className="text-base-content text-xl" />Tags</a></li> 
                        <li data-tip="aaa" className="tooltip tooltip tooltip-open tooltip-bottom"><a><AiOutlineTags className="text-2xl" />Labels</a></li> 
                    </ul>
                </div>
            </div>
        </div>
    )
}

Entry.propTypes = {
    entry : PropTypes.object.isRequired,
}

export default Entry
