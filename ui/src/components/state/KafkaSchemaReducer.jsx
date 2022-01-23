import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialKafkaSchemaState = {    
    allResults : {},
    schemas : {},
    sortedSchemaNames : [],

    loading : false,
    errorMessage : "",

    limit : 2,
    page : 1,
    pagination : null
}

const selectSchemas = (allResults, keys) => {
    const schemas = {}
    keys.forEach(key => {
        schemas[key] = allResults[key]
    })
    return schemas
}

const updatePagination = (state, targetPage, targetLimit) => {
    const itemNo = state.sortedSchemaNames ? state.sortedSchemaNames.length : 0
    return {
        total : itemNo,
        page : targetPage,
        perPage : targetLimit,
        prev : targetPage-1,
        next : targetPage+1,
        totalPage : Math.ceil(itemNo / targetLimit) 
    }
}

const KafkaSchemaReducer = (state = initialKafkaSchemaState, {type, payload}) => {
    switch (type) {
        case 'kafkaschema/get':  
            getKafkaSchema(payload)
            return {
                ...state,
                loading : true
            }
        case 'kafkaschema/gotopage':        
            const gotoPagination = updatePagination(state, payload, state.limit)
            const gotoSelectedKeys = state.sortedSchemaNames.slice((gotoPagination.page-1)*state.limit, gotoPagination.page*state.limit)
            return {...state, page : payload, 
                pagination : gotoPagination,
                schemas : selectSchemas(state.allResults, gotoSelectedKeys) 
            }
        case 'kafkaschema/resizemaxitems':
            const resizePagination = updatePagination(state, 1, payload)
            const resizeSelectedKeys = state.sortedSchemaNames.slice(
                (resizePagination.page-1)*payload, 
                Math.min(Object.keys(state.allResults).length, resizePagination.page*payload)
            )            
            return {...state, page : 1, limit : payload,
                pagination : resizePagination,
                schemas : selectSchemas(state.allResults, resizeSelectedKeys)
            }
        case 'kafkaschema/show':
            const sortedSchemaNames = Object.keys(payload).sort()
            const itemNo = payload ? sortedSchemaNames.length : 0
            return {...state, 
                loading : false,
                errorMessage : "", 
                allResults : payload,
                schemas : selectSchemas(payload, sortedSchemaNames.slice(0, Math.max(itemNo, state.limit))),
                sortedSchemaNames : sortedSchemaNames,
                pagination : {
                    total : itemNo,
                    page : state.page,
                    perPage : state.limit,
                    prev : state.page-1,
                    next : itemNo > state.limit ? state.page+1 : state.page,
                    totalPage : Math.ceil(itemNo / state.limit) 
                }
            }
        case 'kafkaschema/error':
            return {...state, loading : false, errorMessage: payload.statusText}
        default:
            return state
    }
}


const getKafkaSchema = async (schemaId) => {
    try {
        if(schemaId){
            const rsp = await fetch(`${getSvcHost('kafka_schema_registry')}/subjects/${schemaId}/versions/latest`)
            const rspData = await rsp.json()
            if(rsp.ok){
                const tmp = {}
                tmp[schemaId] = rspData
                store.dispatch({type: "kafkaschema/show", payload: tmp}) 
            }else{
                store.dispatch({type: "kafkaschema/error", payload: {status: rsp.status, statusText: `${rsp.statusText}: ${rspData.message}`}})    
            }
        }else{
            const subjectsResponse = await fetch(`${getSvcHost('kafka_schema_registry')}/subjects`)
            const subjects = await subjectsResponse.json()
            if(!subjectsResponse.ok){
                store.dispatch({type: "kafkaschema/error", payload: {status:subjectsResponse.status, statusText: `${subjectsResponse.statusText}: ${subjects.message}`}})
            }else{
                let schemas = await Promise.all(
                    subjects.map(async schemaName => {
                        let tmpResponse = await fetch(`${getSvcHost('kafka_schema_registry')}/subjects/${schemaName}/versions/latest`)
                        if(tmpResponse.ok){
                            return tmpResponse.json()
                        }
                    })
                )
                // convert to object of objects and dispatch
                store.dispatch({type: 'kafkaschema/show', payload: schemas.reduce((acc, v) => ({ ...acc, [v.subject]: v}), {})})
            }
            
        }
    }catch(error){
        store.dispatch({type: "kafkaschema/error", payload: {status: 500, statusText: error.message}})
    }
}


export default KafkaSchemaReducer