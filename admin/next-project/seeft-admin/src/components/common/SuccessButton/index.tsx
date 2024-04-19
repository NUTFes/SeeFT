import React from "react";

interface Props {
    text: string;
    onClick?: () => void;
}

export const SuccessButton = (props: Props): JSX.Element => {
    return(
        <button className='py-3 px-6 rounded-xl text-white-0 bg-accent-2' onClick={props.onClick}>
            {props.text}
        </button>
    );
}