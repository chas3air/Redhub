import styles from './Article.module.css'

export default function Article({article}) {
    return (
        <div key={article.id} className={styles.block}>
            <h2>{article.title}</h2>
            <p><strong>Создано:</strong> {new Date(article.created_at).toLocaleDateString()}</p>
            <p>{article.content}</p>
            <p><strong>Владелец:</strong> {article.owner_id}</p>
        </div>
    );
}