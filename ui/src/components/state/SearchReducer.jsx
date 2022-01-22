import store from "./store"
import {getSvcHost} from "../../SvcUtils"

const initialSearchState = {
    query : null,
    
    limit : 8,
    page : 0,
    pagination : null,

    assets : [],
    loading : false,
    errorMessage : ""
}

const SearchReducer = (state = initialSearchState, {type, payload}) => {
    switch (type) {
        case 'search/submit':
            search(payload, state.limit, state.page)
            return {
                ...state,
                query : payload,
                loading : true,
                errorMessage : "",
                assets : [],
                pagination : null
            }
        case 'search/gotopage':        
            search(state.query, state.limit, payload)
            return {...state, page : payload}
        case 'search/resizemaxitems':
            if(state.pagination && state.pagination.total >= state.limit){
                search(state.query, payload, state.page)
            }
            return {...state, limit : payload}
        case 'search/fetched':
            return {
                ...state, 
                page : 1,
                assets : 'pagination' in payload ? payload.data : [payload],
                pagination : 'pagination' in payload ? payload.pagination : null,
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
            return initialSearchState
        default:
            return state
    }
}

const getRequest = (query, limit, page) => {
    var elements = query.split(",");
    // get by name is 1 element only without #
    if(elements.length === 1 && !elements[0].includes("#") && !elements[0].includes(" ")) {
        return {
            url : `${getSvcHost('catalogue')}/asset/name/${elements[0]}`,
            options : null
        }

    }else if(elements[0].includes("#")){
        // we either have a list of tags (>1) or whatever having # inside
        for (var i = 0; i < elements.length; i++) {
            elements[i] = elements[i].trim().replace("#", "");
        }
        return {
            //url : `${getSvcHost('catalogue')}/assets/tags?limit=${limit}&page=${page}`,
            url : `${getSvcHost('catalogue')}/assets/tags`,
            options : { 
                method : 'POST',
                headers : { 'Content-Type': 'application/json' },
                body : JSON.stringify({ tags: elements, limit : limit, page : page })
            }
        }
    }else{
        return {
            url : `${getSvcHost('catalogue')}/assets/search`,
            options : { 
                method : 'POST',
                headers : { 'Content-Type': 'application/json' },
                body : JSON.stringify({ query: query, limit : limit, page : page })
            }
        }
    }
}

const search = async (query, limit, page) => {
    try {
        const request = getRequest(query, limit, page)
        const response = await fetch(request.url, request.options)
        const data = await response.json()
        if(response.ok){
            store.dispatch({type: "search/fetched", payload : data})
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