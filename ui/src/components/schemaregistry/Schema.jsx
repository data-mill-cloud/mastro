import { useState } from 'react';
import {useDispatch} from 'react-redux';

function Schema({name, schema}) {
    const dispatch = useDispatch()
    //const [selectedVersion, setSelectedVersion] = useState(1)
    const formattedSchema = JSON.stringify(JSON.parse(schema.schema), null, 2)

    return (
        <div className="h-auto card card card-bordered shadow-2xl compact side bg-base-100 border border-2 hover:border-gray-400 mb-2">
            <div className="card-body" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                <div className="flex justify-between">
                    <div className="justify-center"><strong>{name}</strong></div>
                    <div className="flex justify-end">
                        <strong className="ml-1">v. {schema.version}</strong>
                    </div>
                </div>
                <div className="form-control h-full">
                    <label className="label">
                        <span className="label-text">Schema</span>
                    </label> 
                    <textarea defaultValue={formattedSchema} className="h-full w-full textarea textarea-bordered textarea-primary"></textarea>
                </div>
            </div>
        </div>
    )
}

export default Schema
