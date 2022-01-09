import PropTypes from 'prop-types';
import {Link} from 'react-router-dom';
import AssetIcon from './AssetIcon';

function Asset({asset}) {
    return (
            <div className="card shadow-2xl compact side bg-base-100 border-2 hover:border-gray-400">
                    <div className="flex-row items-start space-x-4 card-body">
                    <div>
                        <AssetIcon type={asset.type} />
                    </div>
                    <div>
                        <div className="flex-row">
                            <Link className="card-title" to={`/asset/${asset.name}`}>
                                {asset.name}
                                
                            </Link>
                        </div>
                        <p className="text-base-content text-opacity-40">{asset["last-discovered-at"]}</p>
                    </div>
                </div>
                <div className="card-body flex-row" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                    <p style={{textAlign: 'justify'}}>{asset.description}</p>
                </div>
            </div>
    )
}

Asset.propTypes = {
    asset : PropTypes.object.isRequired,
}

export default Asset
