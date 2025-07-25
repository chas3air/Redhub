import { useNavigate } from 'react-router-dom';
import styles from './Header.module.css';

export default function Header({ title, userRole, onLogout }) {
    const navigate = useNavigate();

    const handleProfileRedirect = () => {
        navigate('/profile');
    };

    const handleTitleClick = () => {
        navigate('/articles');
    };

    const handleLoginRedirect = () => {
        navigate('/login');
    };

    return (
        <header className={styles.header} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '10px 40px' }}>
            <h1 onClick={handleTitleClick} style={{ cursor: 'pointer' }}>{title}</h1>
            <div style={{ display: 'flex', alignItems: 'center' }}>
                {userRole ? (
                    <>
                        <button onClick={handleProfileRedirect} className="badge bg-secondary" style={{ marginRight: '10px' }}>
                            {userRole}
                        </button>
                        <button onClick={onLogout} className="btn btn-danger btn-sm">Logout</button>
                    </>
                ) : (
                    <button onClick={handleLoginRedirect} className="btn btn-primary btn-sm">Login</button>
                )}
            </div>
        </header>
    );
}