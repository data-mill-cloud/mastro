import {useDispatch, useSelector} from 'react-redux';
import Asset from './Asset'

function Assets(){
    const dispatch = useDispatch()
    const searchState = useSelector(state => state.searchState)
    const loading = searchState.loading
    const assets = searchState.assets
    const pagination = searchState.pagination
    const limit = searchState.limit
   
    const handlePaginationClick = (page) => {
        dispatch({type: 'search/gotopage', payload: parseInt(page)})
    }

    const changeMaxItemsPerPage = (e) => {
        dispatch({type: 'search/resizemaxitems', payload: parseInt(e.target.value)})
    }

    if (!loading){
        return (
            <div>
                <div>
                    {pagination && (
                        <div className="flex justify-between mb-10">
                            <div className="justify-start">
                                Page <strong>{pagination.page}</strong> of <strong>{pagination.totalPage}</strong>
                            </div>
                            <div><strong>{pagination.total}</strong> results found</div>
                            <div className="flex justify-end">
                                <div><input onChange={changeMaxItemsPerPage} type="range" min="1" max="30" defaultValue={limit} className="range range-primary"/></div>
                                <div className="ml-10"><strong>{limit}</strong> results per page</div>
                            </div>
                        </div>        
                    )}
                </div>
                <div className="grid grid-cols-1 gap-8 xl:grid-cols-4 lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-1">
                    {assets.map(asset => (
                        <Asset key={asset.name} asset={asset} />
                    ))}
                </div>
                
                {pagination && pagination.totalPage>0 && (  
                    <div>  
                        <div className="justify-center btn-group mt-10">
                            {pagination.page -1 > 0 && (
                                <button onClick={(e) => handlePaginationClick(pagination.page -1)} className="btn">Previous</button> 
                            ) || (
                                <button className="btn btn-disabled">Previous</button>
                            )}

                            {pagination.page - 3 > 0 && (
                                <button onClick={(e) =>handlePaginationClick(pagination.page - 3)} className="btn">{pagination.page - 3}</button>
                            )}
                            {pagination.page - 2 > 0 && (
                                <button onClick={(e) =>handlePaginationClick(pagination.page - 2)} className="btn">{pagination.page - 2}</button>
                            )}
                            {pagination.page -1 > 0 && (
                                <button onClick={(e) =>handlePaginationClick(pagination.page - 1)} className="btn">{pagination.page - 1}</button>
                            )}

                            <button className="btn btn-active btn-disabled">{pagination.page}</button> 

                            {pagination.page + 1 <= pagination.totalPage && (
                                <button onClick={(e) =>handlePaginationClick(pagination.page + 1)} className="btn">{pagination.page + 1}</button>
                            )}

                            {pagination.page + 2 <= pagination.totalPage && (
                                <button onClick={(e) =>handlePaginationClick(pagination.page + 2)} className="btn">{pagination.page + 2}</button>
                            )}

                            {pagination.page + 3 <= pagination.totalPage && (
                                <button onClick={(e) =>handlePaginationClick(pagination.page + 3)} className="btn">{pagination.page + 3}</button>
                            )}

                            {pagination.page + 1 <= pagination.totalPage && (
                                <button onClick={(e) =>handlePaginationClick(pagination.page + 1)}  className="btn">Next</button>
                            ) || (
                                <button className="btn btn-disabled">Next</button>
                            )}
                        </div>
                    </div> 
                )}
            </div>
        )
    }else{
        return null
    }
}

export default Assets