import store from "./store"

const initialSearchState = {
    entries : [],
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
                entries : payload,
                loading : false,
                errorMessage : ""
            }
        case 'search/error':
            return {
                ...state,
                entries : [],
                loading : false,
                errorMessage : payload.statusText
            }
        default:
            return state
    }
}

const getSvcHost = (svcName) => {
    const devEnvVar = `REACT_APP_${svcName.toUpperCase()}_URL`
    return typeof process.env[devEnvVar] !== 'undefined' ? process.env[devEnvVar] : svcName
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
        }else{
            store.dispatch({type: "search/error", payload: {status:response.status, statusText: `${response.statusText}: ${data.message}`}})    
        }
    }catch(error){
        store.dispatch({type: "search/error", payload: {status: 500, statusText: error.message}})
    }
}

export default SearchReducer