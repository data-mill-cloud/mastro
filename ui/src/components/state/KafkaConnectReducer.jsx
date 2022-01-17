import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialKafkaConnectState = {    
    allResults : null,
    connectors : null,
    sortedConnectorNames : [],

    loading : false,
    errorMessage : "",

    limit : 4,
    page : 1,
    pagination : null
}

const selectConnectors = (allConnectors, keys) => {
    const connectors = {}
    keys.forEach(key => {
        connectors[key] = allConnectors[key]
    })
    return connectors
}

const updatePagination = (state, targetPage, targetLimit) => {
    const connectorNumber = state.sortedConnectorNames ? state.sortedConnectorNames.length : 0
    return {
        total : connectorNumber,
        page : targetPage,
        perPage : targetLimit,
        prev : targetPage-1,
        next : targetPage+1,
        totalPage : Math.ceil(connectorNumber / targetLimit) 
    }
}

const KafkaConnectReducer = (state = initialKafkaConnectState, {type, payload}) => {
    switch (type) {
        case 'kafkaconnect/get':  
            getKafkaConnector(payload)
            return {
                ...state,
                loading : true
            }
        case 'kafkaconnect/gotopage':        
            const gotoPagination = updatePagination(state, payload, state.limit)
            const gotoSelectedKeys = state.sortedConnectorNames.slice((gotoPagination.page-1)*state.limit, gotoPagination.page*state.limit)
            return {...state, page : payload, 
                pagination : gotoPagination,
                connectors : selectConnectors(state.allResults, gotoSelectedKeys) 
            }
        case 'kafkaconnect/resizemaxitems':
            const resizePagination = updatePagination(state, 1, payload)
            const resizeSelectedKeys = state.sortedConnectorNames.slice((resizePagination.page-1)*payload, resizePagination.page*payload)
            return {...state, page : 1, limit : payload,
                pagination : resizePagination,
                connectors : selectConnectors(state.allResults, resizeSelectedKeys)
            }
        case 'kafkaconnect/show':
            const sortedConnectorNames = Object.keys(payload).sort()
            const connectorNumber = payload ? sortedConnectorNames.length : 0
            return {...state, 
                loading : false,
                errorMessage : "", 
                allResults : payload,
                sortedConnectorNames : sortedConnectorNames,
                connectors : selectConnectors(payload, sortedConnectorNames.slice(0, state.limit)            ),
                pagination : {
                    total : connectorNumber,
                    page : state.page,
                    perPage : state.limit,
                    prev : state.page-1,
                    next : connectorNumber > state.limit ? state.page+1 : state.page,
                    totalPage : Math.ceil(connectorNumber / state.limit) 
                }
            }
        case 'kafkaconnect/error':
            return {...state, loading : false, errorMessage: payload.statusText}
        case 'kafkaconnect/restart':
            restartConnector(payload)
            return {...state}
        case 'kafkaconnect/restartsuccess':
            alert(payload.statusText)
            return {...state }
        case 'kafkaconnect/restarterror':
            alert(payload.statusText)
            return {...state }
        default:
            return state
    }
}

const getKafkaConnector = async (connectorId) => {
    try {
        const request = {
            method : 'GET',
            url : `${getSvcHost('kafka_connect')}/connectors?expand=status&expand=info`
        }
        const response = await fetch(request.url, request.options)
        const data = await response.json()
        if(response.ok){
            if(Object.keys(data).length > 0){
                if(connectorId){
                    if(connectorId in Object.keys(data)){
                        const tmp = {}
                        tmp[connectorId] = data[connectorId]
                        store.dispatch({type: "kafkaconnect/show", payload: tmp})
                    }else{
                        store.dispatch({type: "kafkaconnect/error", payload: {statusText: `no connector named ${connectorId} found`}})
                    }
                }else{
                    store.dispatch({type: "kafkaconnect/show", payload:  data})
                }
            } else {
                store.dispatch({type: "kafkaconnect/error", payload: {statusText: `no data returned from service at ${getSvcHost('kafka_connect')}`}})
            }
        }else{
            store.dispatch({type: "kafkaconnect/error", payload: {status:response.status, statusText: `${response.statusText}: ${data.message}`}})    
        }
    }catch(error){
        store.dispatch({type: "kafkaconnect/error", payload: {status: 500, statusText: error.message}})
    }
}

const restartConnector = async (connectorID) => {
    try {
        const request = {
            method : 'POST',
            url : `${getSvcHost('kafka_connect')}/connectors/${connectorID}/restart?includeTasks=true&onlyFailed=true`
        }
        const response = await fetch(request.url, request.options)
        if(response.ok){
            store.dispatch({type: "kafkaconnect/restartsuccess", payload: {status:response.status, statusText: `${response.statusText}`}})    
        }else{
            store.dispatch({type: "kafkaconnect/restarterror", payload: {status:response.status, statusText: `${response.statusText}`}})    
        }
    }catch(error){
        store.dispatch({type: "kafkaconnect/restarterror", payload: {statusText: error.message}})
    }
}

export default KafkaConnectReducer