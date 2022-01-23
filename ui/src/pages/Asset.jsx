import { useEffect } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom'
import AssetIcon from '../components/entries/AssetIcon'
import {BiError} from 'react-icons/bi'
import LineageChart from '../components/lineage/LineageChart';
import {Link} from 'react-router-dom';

function Asset() {
    const { v4: uuidv4 } = require('uuid');
    const params = useParams()
    const dispatch = useDispatch()
    const assetDetailState = useSelector(state => state.assetDetailState)
    const asset = assetDetailState.asset
    const loading = assetDetailState.loading
    const errorMessage = assetDetailState.errorMessage

    useEffect(() => {
        dispatch({type: 'assetdetail/get', payload: params.assetid})
    }, [dispatch, params.assetid])

    const getLineage = (asset) => {
        const currentAssetId = uuidv4()
        const lineageData = [{
            id : currentAssetId, data : { label: asset.name }, position : {x : 400, y : 150},
            //sourcePosition: 'top', targetPosition: 'top'
         }
        ]
        if('depends-on' in asset){
            asset['depends-on'].forEach(function (parent, i) {
                let parentId = uuidv4()
                lineageData.push({
                    id: parentId, data : {label : parent}, type : 'input', position : {x : i*250, y:0}, //sourcePosition: 'bottom'
                })
                lineageData.push({
                    id: `edge-${parentId}-to-${currentAssetId}`,
                    source: parentId,
                    type: 'smoothstep',
                    target: currentAssetId,
                    animated: true,
                })
            })
        }
        return lineageData
    }

    if (!loading){
        if(errorMessage){
            return (
                <div className="alert alert-error">
                    <div className="flex-1">
                        <BiError className="text-2xl" />
                        <label>{errorMessage}</label>
                    </div>
                </div>
            )
        }else{
            return (
                <div>
                    <div className="grid grid-cols-1 xl:grid-cols-6 lf:grid-cols-6 md:grid-cols-6 mb-8 md:gap-8">
                        <div className="custom-card-img">
                            <AssetIcon type={asset.type} size="9xl" />
                        </div>
                        <div className="col-span-5">
                            <div className="mb-6">
                                <h1 className="text-3xl card-title">
                                    {asset.name}
                                    <Link className="card-title" to={`/asset/${asset.name}`}>
                                        <div className="ml-2 mr-1 badge badge-primary">asset</div>
                                    </Link>
                                    <Link className="card-title" to={`/${asset.type}/${asset.name}`}>
                                        <div className="ml-2 mr-1 badge badge-primary">{asset.type}</div>
                                    </Link>
                                </h1>
                                <p>{asset.description}</p>
                            </div>
                            <div className="w-full rounded-lg shadow-md bg-base-100 stats">                       
                            
                                <div className="stat">
                                    <div className="stat-title text-md">Tags</div>
                                    <div className="stat-value text-lg overflow-auto">
                                        { asset.tags !== undefined ? asset.tags.map(tag => <span key={uuidv4()} className="badge badge-primary mr-1">{tag}</span>) : ''}
                                    </div> 
                                </div>
                                
                                <div className="stat">
                                    <div className="stat-title text-md">Versions</div>
                                    <div className="stat-value text-lg overflow-auto">
                                        { asset.versions !== undefined ? Object.keys(asset.versions) : ''}
                                    </div> 
                                </div>
                            </div>
                        </div>
                    </div>

                    { asset.labels && Object.keys(asset.labels).map(key => (
                                    <div key={uuidv4()} className="shadow stats">
                                        <div className="stat">
                                            <div className="stat-title">{key}</div> 
                                            <div className="stat-value">{asset.labels[key]}</div> 
                                        </div>
                                    </div>
                                ))}

                    <div className="w-full rounded-lg shadow-md bg-base-100 py-5 mb-6">        
                        <div className="stat-title ml-5">Depends on</div>
                        <div className="stat h-72">
                            <LineageChart lineageData={getLineage(asset)}/>    
                        </div>
                    </div>

                    
                </div>
            )
        }
    }else{
        return <h1>Loading...</h1>
    }
    
}

export default Asset
