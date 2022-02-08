function FeatureViewer({featureset}) {
    const { v4: uuidv4 } = require('uuid');

    return ( 
        <div className="shadow stats">
            { featureset.features.map(f => (
                <div key={uuidv4()} className="stat">
                    <div className="stat-title">{f.name}</div> 
                    {!f['data_type'].includes("dataframe") && (<div className="stat-value">{`${f.value}`}</div>)}
                    <div className="stat-desc">{f['data_type']}</div>
                </div>
            ))}
        </div>
    )
}

export default FeatureViewer
