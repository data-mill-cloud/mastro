import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialMetricsetState = {
    selectedMetricset : null,
    metricsets : [],
    loading : false,
    errorMessage : "", 

    limit : 8,
    page : 0,
    pagination : null
}

const MetricsetReducer = (state = initialMetricsetState, {type, payload}) => {
    switch (type) {
        case 'metricset/get':
            getMetricset(payload, state.limit, state.page)
            return {
                ...state,
                query: payload,
                loading : true,
                selectedMetricset : null,
            }
        case 'metricset/gotopage':        
            getMetricset(state.query, state.limit, payload)
            return {...state, page : payload, selectedMetricset : null,}
        case 'metricset/resizemaxitems':
            if(state.pagination && state.pagination.total >= state.limit){
                getMetricset(state.query, payload, state.page)
            }
            return {...state, limit : payload, selectedMetricset : null,}
        case 'metricset/show':
            return {
                ...state,
                loading : false,
                errorMessage : "",
                page : 1,
                metricsets : 'pagination' in payload ? payload.data : [payload],
                pagination : 'pagination' in payload ? payload.pagination : null,
            }
        case 'metricset/error':
            return {
                ...state, 
                loading : false,
                errorMessage: payload.statusText
            }
        case 'metricset/select':
            return {
                ...state,
                selectedMetricset : payload
            }
        case 'metricset/unselect':
                return {
                    ...state,
                    selectedMetricset : null
                }
        default:
            return state
    }
}


const getMetricset = async (assetId, limit, page) => {
    try {
        const request = {
            method : 'GET',
            //url: `${getSvcHost('metricstore')}/metricstore/name/${assetId}`
            url: `${getSvcHost('metricstore')}/metricstore/name/${assetId}?limit=${limit}&page=${page}`
        }
        const response = await fetch(request.url, request.options)
        const data = await response.json()
        if(response.ok){
            store.dispatch({type: "metricset/show", payload: data})
        }else{
            store.dispatch({type: "metricset/error", payload: {status:response.status, statusText: `${response.statusText}: ${data.message}`}})    
        }
    }catch(error){
        store.dispatch({type: "metricset/error", payload: {status: 500, statusText: error.message}})
    }
}


export default MetricsetReducer