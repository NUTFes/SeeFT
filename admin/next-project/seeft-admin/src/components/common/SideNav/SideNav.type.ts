import { ReactNode } from "react";

export interface NavItemProps {
    icon: ReactNode;
    children?: React.ReactNode;
    href: string;
    currentPath: string;
    isParent?: boolean;
    isShow?: boolean;
    onClick?: () => void;
}