import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialMetricsetState = {
    selectedMetricset : null,
    metricsets : [],
    loading : false,
    errorMessage : ""
}

const MetricsetReducer = (state = initialMetricsetState, {type, payload}) => {
    switch (type) {
        case 'metricset/get':
            getMetricset(payload)
            return {
                ...state,
                loading : true
            }
        case 'metricset/show':
            return {
                ...state,
                loading : false,
                metricsets : payload,
                errorMessage : ""
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


const getMetricset = async (assetId) => {
    try {
        const request = {
            method : 'GET',
            url: `${getSvcHost('metricstore')}/metricstore/name/${assetId}`
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