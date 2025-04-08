import React from 'react';
import '../assets/Logo.css';
import logoImage from '../assets/images/logo.png';

interface LogoProps {
    collapsed: boolean;
}

const Logo: React.FC<LogoProps> = ({ collapsed }) => {
    return (
        <div className={`app-logo ${collapsed ? 'collapsed' : ''}`}>
            <img src={logoImage} alt="BeePF Logo" />
        </div>
    );
};

export default Logo; 