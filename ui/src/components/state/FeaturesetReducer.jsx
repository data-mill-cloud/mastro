import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialFeaturesetState = {
    selectedFeatureset : null,

    featuresets : [],
    loading : false,
    errorMessage : "",

    limit : 4,
    page : 0,
    pagination : null
}

const FeaturesetReducer = (state = initialFeaturesetState, {type, payload}) => {
    switch (type) {
        case 'featureset/get':
            getFeatureset(payload, state.limit, state.page)
            return {
                ...state,
                query : payload,
                loading : true,
                selectedFeatureset : null
            }
        case 'featureset/gotopage':        
            getFeatureset(state.query, state.limit, payload)
            return {...state, page : payload, selectedFeatureset : null}
        case 'featureset/resizemaxitems':
            if(state.pagination && state.pagination.total >= state.limit){
                getFeatureset(state.query, payload, state.page)
            }
            return {...state, limit : payload, selectedFeatureset : null}
        case 'featureset/show':
            return {
                ...state,
                loading : false,
                errorMessage : "", 
                page : 1,
                featuresets : 'pagination' in payload ? payload.data : [payload],
                pagination : 'pagination' in payload ? payload.pagination : null,
            }
        case 'featureset/error':
            return {
                ...state, 
                loading : false,
                errorMessage: payload.statusText
            }
        case 'featureset/select':
            return {
                ...state,
                selectedFeatureset : payload
            }
        case 'featureset/unselect':
                return {
                    ...state,
                    selectedFeatureset : null
                }
        default:
            return state
    }
}


const getFeatureset = async (assetId, limit, page) => {
    try {
        const request = {
            method : 'GET',
            //url: `${getSvcHost('featurestore')}/featureset/name/${assetId}`
            url: `${getSvcHost('featurestore')}/featureset/name/${assetId}?limit=${limit}&page=${page}`
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