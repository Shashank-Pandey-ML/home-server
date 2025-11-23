import { useAuth } from '../context/AuthContext';

function ProfileHeader() {
    const { user } = useAuth();
    
    const handleScrollToPortfolio = (e) => {
        e.preventDefault();
        const element = document.querySelector('#portfolio');
        if (element) {
            const offset = 70;
            const elementPosition = element.getBoundingClientRect().top;
            const offsetPosition = elementPosition + window.pageYOffset - offset;

            window.scrollTo({
                top: offsetPosition,
                behavior: 'smooth'
            });
        }
    };

    return (
        <header id="home" className="header">
            <div className="overlay"></div>
            <div className="header-content container">
                <h1 className="header-title">
                    <span className="up">HI!</span>
                    <span className="down">{user ? 'Welcome Shashank!' : 'I am Shashank Pandey'}</span>
                </h1>
                <p className="header-subtitle">BACKEND DEVELOPER</p>            
                <button 
                    className="btn btn-primary" 
                    onClick={handleScrollToPortfolio}
                >
                    Visit My Works
                </button>
            </div>              
        </header>
    );
}

export default ProfileHeader;
