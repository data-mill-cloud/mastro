import { useState } from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {BiError} from 'react-icons/bi';
import {GrStatusGood} from 'react-icons/gr'
import {MdInput} from 'react-icons/md'
import {VscDebugRestart} from 'react-icons/vsc'

function Connector({name, connector}) {
    const dispatch = useDispatch()
    const [selectedTab, setSelectedTab] = useState('status')
    const { v4: uuidv4 } = require('uuid');

    const onTabClick = (target) => setSelectedTab(target)

    const handleRestart = () => {
        if (window.confirm(`Are you sure you want to restart all failing tasks for the ${name} connector?`)) {
            dispatch({type: 'kafkaconnect/restart', payload: name})
        }
    }

    return (
        <div className="card card card-bordered shadow-2xl compact side bg-base-100 border border-2 hover:border-gray-400 mb-2">
            <div className="card-body" style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                <div className="flex justify-between">
                    <div className="justify-center"><strong>{name}</strong></div>
                    <div className="flex justify-end">
                        <div style={connector.info.type === 'sink' ? {transform: 'scaleX(-1)'} : {}}><MdInput className="text-xl"/></div>
                        <strong className="ml-1">{connector.info.type}</strong>
                    </div>
                </div>
                <div className="justify-end mb-1">{connector.info.config['connector.class']}</div>
                <div className="tabs">
                    <a onClick={(e) => onTabClick('status')} className={`tab tab-bordered ${selectedTab === "status" ? 'tab-active':''}`}><strong>Status</strong></a>                 
                    <a onClick={(e) => onTabClick('info')} className={`tab tab-bordered ${selectedTab === "info" ? 'tab-active':''}`}><strong>Info</strong></a>
                </div>
                <div style={{marginLeft: '0.7rem', marginRight: '0.7rem'}}>
                    {(selectedTab === 'status') && (
                        <div>
                            <table className="table w-full table-compact">
                                <tbody>
                                    <tr><th>Id</th><th>Worker Id</th><th>State</th></tr>
                                {connector.status.tasks.map(task => (
                                <tr key={uuidv4()}>
                                    <td>{task.id}</td>
                                    <td>{task.worker_id}</td>
                                    <th>
                                        <div className={`alert `+(task.state !== 'RUNNING' ? `alert-error` :`alert-success`)}>
                                            {task.state === 'RUNNING' && <GrStatusGood className="text-2xl" />}
                                            {task.state !== 'RUNNING' && <BiError className="text-2xl"/>}
                                            <label className="ml-1">{task.state}</label>
                                        </div>
                                    </th>
                                    </tr>
                                ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                    {(selectedTab === 'info') && (
                        <div>
                            <div><strong>Config:</strong></div>
                            <div>
                                <textarea defaultValue={`${JSON.stringify(connector.info.config, null, 2)}`} className="textarea w-full h-full textarea-bordered textarea-primary"></textarea>
                            </div>
                        </div>
                    )}
                </div>
                <div className="flex justify-end">
                    <VscDebugRestart onClick={handleRestart} className="text-xl" style={{cursor: 'pointer'}} />
                </div>
            </div>
        </div>
    )
}

export default Connector
