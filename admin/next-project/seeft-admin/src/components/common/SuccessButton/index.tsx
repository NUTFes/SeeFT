import React from "react";

interface Props {
    text: string;
    onClick?: () => void;
    disabled?: boolean;
    onSubmit?: () => Promise<void>;
}

export const SuccessButton = (props: Props): JSX.Element => {
    const className = 'py-3 px-6 rounded-xl text-white-0 bg-accent-2' + (props.disabled ? ' cursor-not-allowed opacity-20' : '');
    return (
        <button className={className} disabled={props.disabled} onClick={props.onClick} onSubmit={props.onSubmit ? props.onSubmit : undefined}>
            {props.text}
        </button>
    );
}