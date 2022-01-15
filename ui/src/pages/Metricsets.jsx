import { useEffect } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom'
import {BiError} from 'react-icons/bi'
import {Link} from 'react-router-dom';
import Metricset from '../components/metricsets/Metricset'
import MetricViewer from '../components/metricsets/MetricViewer';
import Slider from '../components/pagination/Slider'
import PageSelector from '../components/pagination/PageSelector'

function Metricsets() {

    const { v4: uuidv4 } = require('uuid');
    const params = useParams()
    const dispatch = useDispatch()
    const metricsetState = useSelector(state => state.metricsetState)
    const metricsets = metricsetState.metricsets
    const selectedMetricset = metricsetState.selectedMetricset
    const loading = metricsetState.loading
    const errorMessage = metricsetState.errorMessage

    const pagination = metricsetState.pagination
    const limit = metricsetState.limit

    useEffect(() => {
        dispatch({type: 'metricset/get', payload: params.assetid})
    }, [dispatch, params.assetid])

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
                            <Link className="card-title" to={`/metricset/${params.assetid}`}>
                                <div className="ml-2 mr-1 badge badge-primary">metricset</div>
                            </Link>
                        </h1>
                    </div>
                    <Slider pagination={pagination} limit={limit} resizeTarget={'metricset/resizemaxitems'} />
                    <div className="grid grid-cols-1 gap-8 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                        {metricsets.map(metricset => <Metricset key={uuidv4()} metricset={metricset} />)}    
                    </div>
                    { selectedMetricset && (
                        <div className="mt-6">
                            <MetricViewer metricset={selectedMetricset} />
                        </div>
                    )}
                    <PageSelector pagination={pagination} pageTarget={'metricset/gotopage'}/>   
                </div>
            )
        }
    }else{
        return <h1>Loading...</h1>
    }
}

export default Metricsets
