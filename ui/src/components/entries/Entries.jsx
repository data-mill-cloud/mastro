import {useSelector} from 'react-redux'
import Entry from './Entry'
import Spinner from '../layout/Spinner'

function Entries(){
    const searchState = useSelector(state => state.searchState)
    const loading = searchState.loading
    const entries = searchState.entries
    
    /*
    useEffect(() => {
        fetchEntries()
    }, [])

    const fetchEntries = async () => {
        const response = await fetch('/api/entries')
        const data = await response.json()
        setEntries(data)
        setLoading(false)

        dispatch({type: 'search/fetched', payload: data})
    }
    */
    if (loading){
        return (
            <Spinner/>
        )
    }else{
        
        return (
            <div className="grid grid-cols-1 gap-8 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                {entries.map(entry => (
                    <Entry key={entry.name} entry={entry} />
                ))}
            </div>
        )
    }
}

export default Entries