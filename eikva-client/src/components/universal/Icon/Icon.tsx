import React from "react";
import "./Icon.css";

interface IconProps
    extends React.DetailedHTMLProps<
        React.HTMLAttributes<HTMLSpanElement>,
        HTMLSpanElement
    > {
    name: string;
}

const Icon: React.FC<IconProps> = ({ className, name = "star", ...props }) => {
    return (
        <span className={`icon${className ? " " + className : ""}`} {...props}>
            {name}
        </span>
    );
};

export default Icon;
