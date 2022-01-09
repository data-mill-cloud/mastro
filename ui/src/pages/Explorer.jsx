import {getSvcHost} from "../SvcUtils.js"
import { useState } from 'react';
import {BiError} from 'react-icons/bi';

function Explorer() {
    const [errorMessage, setErrorMessage] = useState('')
    const [method, setMethod] = useState('GET')
    const [selectedService, setSelectedService] = useState(getSvcHost('catalogue'))
    const [path, setPath] = useState('/')
    const [header, setHeader] = useState('')
    const [body, setBody] = useState('')
    const [responseBody, setResponseBody] = useState('')

    const handleMethodChange = (e) => setMethod(e.target.value)
    const handleSelectedServiceChange = (e) => setSelectedService(e.target.value)
    const handlePathChange = (e) => setPath(e.target.value)
    const handleHeaderChange = (e) => setHeader(e.target.value)
    const handleBodyChange = (e) => setBody(e.target.value)

    const sendRequest = async () => {
        try {
            const request = {
                url : `${selectedService}${path}`,
                options : { 
                    method : `${method}`,
                    headers : new Map(header ? Object.entries(JSON.parse(`${header}`)) : []),
                }
            }
            if(body) {
                request.options.body = body
            }
            const response = await fetch(request.url, request.options)
            const data = await response.json()

            if(response.ok){
                setResponseBody(JSON.stringify(data, null, 2))
            }else{
                setErrorMessage(`${response.statusText}: ${data.message}`)
            }
        }catch(error){
            setErrorMessage(error.message)
        }
    }

    const handleSend = (e) => {
        e.preventDefault();
        setErrorMessage('')
        sendRequest()
    }


    return (
        <div>
            <div className="grid grid-cols-1 xl:grid-cols-6 lf:grid-cols-6 md:grid-cols-6 mb-8 md:gap-8">
                <div className="col-span-3 card-body rounded-lg shadow-md bg-base-100 border-2">
                    <form onSubmit={handleSend}>
                        <h2 className="card-title">Request</h2>
                        <div className="grid grid-cols-1 xl:grid-cols-5 lf:grid-cols-5 md:grid-cols-5 mb-6 md:gap-6 content-end">
                            <div className="col-span-1 form-control">
                                <label className="label">
                                    <span className="label-text">Method</span>
                                </label> 
                                <select defaultValue={method} onChange={handleMethodChange} className="select select-bordered select-primary w-full max-w-xs">
                                    <option>GET</option> 
                                    <option>POST</option> 
                                    <option>PUT</option> 
                                </select>
                            </div>
                            <div className="col-span-2">
                                <label className="label">
                                    <span className="label-text">Service</span>
                                </label> 
                                <select defaultValue={selectedService} onChange={handleSelectedServiceChange} className="select select-bordered select-primary w-full max-w-xs">
                                    <option value={getSvcHost('catalogue')}>Catalogue @ {getSvcHost('catalogue')}</option> 
                                    <option value={getSvcHost('featurestore')}>Featurestore @ {getSvcHost('featurestore')}</option> 
                                    <option value={getSvcHost('metricstore')}>Metricstore @ {getSvcHost('metricstore')}</option> 
                                </select>
                            </div>
                            <div className="col-span-2 form-control">
                                <label className="label">
                                    <span className="label-text">Path</span>
                                </label> 
                                <input type="text" onChange={handlePathChange} defaultValue={path} className="input input-bordered input-primary"/>
                            </div>
                        </div>
                        <div className="form-control">
                            <label className="label">
                                <span className="label-text">Header</span>
                            </label> 
                            <textarea  onChange={handleHeaderChange} className="textarea h-24 textarea-bordered textarea-primary" defaultValue={header}></textarea>
                        </div>
                        <div className="form-control">
                            <label className="label">
                                <span className="label-text">Body</span>
                            </label> 
                            <textarea onChange={handleBodyChange} className="textarea h-24 textarea-bordered textarea-primary" defaultValue={body}></textarea>
                        </div>
                        <div className="flex justify-center">
                            <button type="submit" className="btn btn-primary">Send</button> 
                        </div>
                    </form>
                </div>
                <div className="col-span-3 card-body rounded-lg shadow-md bg-base-100 border-2">
                    <h2 className="card-title">Response by {`${selectedService}${path}`}</h2>
                    <div className="form-control h-full">
                        { errorMessage !== "" && (
                            <div className="alert alert-error">
                                <div className="flex-1">
                                    <BiError className="text-2xl" />
                                    <label>{errorMessage}</label>
                                </div>
                            </div>
                        ) || (
                            <div className="form-control h-full">
                                <label className="label">
                                    <span className="label-text">Body</span>
                                </label> 
                                <textarea defaultValue={responseBody} className="textarea h-full textarea-bordered textarea-primary"></textarea>
                            </div>
                        )}
                    </div>
                </div>
            </div>    
        </div>
    )
}

export default Explorer
