import {useSelector} from 'react-redux'
import Asset from './Asset'

function Assets(){
    const searchState = useSelector(state => state.searchState)
    const loading = searchState.loading
    const assets = searchState.assets
    
    if (!loading){
        return (
            <div className="grid grid-cols-1 gap-8 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                {assets.map(asset => (
                    <Asset key={asset.name} asset={asset} />
                ))}
            </div>
        )
    }else{
        return null
    }
}

export default Assets