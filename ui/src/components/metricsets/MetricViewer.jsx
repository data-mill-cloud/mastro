import MetricView from './MetricView';

function MetricViewer({metricset}) {
    const { v4: uuidv4 } = require('uuid');
    
    return (
        metricset.metrics.map(metric => (
            <MetricView key={uuidv4()} metric={metric} />
        ))
    )
}

export default MetricViewer
