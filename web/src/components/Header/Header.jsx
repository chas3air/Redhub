import styles from './Header.module.css';

export default function Header({title}) {
    return (
        <header className={styles.header}>
            {title}
        </header>
    );
}