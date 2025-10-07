import React from "react";
import "./Badge.css";

interface BadgeProps {
    label: string;
    value: string;
}

const Badge: React.FC<BadgeProps> = ({label, value}) => {
    return (
        <div className="badge">
            <div className="badge-label"><span>{label}</span></div>
            <div className="badge-value"><span>{value}</span></div>
        </div>
    );
};

export default Badge;
