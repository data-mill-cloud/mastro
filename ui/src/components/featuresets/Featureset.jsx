import {Link} from 'react-router-dom';
import {FaFileCsv} from 'react-icons/fa'
import { CSVLink } from "react-csv";

function Featureset({featureset}) {
    const { v4: uuidv4 } = require('uuid');
    const featureKeys = Array.from(new Set(featureset.features.map(feature => Object.keys(feature)).flat().sort()));
    const csvData = featureset.features //.map(feature => featureKeys.map(fk => feature[fk]))
    
    return (
        <div className="card shadow-2xl compact side bg-base-100 border-2 hover:border-gray-400">
            <div className="card-body" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                <div className="grid grid-cols-4 md:gap-8">
                    <div className="items-start flex-row col-span-3">
                        <p className="text-base-content text-opacity-40">{featureset["inserted_at"]}</p>
                        <Link className="card-title" to={`/featureset/${featureset.name}`}>
                            {featureset.version}
                        </Link>
                    </div>
                    <div className="items-end col-span-1">                        
                        <CSVLink filename={`${featureset.name}_${featureset.version}.csv`} headers={featureKeys} data={csvData} separator={","} target="_blank"><FaFileCsv className="text-2xl" /></CSVLink>
                    </div>
                </div>
                
            </div>
            <div className="card-body flex-row" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                <p style={{textAlign: 'justify'}}>{featureset.description}</p>
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
                            { featureset.labels !== undefined ? Object.keys(featureset.labels).map(key => (<tr key={uuidv4()}><th key={uuidv4()}>{key}</th><td>{featureset.labels[key]}</td></tr>)): '' }
                        </tbody>
                    </table>
                </div>
                <div className="card-body flex-row" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                    <table className="text-center table-compact w-full">
                        <thead>
                        <tr>{ featureKeys.map(fk => <th key={uuidv4()}>{fk}</th> )}</tr>
                        </thead> 
                        <tbody>
                            {featureset.features.map(feature => <tr key={uuidv4()}>{featureKeys.map(fk => <td key={uuidv4()}>{feature[fk]}</td>)}</tr>)}
                        </tbody>
                    </table>
                </div>
            </div>
        </div> 
    )
}

export default Featureset
