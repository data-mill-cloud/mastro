import {useDispatch} from 'react-redux';

function PageSelector({pagination, pageTarget}) {
    const dispatch = useDispatch()

    const handlePaginationClick = (page) => {
        dispatch({type: pageTarget, payload: parseInt(page)})
    }

    if(pagination && pagination.totalPage > 0){
        return (
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
        )
    }else{
        return null
    }
}

export default PageSelector
