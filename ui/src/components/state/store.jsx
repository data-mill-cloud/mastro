import { createStore, combineReducers } from 'redux'
import SearchReducer from './SearchReducer'
import FeaturesetReducer from './FeaturesetReducer'

const rootReducer = combineReducers({
        searchState : SearchReducer,
        featuresetState : FeaturesetReducer,
    }
)

const store = createStore(rootReducer)
export default store
