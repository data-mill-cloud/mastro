import { useEffect } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom'
import {BiError} from 'react-icons/bi'
import Featureset from '../components/featuresets/Featureset'
import {Link} from 'react-router-dom';
import Slider from '../components/pagination/Slider'
import PageSelector from '../components/pagination/PageSelector'
import FeatureViewer from '../components/featuresets/FeatureViewer';

function Featuresets() {
    const { v4: uuidv4 } = require('uuid');
    const params = useParams()
    const dispatch = useDispatch()
    const featuresetState = useSelector(state => state.featuresetState)
    const selectedFeatureset = featuresetState.selectedFeatureset
    const featuresets = featuresetState.featuresets
    const loading = featuresetState.loading
    const errorMessage = featuresetState.errorMessage

    const pagination = featuresetState.pagination
    const limit = featuresetState.limit

    useEffect(() => {
        dispatch({type: 'featureset/get', payload: params.assetid})
    },[dispatch, params.assetid])

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
                    <div className="mb-6">
                        <h1 className="text-3xl card-title">
                            {params.assetid}
                            <Link className="card-title" to={`/asset/${params.assetid}`}>
                                <div className="ml-2 mr-1 badge badge-primary">asset</div>
                            </Link>
                            <Link className="card-title" to={`/featureset/${params.assetid}`}>
                                <div className="ml-2 mr-1 badge badge-primary">featureset</div>
                            </Link>
                        </h1>
                    </div>
                    <Slider pagination={pagination} limit={limit} resizeTarget={'featureset/resizemaxitems'} />
                    <div className="grid grid-cols-1 gap-8 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                        {featuresets.map(featureset => <Featureset key={uuidv4()} featureset={featureset} />)}    
                    </div>
                    { selectedFeatureset && (
                        <div className="mt-6">
                            <FeatureViewer featureset={selectedFeatureset} />
                        </div>
                    )}
                    <PageSelector pagination={pagination} pageTarget={'featureset/gotopage'}/>         
                </div>
            )
        }
    }else{
        return <h1>Loading...</h1>
    }
}

export default Featuresets
