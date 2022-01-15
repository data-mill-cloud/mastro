
import {useDispatch} from 'react-redux';

function Slider({pagination, limit, resizeTarget}) {
    const dispatch = useDispatch()

    const changeMaxItemsPerPage = (e) => {
        dispatch({type: resizeTarget, payload: parseInt(e.target.value)})
    }

    return (
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
    )
}

export default Slider
