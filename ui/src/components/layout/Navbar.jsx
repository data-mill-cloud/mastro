import {Link} from 'react-router-dom'
import PropTypes from 'prop-types'
import logo from '../../img/logo.png';

function Navbar({title}) {
    return (
        <nav className='navbar mb-12 shadow-lg bg-neutral text-neutral-content'>
            <div className='container mx-auto'>  
                <div className='flex-none px-2 mx-2'>
                    <img className='inline pr-2 text-3xl' src={logo} height={48} width={48} alt='logo' />
                    <Link to='/' className='text-lg font-bold align-middle'>{title}</Link>
                </div>
                

                <div className="flex-1 px-2 mx-2">
                    <div className="flex justify-end">
                        <Link to='/' className='btn btn-ghost btn-sm rounded-btn'>Search</Link>
                        <Link to='/explorer' className='btn btn-ghost btn-sm rounded-btn'>Explore</Link>
                        <Link to='/about' className='btn btn-ghost btn-sm rounded-btn'>About</Link>
                    </div>
                </div>
            </div>
        </nav>
    )
}

Navbar.defaultProps = {
    title: 'MASTRO',
}

Navbar.propTypes = {
    title : PropTypes.string,
}

export default Navbar
