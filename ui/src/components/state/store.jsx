import { createStore, combineReducers } from 'redux'
import SearchReducer from './SearchReducer'

const rootReducer = combineReducers({
        searchState : SearchReducer
    }
)

const store = createStore(rootReducer)
export default store
