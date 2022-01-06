import { useState, useEffect } from 'react';
import {AiOutlineSearch} from 'react-icons/ai';
import {useDispatch, useSelector} from 'react-redux';

function Search() {
    const [errorMessageText, setErrorMessageText] = useState('')
    const [searchText, setSearchText] = useState('')

    const dispatch = useDispatch()
    const searchState = useSelector(state => state.searchState)

    useEffect(() => {
        setErrorMessageText(searchState.errorMessage && searchState.errorMessage.length > 0 ? searchState.errorMessage : '')
    }, [searchState.errorMessage])

    /*
    useEffect(() => {
        setSearchText(searchState.searchText)
    }, [searchState.searchText])
    */

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

    return (
        <div>
            <div className="grid grid-cols-1 xl:grid-cols-2 lg:grid-cols-2 md:grid-cols-2 mb-8 gap-8">
                <div>
                    <form onSubmit={handleSubmit}>
                        <div className="form-control">
                            <div className="relative">
                                <input type="text" 
                                    className="w-full pr-40 bg-gray-200 input input-lg text-black"
                                    placeholder="Search"
                                    onChange={handleChange}
                                />
                                <button type="submit" className="absolute top-0 right-0 rouded-l-none w-36 btn btn-lg">
                                    <AiOutlineSearch className="text-4xl" />
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
                { searchText !== "" && (
                    <div>
                        <button className="btn btn-ghost btn-lg rounded-btn">Clear</button>
                    </div>
                )}   
            </div>
            { errorMessageText !== "" && (
                <div className="alert alert-error">
                    <div className="flex-1">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-6 h-6 mx-2 stroke-current">    
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path>                      
                        </svg> 
                        <label>{errorMessageText}</label>
                    </div>
                </div>
            )}
        </div>
    )
}

export default Search
