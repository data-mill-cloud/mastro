import {useSelector} from 'react-redux';
import Asset from './Asset'
import Slider from '../pagination/Slider'
import PageSelector from '../pagination/PageSelector'

function Assets(){
    const searchState = useSelector(state => state.searchState)
    const loading = searchState.loading
    const assets = searchState.assets
    const pagination = searchState.pagination
    const limit = searchState.limit
   
    if (!loading){
        return (
            <div>
                <Slider pagination={pagination} limit={limit} resizeTarget={'search/resizemaxitems'} />
                <div className="grid grid-cols-1 gap-8 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                    {assets.map(asset => (
                        <Asset key={asset.name} asset={asset} />
                    ))}
                </div>
                <PageSelector pagination={pagination} pageTarget={'search/gotopage'}/> 
            </div>
        )
    }else{
        return null
    }
}

export default Assets