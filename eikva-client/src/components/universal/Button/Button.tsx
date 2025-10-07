import React from "react";
import Icon from "../Icon/Icon";
import "./Button.css";

interface ButtonProps
    extends React.DetailedHTMLProps<
        React.ButtonHTMLAttributes<HTMLButtonElement>,
        HTMLButtonElement
    > {
    icon?: string;
}

const Button: React.FC<ButtonProps> = ({
    className,
    children,
    type = "button",
    icon,
    ...props
}) => {
    return (
        <button
            className={`button display-flex gap-small align-items-center${className ? " " + className : ""}`}
            type={type}
            {...props}
        >
            {icon && <Icon name={icon} />}
            {children && <div>{children}</div>}
        </button>
    );
};

export default Button;
