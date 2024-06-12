import React from "react";

interface Props {
    text: string;
    onClick?: () => void;
    disabled?: boolean;
}

export const SuccessButton = (props: Props): JSX.Element => {
    const className = 'py-3 px-6 rounded-xl text-white-0 bg-accent-2' + (props.disabled ? ' cursor-not-allowed opacity-20' : '');
    return (
        <button className={className} disabled={props.disabled} onClick={props.onClick}>
            {props.text}
        </button>
    );
}