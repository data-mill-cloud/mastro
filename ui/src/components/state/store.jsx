import { createStore, combineReducers } from 'redux'
import SearchReducer from './SearchReducer'
import AssetDetailReducer from './AssetDetailReducer'
import FeaturesetReducer from './FeaturesetReducer'
import MetricsetReducer from './MetricsetReducer'
import KafkaConnectReducer from './KafkaConnectReducer'
import KafkaSchemaReducer from './KafkaSchemaReducer'

const rootReducer = combineReducers({
        searchState : SearchReducer,
        assetDetailState : AssetDetailReducer,
        featuresetState : FeaturesetReducer,
        metricsetState : MetricsetReducer,
        kafkaconnectState : KafkaConnectReducer,
        kafkaSchemaState : KafkaSchemaReducer,
    }
)

const store = createStore(rootReducer)
export default store
