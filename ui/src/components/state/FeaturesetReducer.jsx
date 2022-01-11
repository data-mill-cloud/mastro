import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialFeaturesetState = {
    featuresets : [],
    loading : false,
    errorMessage : ""
}

const FeaturesetReducer = (state = initialFeaturesetState, {type, payload}) => {
    switch (type) {
        case 'featureset/get':
            getFeatureset(payload)
            return {
                ...state,
                loading : true
            }
        case 'featureset/show':
            return {
                ...state,
                loading : false,
                featuresets : payload,
                errorMessage : ""
            }
        case 'featureset/error':
            return {
                ...state, 
                loading : false,
                errorMessage: payload.statusText
            }
        default:
            return state
    }
}


const getFeatureset = async (assetId) => {
    try {
        const request = {
            method : 'GET',
            url: `${getSvcHost('featurestore')}/featureset/name/${assetId}`
        }
        const response = await fetch(request.url, request.options)
        const data = await response.json()
        if(response.ok){
            store.dispatch({type: "featureset/show", payload: data})
        }else{
            store.dispatch({type: "featureset/error", payload: {status:response.status, statusText: `${response.statusText}: ${data.message}`}})    
        }
    }catch(error){
        store.dispatch({type: "featureset/error", payload: {status: 500, statusText: error.message}})
    }
}


export default FeaturesetReducer