import styles from './Button.module.css';

export default function Button({children, onClick, isActive}) {
    return (
        <button className={`${styles.button} ${isActive ? styles.active : ''}`} onClick={onClick}>
            {children}
        </button>
    );
}