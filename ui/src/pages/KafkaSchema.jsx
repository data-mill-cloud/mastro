import { useEffect } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom'
import {BiError} from 'react-icons/bi'
import Slider from '../components/pagination/Slider'
import PageSelector from '../components/pagination/PageSelector'
import Schema from '../components/schemaregistry/Schema'

function KafkaSchema() {
    const { v4: uuidv4 } = require('uuid');
    const params = useParams()
    const dispatch = useDispatch()
    const kafkaSchemaState = useSelector(state => state.kafkaSchemaState)
    const schemas = kafkaSchemaState.schemas
    const loading = kafkaSchemaState.loading
    const errorMessage = kafkaSchemaState.errorMessage
    const pagination = kafkaSchemaState.pagination
    const limit = kafkaSchemaState.limit


    useEffect(() => {
        dispatch({type: 'kafkaschema/get', payload: params.assetid ? params.assetid : null})
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
        }else if(schemas && Object.keys(schemas).length > 0){
            return (
                <div className="h-full">
                    <Slider pagination={pagination} limit={limit} resizeTarget={'kafkaschema/resizemaxitems'} />
                    <div className="grid h-full grid-cols-1 gap-8 xl:grid-cols-2 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                        {Object.keys(schemas).map(schema => (
                            <Schema key={uuidv4()} name={schema} schema={schemas[schema]}/>
                        ))}
                    </div>
                    <PageSelector pagination={pagination} pageTarget={'kafkaschema/gotopage'}/>
                </div>
            )
        }else{
            return (
                <div className="alert alert-info">
                    <div className="flex-1">
                        <BiError className="text-2xl" />
                        <label>No schema found</label>
                    </div>
                </div>
            )
        }
    }else{
        return <h1>Loading...</h1>
    }
}

export default KafkaSchema
