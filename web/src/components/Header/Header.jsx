import styles from './Header.module.css';

export default function Header({ title, userRole, onLogout }) {
    return (
        <header className={styles.header} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '10px 40px' }}>
            <h1>{title}</h1>
            <div style={{ display: 'flex', alignItems: 'center' }}>
                {userRole && <span className="badge bg-secondary" style={{ marginRight: '10px' }}>{userRole}</span>}
                {userRole && (
                    <button onClick={onLogout} className="btn btn-danger btn-sm">Logout</button>
                )}
            </div>
        </header>
    );
}