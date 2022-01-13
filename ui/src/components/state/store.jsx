import { createStore, combineReducers } from 'redux'
import SearchReducer from './SearchReducer'
import FeaturesetReducer from './FeaturesetReducer'
import MetricsetReducer from './MetricsetReducer'

const rootReducer = combineReducers({
        searchState : SearchReducer,
        featuresetState : FeaturesetReducer,
        metricsetState : MetricsetReducer,
    }
)

const store = createStore(rootReducer)
export default store
