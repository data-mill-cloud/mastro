function DeequMetricsView({metric}) {
    const { v4: uuidv4 } = require('uuid');
    const metrics = metric.analyzerContext.metricMap
    
    return (
        <div className="shadow stats">
            { metrics.map(m => (
                
                <div key={uuidv4()} className="stat">
                    <div className="stat-title">{`${m.metric.entity} ${m.metric.name} ${m.metric.instance}`}</div> 
                    <div className="stat-value">{m.metric.value}</div>
                    <div className="stat-desc">{m.metric.metricName}</div>
                </div>
            
            ))}
            
        </div>
    )
}

export default DeequMetricsView
