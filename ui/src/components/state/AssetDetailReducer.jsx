import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialAssetDetailState = {
    asset : {},
    loading : false,
    errorMessage : ""
}

const AssetDetailReducer = (state = initialAssetDetailState, {type, payload}) => {
    switch (type) {
        case 'assetdetail/get':
            getAsset(payload)
            return {
                ...state,
                loading : true
            }
        case 'assetdetail/show':
            return {
                ...state,
                loading : false,
                asset : payload,
                errorMessage : ""
            }
        case 'assetdetail/error':
            return {
                ...state, 
                loading : false,
                errorMessage: payload.statusText
            }
        default:
            return state
    }
}

const getAsset = async (query) => {
    try {
        const request =  {
            url : `${getSvcHost('catalogue')}/asset/name/${query}`,
            options : null
        }
        const response = await fetch(request.url, request.options)
        const data = await response.json()
        if(response.ok){
            store.dispatch({type: "assetdetail/show", payload: data})
        }else{
            store.dispatch({type: "assetdetail/error", payload: {status:response.status, statusText: `${response.statusText}: ${data.message}`}})    
        }
    }catch(error){
        store.dispatch({type: "assetdetail/error", payload: {status: 500, statusText: error.message}})
    }
}


export default AssetDetailReducer