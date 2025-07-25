import styles from './User.module.css';

export default function User({ user, onDelete }) {
    return (
        <div key={user.id} className={styles.userBlock}>
            <div className={styles.info}>
                <div className={styles.nick}><strong>Ник:</strong> {user.nick}</div>
                <div className={styles.email}><strong>Email:</strong> {user.email}</div>
            </div>
            <button className={styles.deleteButton} onClick={() => onDelete(user.id)}>Удалить</button> {/* Кнопка удаления */}
        </div>
    );
}