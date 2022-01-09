import { useEffect } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom'
import AssetIcon from '../components/entries/AssetIcon'
import {FaLink} from 'react-icons/fa'
import {BiError} from 'react-icons/bi'
import LineageChart from '../components/lineage/LineageChart';

function Asset() {
    const params = useParams()
    const dispatch = useDispatch()
    const searchState = useSelector(state => state.searchState)
    const asset = searchState.asset
    const loading = searchState.loading
    const errorMessage = searchState.errorMessage

    useEffect(() => {
        // get asset info
        // params.assetid
        dispatch({type: 'assetdetail/get', payload: params.assetid})
    }, [])

    const getUpwardLineage = (asset) => {
        var lineageData = {name : asset.name, children : []}
        if (asset['depends-on']){
            lineageData.children = asset['depends-on'].map(parent => {
                return {name : parent, children : []}
            })
        }
        return lineageData
    }
    
    const getDownwardLineage = (asset) => {
        var lineageData = {name : "root", children : []}

        if (asset['depends-on']){
            lineageData.children = asset['depends-on'].map(parentDependency => {
                return {
                    name : parentDependency,
                    children : [
                        { name : asset.name, children : [] }
                    ]
                }
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
                                    <div className="ml-2 mr-1 badge badge-success">{asset.type}</div>
                                </h1>
                                <p>{asset.description}</p>
                            </div>
                            <div className="w-full rounded-lg shadow-md bg-base-100 stats">
                                <div className="stat">
                                    <div className="stat-title text-md">Labels</div>
                                    <div className="stat-value text-lg">
                                        <div className="overflow-x-auto">
                                            <table className="table w-full">
                                                <thead>
                                                <tr>
                                                    <th>Key</th> 
                                                    <th>Value</th> 
                                                </tr>
                                                </thead> 
                                                <tbody>
                                                    { asset.labels !== undefined ? Object.keys(asset.labels).map(key => (<tr><th>{key}</th><td>{asset.labels[key]}</td></tr>)): '' }
                                                </tbody>
                                            </table>
                                        </div>
                                    </div> 
                                </div>
                            
                                <div className="stat">
                                    <div className="stat-title text-md">Tags</div>
                                    <div className="stat-value text-lg">
                                        { asset.tags !== undefined ? asset.tags.map(tag => <span className="badge badge-primary mr-1">{tag}</span>) : ''}
                                    </div> 
                                </div>
                                
                                <div className="stat">
                                    <div className="stat-title text-md">Versions</div>
                                    <div className="stat-value text-lg">
                                        { asset.versions !== undefined ? Object.keys(asset.versions) : ''}
                                    </div> 
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="w-full rounded-lg shadow-md bg-base-100 stats py-5 mb-6">
                        <div className="stat">
                            <div className="stat-figure text-secondary">
                                <FaLink className="text-3xl md:text-5xl" />
                            </div>
                            <div className="stat-title text-md">Depends on</div>
                            <LineageChart lineageData={getUpwardLineage(asset)} width={1024} height={400}/>    
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
