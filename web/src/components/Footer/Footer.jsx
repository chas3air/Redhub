import styles from "./Footer.module.css";

export default function Footer() {
    return (
        <footer className={styles.footer}>
            <div className={styles.info}>
                <ul>
                    <li>Github: https://github.com/chas3air</li>
                    <li>telegram: @chas3air</li>
                </ul>
            </div>
        </footer>
    );
}