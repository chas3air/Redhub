import Button from '../Button/Button'

export default function TabChoser({active, onChange}) {
    return (
        <section>
            <Button
                isActive={active === "articles"}
                onClick={() => onChange("articles")}>
                Articles
            </Button>
            <Button
                isActive={active === "profile"}
                onClick={() => onChange("profile")}>
                Profile
            </Button>
        </section>
    );
}