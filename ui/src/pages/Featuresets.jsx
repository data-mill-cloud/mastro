import { useEffect } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom'
import {BiError} from 'react-icons/bi'
import Featureset from '../components/featuresets/Featureset'
import {Link} from 'react-router-dom';

function Featuresets() {
    const { v4: uuidv4 } = require('uuid');
    const params = useParams()
    const dispatch = useDispatch()
    const featuresetState = useSelector(state => state.featuresetState)
    const featuresets = featuresetState.featuresets
    const loading = featuresetState.loading
    const errorMessage = featuresetState.errorMessage

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
                    <div className="grid grid-cols-1 gap-8 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                        {featuresets.map(featureset => <Featureset key={uuidv4()} featureset={featureset} />)}    
                    </div>
                </div>
            )
        }
    }else{
        return <h1>Loading...</h1>
    }
}

export default Featuresets
