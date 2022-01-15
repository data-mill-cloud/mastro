import {useDispatch, useSelector} from 'react-redux';

function Metricset({metricset}) {
    const { v4: uuidv4 } = require('uuid');
    const dispatch = useDispatch()
    const metricsetState = useSelector(state => state.metricsetState)
    const selectedMetricset = metricsetState.selectedMetricset

    const handleSelect = (e) => {
        e.preventDefault();
        if(metricset !== selectedMetricset){
            dispatch({type: 'metricset/select', payload: metricset})
        }else{
            dispatch({type: 'metricset/unselect'})
        }
    }

    return (
        <div onClick={handleSelect} className={`card card-bordered shadow-2xl compact side bg-base-100 border border-2 hover:border-gray-400 ` + (metricset === selectedMetricset ? 'border-gray-400': 'border-gray-200')}>
            <div className="card-body" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                <div className="grid grid-cols-4 md:gap-8">
                    <div className="items-start flex-row col-span-3">
                        <p className="text-base-content text-opacity-40">{metricset["inserted_at"]}</p>
                        <p className="card-title">{metricset.version}</p>
                    </div>
                </div>
            </div>
            <div className="card-body flex-row" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                <p style={{textAlign: 'justify'}}>{metricset.description}</p>
            </div>
            <div className="overflow-auto overflow-scroll max-h-80">
                <div className="card-body flex-row" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                    <table className="table-compact w-full">
                        <thead>
                        <tr>
                            <th>Key</th> 
                            <th>Value</th> 
                        </tr>
                        </thead> 
                        <tbody>
                            { metricset.labels !== undefined ? Object.keys(metricset.labels).map(key => (<tr key={uuidv4()}><th>{key}</th><td>{metricset.labels[key]}</td></tr>)): '' }
                        </tbody>
                    </table>
                </div>
            </div>
        </div> 
    )
}

export default Metricset
