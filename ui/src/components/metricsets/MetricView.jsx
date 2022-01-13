import DeequMetricsView from './DeequMetricsView'

function MetricView({metric}) {

    if(Object.keys(metric).includes('analyzerContext')) {
        return <DeequMetricsView metric={metric} />
    }
    
}

export default MetricView
