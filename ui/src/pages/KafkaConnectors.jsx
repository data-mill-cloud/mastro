import { useEffect } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom'
import {BiError} from 'react-icons/bi'
import Connector from '../components/kafkaconnectors/Connector'

//import {Link} from 'react-router-dom';
import Slider from '../components/pagination/Slider'
import PageSelector from '../components/pagination/PageSelector'

function KafkaConnectors() {
    const { v4: uuidv4 } = require('uuid');
    const params = useParams()
    const dispatch = useDispatch()
    const kafkaconnectState = useSelector(state => state.kafkaconnectState)
    const connectors = kafkaconnectState.connectors
    
    const loading = kafkaconnectState.loading
    const errorMessage = kafkaconnectState.errorMessage

    const pagination = kafkaconnectState.pagination
    const limit = kafkaconnectState.limit

    useEffect(() => {
        dispatch({type: 'kafkaconnect/get', payload: params.assetid ? params.assetid : null})
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
        }else if(connectors && Object.keys(connectors).length > 0){
            return (
                <div>
                    <Slider pagination={pagination} limit={limit} resizeTarget={'kafkaconnect/resizemaxitems'} />
                    <div className="grid grid-cols-1 gap-8 xl:grid-cols-2 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                        {Object.keys(connectors).map(connector => (
                            <Connector key={uuidv4()} name={connector} connector={connectors[connector]}/>
                        ))}
                    </div>
                    <PageSelector pagination={pagination} pageTarget={'kafkaconnect/gotopage'}/>
                </div>
            )
        }else{
            return (
                <div className="alert alert-info">
                    <div className="flex-1">
                        <BiError className="text-2xl" />
                        <label>No connectors found</label>
                    </div>
                </div>
            )
        }
    }else{
        return <h1>Loading...</h1>
    }
}

export default KafkaConnectors
