import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialSearchState = {
    asset : {},
    assets : [],
    loading : false,
    errorMessage : ""
}

const SearchReducer = (state = initialSearchState, {type, payload}) => {
    switch (type) {
        case 'search/submit':
            search(payload)
            return {
                ...state,
                loading : true,
                errorMessage : ""
            }
        case 'search/fetched':
            return {
                ...state, 
                assets : payload,
                loading : false,
                errorMessage : ""
            }
        case 'search/error':
            return {
                ...state,
                assets : [],
                loading : false,
                errorMessage : payload.statusText
            }
        case 'search/clear':
            return {
                ...state,
                assets: [],
                loading : false,
                errorMessage: ""
            }
        default:
            return state
    }
}

const getRequest = (query) => {
    var elements = query.split(",");
    // get by name is 1 element only without #
    if(elements.length === 1 && !elements[0].includes("#")) {
        return {
            url : `${getSvcHost('catalogue')}/asset/name/${elements[0]}`,
            options : null
        }
    }else{
        // we either have a list of tags (>1) or whatever having # inside
        for (var i = 0; i < elements.length; i++) {
            elements[i] = elements[i].trim().replace("#", "");
        }
        return {
            url : `${getSvcHost('catalogue')}/assets/tags`,
            options : { 
                method : 'POST',
                headers : { 'Content-Type': 'application/json' },
                body : JSON.stringify({ tags: elements })
            }
        }
    }
}

const search = async (query) => {
    try {
        const request = getRequest(query)
        const response = await fetch(request.url, request.options)
        const data = await response.json()
        if(response.ok){
            store.dispatch({type: "search/fetched", payload: data.constructor !== Array ? [data] : data})
        }/*else if(response.status === 404){
            window.location = "/notfound"
        }*/else{
            store.dispatch({type: "search/error", payload: {status:response.status, statusText: `${response.statusText}: ${data.message}`}})    
        }
    }catch(error){
        store.dispatch({type: "search/error", payload: {status: 500, statusText: error.message}})
    }
}


export default SearchReducer