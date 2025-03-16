import Button from "../Button/Button";

const AuthChooser = ({ active, onChange }) => {
    return (
        <section>
            <Button
                isActive={active === "login"}
                onClick={() => onChange("login")}
            >
                Вход
            </Button>
            <Button
                isActive={active === "register"}
                onClick={() => onChange("register")}
            >
                Регистрация
            </Button>
        </section>
    );
};

export default AuthChooser;