import React from "react";

interface Props {
    text: string;
    onClick?: () => void;
}

export const FailerButton = (props: Props): JSX.Element => {
    return(
        <button className='py-3 px-6 rounded-xl text-white-0 bg-caution' onClick={props.onClick}>
            {props.text}
        </button>
    );
}