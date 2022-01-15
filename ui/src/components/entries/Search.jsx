import { useState, useEffect } from 'react';
import {AiOutlineSearch} from 'react-icons/ai';
import {useDispatch, useSelector} from 'react-redux';
import {BiError} from 'react-icons/bi';

function Search() {
    const [errorMessageText, setErrorMessageText] = useState('')
    const [searchText, setSearchText] = useState('')

    const dispatch = useDispatch()
    const searchState = useSelector(state => state.searchState)
    const loading = searchState.loading

    useEffect(() => {
        setErrorMessageText(searchState.errorMessage && searchState.errorMessage.length > 0 ? searchState.errorMessage : '')
    }, [searchState.errorMessage])

    useEffect(() => {
        setSearchText(searchState.query)
    }, [searchState.query])
    
    const handleChange = (e) => {
        setSearchText(e.target.value);
    }

    const handleSubmit = (e) => {
        e.preventDefault();
        if (searchText === '') {
            setErrorMessageText("Please enter a search term!")
        }else{
            setErrorMessageText("")
            dispatch({type: 'search/submit', payload: searchText})
            //setSearchText('');
        }
    }

    const handleClearSearch = (e)=> {
        dispatch({type: 'search/clear'})
        setSearchText('')
    } 

    return (
        <div>
            <div className="grid grid-cols-1 xl:grid-cols-2 lg:grid-cols-2 md:grid-cols-2 mb-8 gap-8 justify-center">
                <div>
                    <form onSubmit={handleSubmit}>
                        <div className="form-control">
                            <div className="relative">
                                <input type="text" 
                                    className="w-full pr-40 bg-gray-200 input input-lg text-black"
                                    placeholder="Search"
                                    onChange={handleChange}
                                    value={searchText}
                                />
                                <button type="submit" className={`absolute top-0 right-0 rouded-l-none w-36 btn btn-lg ${loading ? 'loading' : ''}`}>
                                    <AiOutlineSearch className="text-4xl" />
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
                { searchText !== "" && (
                    <div>
                        <button onClick={handleClearSearch} className="btn btn-ghost btn-lg rounded-btn">Clear</button>
                    </div>
                )}   
            </div>
            { errorMessageText !== "" && (
                <div className="alert alert-error">
                    <div className="flex-1">
                        <BiError className="text-2xl" />
                        <label>{errorMessageText}</label>
                    </div>
                </div>
            )}
        </div>
    )
}

export default Search
